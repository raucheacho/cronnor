package storage

import (
	"database/sql"
	"fmt"

	"github.com/rauche/cronnor/internal/models"
)

// CreateJobLog creates a new job log entry
func (r *Repository) CreateJobLog(log models.JobLog) error {
	query := `
		INSERT INTO job_logs (job_id, status, http_code, duration_ms, response_body, error_message)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(query, log.JobID, log.Status, log.HTTPCode, log.DurationMs, log.ResponseBody, log.ErrorMessage)
	if err != nil {
		return fmt.Errorf("failed to create job log: %w", err)
	}

	return nil
}

// GetJobLogs retrieves logs for a specific job
func (r *Repository) GetJobLogs(jobID int64, limit int) ([]models.JobLog, error) {
	if limit <= 0 {
		limit = 50
	}

	query := `
		SELECT id, job_id, status, http_code, duration_ms, response_body, error_message, created_at
		FROM job_logs
		WHERE job_id = ?
		ORDER BY created_at DESC
		LIMIT ?
	`

	rows, err := r.db.Query(query, jobID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query job logs: %w", err)
	}
	defer rows.Close()

	var logs []models.JobLog
	for rows.Next() {
		var log models.JobLog
		err := rows.Scan(
			&log.ID, &log.JobID, &log.Status, &log.HTTPCode,
			&log.DurationMs, &log.ResponseBody, &log.ErrorMessage, &log.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan job log: %w", err)
		}
		logs = append(logs, log)
	}

	return logs, rows.Err()
}

// GetLatestJobLog retrieves the most recent log for a job
func (r *Repository) GetLatestJobLog(jobID int64) (*models.JobLog, error) {
	query := `
		SELECT id, job_id, status, http_code, duration_ms, response_body, error_message, created_at
		FROM job_logs
		WHERE job_id = ?
		ORDER BY created_at DESC
		LIMIT 1
	`

	var log models.JobLog
	err := r.db.QueryRow(query, jobID).Scan(
		&log.ID, &log.JobID, &log.Status, &log.HTTPCode,
		&log.DurationMs, &log.ResponseBody, &log.ErrorMessage, &log.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No logs yet
		}
		return nil, fmt.Errorf("failed to get latest job log: %w", err)
	}

	return &log, nil
}
