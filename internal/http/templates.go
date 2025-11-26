package http

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"github.com/robfig/cron/v3"
)

// TemplateRenderer handles template rendering
type TemplateRenderer struct {
	templates *template.Template
}

// NewTemplateRenderer creates a new template renderer
func NewTemplateRenderer(templatesDir string) (*TemplateRenderer, error) {
	funcMap := template.FuncMap{
		"formatTime":  formatTime,
		"statusClass": statusClass,
		"nextRun":     nextRun,
		"eq":          func(a, b string) bool { return a == b },
	}

	tmpl := template.New("").Funcs(funcMap)
	
	// Parse all templates
	_, err := tmpl.ParseGlob(filepath.Join(templatesDir, "*.html"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}

	return &TemplateRenderer{templates: tmpl}, nil
}

// Render renders a template with data
func (tr *TemplateRenderer) Render(w http.ResponseWriter, name string, data interface{}) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := tr.templates.ExecuteTemplate(w, name, data); err != nil {
		return fmt.Errorf("failed to execute template %s: %w", name, err)
	}

	return nil
}

// Template helper functions

func formatTime(t time.Time) string {
	if t.IsZero() {
		return "Never"
	}
	return t.Format("2006-01-02 15:04:05")
}

func statusClass(status string) string {
	switch status {
	case "SUCCESS":
		return "status-success"
	case "FAILED":
		return "status-failed"
	case "ERROR":
		return "status-error"
	default:
		return "status-pending"
	}
}

func nextRun(cronExpr string) string {
	parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	schedule, err := parser.Parse(cronExpr)
	if err != nil {
		return "Invalid cron expression"
	}

	next := schedule.Next(time.Now())
	return formatTime(next)
}
