package http

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
)

// TemplateRenderer handles template rendering
type TemplateRenderer struct {
	templates map[string]*template.Template
}

// NewTemplateRenderer creates a new template renderer
func NewTemplateRenderer(templatesDir string) (*TemplateRenderer, error) {
	funcMap := template.FuncMap{
		"formatTime":  formatTime,
		"statusClass": statusClass,
		"nextRun":     nextRun,
		"eq":          func(a, b string) bool { return a == b },
	}

	// 1. Identify files
	files, err := filepath.Glob(filepath.Join(templatesDir, "*.html"))
	if err != nil {
		return nil, fmt.Errorf("failed to list templates: %w", err)
	}

	var baseFiles []string
	var pageFiles []string

	for _, file := range files {
		name := filepath.Base(file)
		if name == "layout.html" || strings.HasPrefix(name, "_") {
			baseFiles = append(baseFiles, file)
		} else {
			pageFiles = append(pageFiles, file)
		}
	}

	// 2. Parse base templates (layout + partials)
	baseTmpl := template.New("base").Funcs(funcMap)
	if len(baseFiles) > 0 {
		baseTmpl, err = baseTmpl.ParseFiles(baseFiles...)
		if err != nil {
			return nil, fmt.Errorf("failed to parse base templates: %w", err)
		}
	}

	// 3. Create isolated template instances for each page
	templates := make(map[string]*template.Template)

	// Add partials/layout directly to map (for rendering partials like _job_list.html)
	for _, file := range baseFiles {
		name := filepath.Base(file)
		templates[name] = baseTmpl
	}

	// Parse each page into a clone of the base
	for _, file := range pageFiles {
		name := filepath.Base(file)
		// Clone the base template
		tmpl, err := baseTmpl.Clone()
		if err != nil {
			return nil, fmt.Errorf("failed to clone base template for %s: %w", name, err)
		}

		// Parse the page file into the clone
		_, err = tmpl.ParseFiles(file)
		if err != nil {
			return nil, fmt.Errorf("failed to parse page template %s: %w", name, err)
		}

		templates[name] = tmpl
	}

	return &TemplateRenderer{templates: templates}, nil
}

// Render renders a template with data
func (tr *TemplateRenderer) Render(w http.ResponseWriter, name string, data interface{}) error {
	tmpl, ok := tr.templates[name]
	if !ok {
		return fmt.Errorf("template %s not found", name)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// For pages (which use layout), we usually execute the file itself (which calls layout)
	// But since we parsed the file into the set, executing "name" works.
	// However, if "name" is a partial (like _job_list.html), it's also in the map.

	if err := tmpl.ExecuteTemplate(w, name, data); err != nil {
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
