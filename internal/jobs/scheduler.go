package jobs

import (
	"fmt"
	"log"
	"sync"

	"github.com/rauche/cronnor/internal/models"
	"github.com/rauche/cronnor/internal/storage"
	"github.com/robfig/cron/v3"
)

// Scheduler manages cron jobs
type Scheduler struct {
	cron     *cron.Cron
	repo     *storage.Repository
	executor *Executor
	entries  map[int64]cron.EntryID // job ID -> cron entry ID
	mu       sync.RWMutex
}

// NewScheduler creates a new scheduler
func NewScheduler(repo *storage.Repository) *Scheduler {
	return &Scheduler{
		cron:     cron.New(cron.WithSeconds()),
		repo:     repo,
		executor: NewExecutor(repo),
		entries:  make(map[int64]cron.EntryID),
	}
}

// Start starts the scheduler and loads active jobs
func (s *Scheduler) Start() error {
	// Load active jobs from database
	jobs, err := s.repo.GetActiveJobs()
	if err != nil {
		return fmt.Errorf("failed to load active jobs: %w", err)
	}

	// Schedule each job
	for _, job := range jobs {
		if err := s.AddJob(job); err != nil {
			log.Printf("Warning: failed to schedule job %d (%s): %v", job.ID, job.Name, err)
		}
	}

	// Start cron scheduler
	s.cron.Start()
	log.Printf("Scheduler started with %d active jobs", len(jobs))

	return nil
}

// Stop stops the scheduler
func (s *Scheduler) Stop() {
	ctx := s.cron.Stop()
	<-ctx.Done()
	log.Println("Scheduler stopped")
}

// AddJob adds a job to the scheduler
func (s *Scheduler) AddJob(job models.Job) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Remove existing entry if present
	if entryID, exists := s.entries[job.ID]; exists {
		s.cron.Remove(entryID)
		delete(s.entries, job.ID)
	}

	// Only schedule if active
	if !job.IsActive {
		return nil
	}

	// Add job to cron
	entryID, err := s.cron.AddFunc(job.CronExpr, func() {
		s.executeJob(job)
	})
	if err != nil {
		return fmt.Errorf("failed to add cron job: %w", err)
	}

	s.entries[job.ID] = entryID
	log.Printf("Scheduled job %d (%s) with cron expression: %s", job.ID, job.Name, job.CronExpr)

	return nil
}

// RemoveJob removes a job from the scheduler
func (s *Scheduler) RemoveJob(jobID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if entryID, exists := s.entries[jobID]; exists {
		s.cron.Remove(entryID)
		delete(s.entries, jobID)
		log.Printf("Removed job %d from scheduler", jobID)
	}
}

// ExecuteNow executes a job immediately (bypassing the cron schedule)
func (s *Scheduler) ExecuteNow(jobID int64) error {
	job, err := s.repo.GetJob(jobID)
	if err != nil {
		return fmt.Errorf("failed to get job: %w", err)
	}

	go s.executeJob(*job)
	return nil
}

// executeJob executes a job
func (s *Scheduler) executeJob(job models.Job) {
	log.Printf("Executing job %d (%s): %s %s", job.ID, job.Name, job.Method, job.URL)

	if err := s.executor.Execute(job); err != nil {
		log.Printf("Job %d (%s) execution failed: %v", job.ID, job.Name, err)
	} else {
		log.Printf("Job %d (%s) executed successfully", job.ID, job.Name)
	}
}

// ReloadJob reloads a job (e.g., after update)
func (s *Scheduler) ReloadJob(jobID int64) error {
	job, err := s.repo.GetJob(jobID)
	if err != nil {
		return fmt.Errorf("failed to get job: %w", err)
	}

	return s.AddJob(*job)
}
