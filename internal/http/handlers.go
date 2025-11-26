package http

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rauche/cronnor/internal/models"
)

// handleDashboard shows the main dashboard
func (s *Server) handleDashboard(w http.ResponseWriter, r *http.Request) {
	jobs, err := s.repo.GetAllJobs()
	if err != nil {
		http.Error(w, "Failed to load jobs", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Jobs": jobs,
	}

	if err := s.templates.Render(w, "dashboard.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleJobsList returns the jobs list partial (for HTMX)
func (s *Server) handleJobsList(w http.ResponseWriter, r *http.Request) {
	jobs, err := s.repo.GetAllJobs()
	if err != nil {
		http.Error(w, "Failed to load jobs", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Jobs": jobs,
	}

	if err := s.templates.Render(w, "_job_list.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleJobDetail shows job details and execution history
func (s *Server) handleJobDetail(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
		return
	}

	job, err := s.repo.GetJob(id)
	if err != nil {
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	logs, err := s.repo.GetJobLogs(id, 50)
	if err != nil {
		http.Error(w, "Failed to load logs", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Job":  job,
		"Logs": logs,
	}

	if err := s.templates.Render(w, "job_detail.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleJobForm shows the new job form
func (s *Server) handleJobForm(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Job": nil,
	}

	if err := s.templates.Render(w, "job_form.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleJobEditForm shows the edit job form
func (s *Server) handleJobEditForm(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
		return
	}

	job, err := s.repo.GetJob(id)
	if err != nil {
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	data := map[string]interface{}{
		"Job": job,
	}

	if err := s.templates.Render(w, "job_form.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleCreateJob creates a new job
func (s *Server) handleCreateJob(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	payload := sql.NullString{}
	if p := r.FormValue("payload"); p != "" {
		payload.String = p
		payload.Valid = true
	}

	params := models.CreateJobParams{
		Name:     r.FormValue("name"),
		CronExpr: r.FormValue("cron_expr"),
		URL:      r.FormValue("url"),
		Method:   r.FormValue("method"),
		Payload:  payload,
	}

	id, err := s.repo.CreateJob(params)
	if err != nil {
		http.Error(w, "Failed to create job", http.StatusInternalServerError)
		return
	}

	// Load job into scheduler
	job, err := s.repo.GetJob(id)
	if err == nil {
		s.scheduler.AddJob(*job)
	}

	http.Redirect(w, r, "/jobs", http.StatusSeeOther)
}

// handleUpdateJob updates an existing job
func (s *Server) handleUpdateJob(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	payload := sql.NullString{}
	if p := r.FormValue("payload"); p != "" {
		payload.String = p
		payload.Valid = true
	}

	params := models.UpdateJobParams{
		ID:       id,
		Name:     r.FormValue("name"),
		CronExpr: r.FormValue("cron_expr"),
		URL:      r.FormValue("url"),
		Method:   r.FormValue("method"),
		Payload:  payload,
	}

	if err := s.repo.UpdateJob(params); err != nil {
		http.Error(w, "Failed to update job", http.StatusInternalServerError)
		return
	}

	// Reload job in scheduler
	s.scheduler.ReloadJob(id)

	http.Redirect(w, r, "/jobs/"+strconv.FormatInt(id, 10), http.StatusSeeOther)
}

// handleToggleJob toggles a job's active status
func (s *Server) handleToggleJob(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
		return
	}

	if err := s.repo.ToggleJob(id); err != nil {
		http.Error(w, "Failed to toggle job", http.StatusInternalServerError)
		return
	}

	// Reload job in scheduler
	s.scheduler.ReloadJob(id)

	// Return updated job list for HTMX
	s.handleJobsList(w, r)
}

// handleRunJob executes a job immediately
func (s *Server) handleRunJob(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
		return
	}

	if err := s.scheduler.ExecuteNow(id); err != nil {
		http.Error(w, "Failed to execute job", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Job execution started"))
}

// handleDeleteJob deletes a job
func (s *Server) handleDeleteJob(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
		return
	}

	// Remove from scheduler first
	s.scheduler.RemoveJob(id)

	if err := s.repo.DeleteJob(id); err != nil {
		http.Error(w, "Failed to delete job", http.StatusInternalServerError)
		return
	}

	// Return updated job list for HTMX
	s.handleJobsList(w, r)
}
