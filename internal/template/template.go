// Package template renders .env files from a Go text/template string,
// allowing users to define custom output formats for synced secrets.
package template

import (
	"bytes"
	"fmt"
	"os"
	"text/template"
)

// Renderer renders secrets using a user-supplied template.
type Renderer struct {
	tmpl *template.Template
}

// New parses the template string and returns a Renderer.
func New(tmplStr string) (*Renderer, error) {
	t, err := template.New("env").Funcs(template.FuncMap{
		"upper": func(s string) string {
			b := make([]byte, len(s))
			for i := range s {
				c := s[i]
				if c >= 'a' && c <= 'z' {
					c -= 32
				}
				b[i] = c
			}
			return string(b)
		},
	}).Parse(tmplStr)
	if err != nil {
		return nil, fmt.Errorf("template parse error: %w", err)
	}
	return &Renderer{tmpl: t}, nil
}

// NewFromFile reads a template file and returns a Renderer.
func NewFromFile(path string) (*Renderer, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading template file: %w", err)
	}
	return New(string(data))
}

// Render executes the template with the provided secrets map and returns the result.
func (r *Renderer) Render(secrets map[string]string) ([]byte, error) {
	var buf bytes.Buffer
	if err := r.tmpl.Execute(&buf, secrets); err != nil {
		return nil, fmt.Errorf("template execute error: %w", err)
	}
	return buf.Bytes(), nil
}

// RenderToFile executes the template and writes the output to path.
func (r *Renderer) RenderToFile(secrets map[string]string, path string) error {
	out, err := r.Render(secrets)
	if err != nil {
		return err
	}
	return os.WriteFile(path, out, 0600)
}
