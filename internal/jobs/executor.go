package jobs

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rauche/cronnor/internal/models"
	"github.com/rauche/cronnor/internal/storage"
)

// Executor handles HTTP job execution
type Executor struct {
	repo   *storage.Repository
	client *http.Client
}

// NewExecutor creates a new executor
func NewExecutor(repo *storage.Repository) *Executor {
	return &Executor{
		repo: repo,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Execute runs a job and logs the result
func (e *Executor) Execute(job models.Job) error {
	start := time.Now()

	// Prepare request
	var body io.Reader
	if job.Payload.Valid && job.Payload.String != "" {
		body = bytes.NewBufferString(job.Payload.String)
	}

	req, err := http.NewRequest(job.Method, job.URL, body)
	if err != nil {
		return e.logError(job.ID, start, fmt.Errorf("failed to create request: %w", err))
	}

	// Set headers
	if job.Payload.Valid && job.Payload.String != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("User-Agent", "Cronnor/1.0")

	// Execute request
	resp, err := e.client.Do(req)
	if err != nil {
		return e.logError(job.ID, start, fmt.Errorf("failed to execute request: %w", err))
	}
	defer resp.Body.Close()

	// Read response
	responseBody, err := io.ReadAll(io.LimitReader(resp.Body, 10*1024)) // Limit to 10KB
	if err != nil {
		return e.logError(job.ID, start, fmt.Errorf("failed to read response: %w", err))
	}

	duration := time.Since(start).Milliseconds()

	// Determine status
	status := "SUCCESS"
	if resp.StatusCode >= 400 {
		status = "FAILED"
	}

	// Log execution
	log := models.JobLog{
		JobID:      job.ID,
		Status:     status,
		HTTPCode:   sql.NullInt64{Int64: int64(resp.StatusCode), Valid: true},
		DurationMs: sql.NullInt64{Int64: duration, Valid: true},
		ResponseBody: sql.NullString{
			String: string(responseBody),
			Valid:  len(responseBody) > 0,
		},
	}

	if err := e.repo.CreateJobLog(log); err != nil {
		return fmt.Errorf("failed to create job log: %w", err)
	}

	// Update job status
	if err := e.repo.UpdateJobStatus(job.ID, status); err != nil {
		return fmt.Errorf("failed to update job status: %w", err)
	}

	return nil
}

// logError logs an error execution
func (e *Executor) logError(jobID int64, start time.Time, err error) error {
	duration := time.Since(start).Milliseconds()

	log := models.JobLog{
		JobID:      jobID,
		Status:     "ERROR",
		DurationMs: sql.NullInt64{Int64: duration, Valid: true},
		ErrorMessage: sql.NullString{
			String: err.Error(),
			Valid:  true,
		},
	}

	if logErr := e.repo.CreateJobLog(log); logErr != nil {
		return fmt.Errorf("failed to log error: %w (original error: %v)", logErr, err)
	}

	if statusErr := e.repo.UpdateJobStatus(jobID, "ERROR"); statusErr != nil {
		return fmt.Errorf("failed to update status: %w (original error: %v)", statusErr, err)
	}

	return err
}
