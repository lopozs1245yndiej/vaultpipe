package template

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNew_InvalidTemplate(t *testing.T) {
	_, err := New("{{ .Unclosed")
	if err == nil {
		t.Fatal("expected parse error, got nil")
	}
}

func TestRender_BasicTemplate(t *testing.T) {
	r, err := New("{{range $k,$v := .}}{{$k}}={{$v}}\n{{end}}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	secrets := map[string]string{"FOO": "bar"}
	out, err := r.Render(secrets)
	if err != nil {
		t.Fatalf("render error: %v", err)
	}
	if string(out) != "FOO=bar\n" {
		t.Errorf("unexpected output: %q", string(out))
	}
}

func TestRender_UpperFunc(t *testing.T) {
	r, err := New("{{range $k,$v := .}}{{upper $k}}={{$v}}\n{{end}}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out, err := r.Render(map[string]string{"foo": "baz"})
	if err != nil {
		t.Fatalf("render error: %v", err)
	}
	if string(out) != "FOO=baz\n" {
		t.Errorf("unexpected output: %q", string(out))
	}
}

func TestRenderToFile_WritesFile(t *testing.T) {
	r, err := New("{{range $k,$v := .}}{{$k}}={{$v}}\n{{end}}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tmp := filepath.Join(t.TempDir(), ".env")
	if err := r.RenderToFile(map[string]string{"KEY": "val"}, tmp); err != nil {
		t.Fatalf("render to file error: %v", err)
	}
	data, _ := os.ReadFile(tmp)
	if string(data) != "KEY=val\n" {
		t.Errorf("unexpected file content: %q", string(data))
	}
}

func TestNewFromFile_MissingFile(t *testing.T) {
	_, err := NewFromFile("/nonexistent/template.tmpl")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestNewFromFile_ValidFile(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "tmpl.txt")
	_ = os.WriteFile(tmp, []byte("{{range $k,$v := .}}{{$k}}={{$v}}\n{{end}}"), 0600)
	r, err := NewFromFile(tmp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out, _ := r.Render(map[string]string{"X": "1"})
	if string(out) != "X=1\n" {
		t.Errorf("unexpected output: %q", string(out))
	}
}
