package http

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rauche/cronnor/internal/jobs"
	"github.com/rauche/cronnor/internal/storage"
)

// Server represents the HTTP server
type Server struct {
	router    *chi.Mux
	repo      *storage.Repository
	scheduler *jobs.Scheduler
	templates *TemplateRenderer
}

// NewServer creates a new HTTP server
func NewServer(repo *storage.Repository, scheduler *jobs.Scheduler) (*Server, error) {
	templates, err := NewTemplateRenderer("./web/templates")
	if err != nil {
		return nil, fmt.Errorf("failed to load templates: %w", err)
	}

	s := &Server{
		router:    chi.NewRouter(),
		repo:      repo,
		scheduler: scheduler,
		templates: templates,
	}

	s.setupRoutes()
	return s, nil
}

// setupRoutes configures all HTTP routes
func (s *Server) setupRoutes() {
	r := s.router

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	// Static files
	workDir, _ := filepath.Abs("./web/static")
	filesDir := http.Dir(workDir)
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(filesDir)))

	// Web routes
	r.Get("/", s.handleDashboard)
	r.Get("/jobs", s.handleDashboard)
	r.Get("/jobs/list", s.handleJobsList)            // API: Job list partial
	r.Get("/jobs/new", s.handleJobForm)              // New job form
	r.Post("/jobs", s.handleCreateJob)               // Create job
	r.Get("/jobs/{id}/edit", s.handleJobEditForm)    // Edit form (must be before /jobs/{id})
	r.Get("/jobs/{id}", s.handleJobDetail)           // Job details
	r.Post("/jobs/{id}", s.handleUpdateJob)          // Update job
	r.Post("/jobs/{id}/toggle", s.handleToggleJob)   // Toggle active
	r.Post("/jobs/{id}/run", s.handleRunJob)         // Run now
	r.Post("/jobs/{id}/delete", s.handleDeleteJob)   // Delete job (POST)
	r.Delete("/jobs/{id}", s.handleDeleteJob)        // Delete job (DELETE)
}

// Start starts the HTTP server
func (s *Server) Start(port string) error {
	addr := ":" + port
	log.Printf("Starting HTTP server on %s", addr)
	return http.ListenAndServe(addr, s.router)
}
