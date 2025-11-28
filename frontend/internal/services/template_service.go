package services

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"time"
)

type TemplateService struct {
	templates *template.Template
}

// NewTemplateService creates a new template service with helper functions
func NewTemplateService(templatesDir string) (*TemplateService, error) {
	// Create a function map with helper functions for templates
	funcMap := template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"sub": func(a, b int) int {
			return a - b
		},
		"mul": func(a, b int) int {
			return a * b
		},
		"div": func(a, b int) int {
			if b == 0 {
				return 0
			}
			return a / b
		},
		"eq": func(a, b interface{}) bool {
			return a == b
		},
		"ne": func(a, b interface{}) bool {
			return a != b
		},
		"lt": func(a, b int) bool {
			return a < b
		},
		"le": func(a, b int) bool {
			return a <= b
		},
		"gt": func(a, b int) bool {
			return a > b
		},
		"ge": func(a, b int) bool {
			return a >= b
		},
		"range_until": func(max int) []int {
			result := make([]int, max)
			for i := 0; i < max; i++ {
				result[i] = i + 1
			}
			return result
		},
		"printf": func(format string, args ...interface{}) string {
			return fmt.Sprintf(format, args...)
		},
		// ✅ ADD: Custom function to handle TimeAgo for notifications
		"timeAgo": func(createdAt time.Time) string {
			now := time.Now()
			diff := now.Sub(createdAt)

			if diff < time.Minute {
				return "just now"
			}
			if diff < time.Hour {
				minutes := int(diff.Minutes())
				if minutes == 1 {
					return "1 minute ago"
				}
				return fmt.Sprintf("%d minutes ago", minutes)
			}
			if diff < 24*time.Hour {
				hours := int(diff.Hours())
				if hours == 1 {
					return "1 hour ago"
				}
				return fmt.Sprintf("%d hours ago", hours)
			}
			if diff < 7*24*time.Hour {
				days := int(diff.Hours() / 24)
				if days == 1 {
					return "1 day ago"
				}
				return fmt.Sprintf("%d days ago", days)
			}
			return createdAt.Format("Jan 2, 2006")
		},
		// ✅ ADD: Custom function to check if notification is unread
		"isUnread": func(isRead bool) bool {
			return !isRead
		},
	}

	// Parse all templates in the directory with the function map
	tmpl, err := template.New("").Funcs(funcMap).ParseGlob(filepath.Join(templatesDir, "*.html"))
	if err != nil {
		return nil, err
	}

	return &TemplateService{
		templates: tmpl,
	}, nil
}

// Render executes a template with the given data
func (ts *TemplateService) Render(w http.ResponseWriter, templateName string, data interface{}) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return ts.templates.ExecuteTemplate(w, templateName, data)
}
