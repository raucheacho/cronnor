package models

import (
	"database/sql"
	"time"
)

// Job represents a scheduled HTTP job
type Job struct {
	ID         int64          `json:"id"`
	Name       string         `json:"name"`
	CronExpr   string         `json:"cron_expr"`
	URL        string         `json:"url"`
	Method     string         `json:"method"`
	Payload    sql.NullString `json:"payload,omitempty"`
	IsActive   bool           `json:"is_active"`
	CreatedAt  time.Time      `json:"created_at"`
	LastRunAt  sql.NullTime   `json:"last_run_at,omitempty"`
	LastStatus sql.NullString `json:"last_status,omitempty"`
}

// JobLog represents an execution log entry
type JobLog struct {
	ID           int64          `json:"id"`
	JobID        int64          `json:"job_id"`
	Status       string         `json:"status"`
	HTTPCode     sql.NullInt64  `json:"http_code,omitempty"`
	DurationMs   sql.NullInt64  `json:"duration_ms,omitempty"`
	ResponseBody sql.NullString `json:"response_body,omitempty"`
	ErrorMessage sql.NullString `json:"error_message,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
}

// CreateJobParams represents parameters for creating a new job
type CreateJobParams struct {
	Name     string
	CronExpr string
	URL      string
	Method   string
	Payload  sql.NullString
}

// UpdateJobParams represents parameters for updating a job
type UpdateJobParams struct {
	ID       int64
	Name     string
	CronExpr string
	URL      string
	Method   string
	Payload  sql.NullString
}
