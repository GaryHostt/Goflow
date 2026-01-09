package utils

import (
	"regexp"
	"strings"

	"github.com/tidwall/gjson"
)

// TemplateEngine handles data mapping from trigger to action
// Allows users to use dynamic values like {{user.name}} in their workflows
type TemplateEngine struct {
	templatePattern *regexp.Regexp
}

// NewTemplateEngine creates a new template engine
func NewTemplateEngine() *TemplateEngine {
	return &TemplateEngine{
		// Matches {{path.to.value}} or {{path}}
		templatePattern: regexp.MustCompile(`\{\{([^}]+)\}\}`),
	}
}

// Render replaces template variables with actual values from JSON data
func (te *TemplateEngine) Render(template string, data string) string {
	return te.templatePattern.ReplaceAllStringFunc(template, func(match string) string {
		// Extract the path from {{path}}
		path := strings.TrimSpace(match[2 : len(match)-2])
		
		// Use gjson to extract value from JSON
		result := gjson.Get(data, path)
		
		if !result.Exists() {
			// Path not found, keep original
			return match
		}
		
		return result.String()
	})
}

// RenderMap processes an entire config map with templates
func (te *TemplateEngine) RenderMap(config map[string]interface{}, data string) map[string]interface{} {
	rendered := make(map[string]interface{})
	
	for key, value := range config {
		switch v := value.(type) {
		case string:
			// Replace templates in string values
			rendered[key] = te.Render(v, data)
		case map[string]interface{}:
			// Recursively process nested maps
			rendered[key] = te.RenderMap(v, data)
		default:
			// Keep non-string values as-is
			rendered[key] = value
		}
	}
	
	return rendered
}

// ExtractValue is a helper to extract a specific value from JSON
func ExtractValue(data string, path string) string {
	result := gjson.Get(data, path)
	if !result.Exists() {
		return ""
	}
	return result.String()
}

// ValidateTemplate checks if a template string is valid
func (te *TemplateEngine) ValidateTemplate(template string) []string {
	var paths []string
	
	matches := te.templatePattern.FindAllStringSubmatch(template, -1)
	for _, match := range matches {
		if len(match) > 1 {
			path := strings.TrimSpace(match[1])
			paths = append(paths, path)
		}
	}
	
	return paths
}

// Example usage:
// template := "Hello {{user.name}}, your email is {{user.email}}"
// data := `{"user": {"name": "Alex", "email": "alex@example.com"}}`
// result := engine.Render(template, data)
// Output: "Hello Alex, your email is alex@example.com"

