package storage

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/rauche/cronnor/internal/models"
)

// GetAllJobs retrieves all jobs
func (r *Repository) GetAllJobs() ([]models.Job, error) {
	query := `
		SELECT id, name, cron_expr, url, method, payload, is_active, 
		       created_at, last_run_at, last_status
		FROM jobs
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query jobs: %w", err)
	}
	defer rows.Close()

	var jobs []models.Job
	for rows.Next() {
		var job models.Job
		err := rows.Scan(
			&job.ID, &job.Name, &job.CronExpr, &job.URL, &job.Method,
			&job.Payload, &job.IsActive, &job.CreatedAt, &job.LastRunAt, &job.LastStatus,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan job: %w", err)
		}
		jobs = append(jobs, job)
	}

	return jobs, rows.Err()
}

// GetActiveJobs retrieves all active jobs
func (r *Repository) GetActiveJobs() ([]models.Job, error) {
	query := `
		SELECT id, name, cron_expr, url, method, payload, is_active, 
		       created_at, last_run_at, last_status
		FROM jobs
		WHERE is_active = 1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query active jobs: %w", err)
	}
	defer rows.Close()

	var jobs []models.Job
	for rows.Next() {
		var job models.Job
		err := rows.Scan(
			&job.ID, &job.Name, &job.CronExpr, &job.URL, &job.Method,
			&job.Payload, &job.IsActive, &job.CreatedAt, &job.LastRunAt, &job.LastStatus,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan job: %w", err)
		}
		jobs = append(jobs, job)
	}

	return jobs, rows.Err()
}

// GetJob retrieves a job by ID
func (r *Repository) GetJob(id int64) (*models.Job, error) {
	query := `
		SELECT id, name, cron_expr, url, method, payload, is_active, 
		       created_at, last_run_at, last_status
		FROM jobs
		WHERE id = ?
	`

	var job models.Job
	err := r.db.QueryRow(query, id).Scan(
		&job.ID, &job.Name, &job.CronExpr, &job.URL, &job.Method,
		&job.Payload, &job.IsActive, &job.CreatedAt, &job.LastRunAt, &job.LastStatus,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("job not found")
		}
		return nil, fmt.Errorf("failed to get job: %w", err)
	}

	return &job, nil
}

// CreateJob creates a new job
func (r *Repository) CreateJob(params models.CreateJobParams) (int64, error) {
	query := `
		INSERT INTO jobs (name, cron_expr, url, method, payload)
		VALUES (?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(query, params.Name, params.CronExpr, params.URL, params.Method, params.Payload)
	if err != nil {
		return 0, fmt.Errorf("failed to create job: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get insert id: %w", err)
	}

	return id, nil
}

// UpdateJob updates an existing job
func (r *Repository) UpdateJob(params models.UpdateJobParams) error {
	query := `
		UPDATE jobs
		SET name = ?, cron_expr = ?, url = ?, method = ?, payload = ?
		WHERE id = ?
	`

	result, err := r.db.Exec(query, params.Name, params.CronExpr, params.URL, params.Method, params.Payload, params.ID)
	if err != nil {
		return fmt.Errorf("failed to update job: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("job not found")
	}

	return nil
}

// ToggleJob toggles the is_active status of a job
func (r *Repository) ToggleJob(id int64) error {
	query := `
		UPDATE jobs
		SET is_active = NOT is_active
		WHERE id = ?
	`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to toggle job: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("job not found")
	}

	return nil
}

// UpdateJobStatus updates the last run status of a job
func (r *Repository) UpdateJobStatus(id int64, status string) error {
	query := `
		UPDATE jobs
		SET last_run_at = ?, last_status = ?
		WHERE id = ?
	`

	_, err := r.db.Exec(query, time.Now(), status, id)
	if err != nil {
		return fmt.Errorf("failed to update job status: %w", err)
	}

	return nil
}

// DeleteJob deletes a job
func (r *Repository) DeleteJob(id int64) error {
	query := `DELETE FROM jobs WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete job: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("job not found")
	}

	return nil
}
