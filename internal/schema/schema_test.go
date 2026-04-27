package schema_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/schema"
)

func TestNew_EmptyKey_ReturnsError(t *testing.T) {
	_, err := schema.New([]schema.Field{{Key: "", Type: schema.TypeString}})
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestNew_InvalidPattern_ReturnsError(t *testing.T) {
	_, err := schema.New([]schema.Field{{Key: "FOO", Pattern: "[invalid"}})
	if err == nil {
		t.Fatal("expected error for invalid regex pattern")
	}
}

func TestNew_Valid(t *testing.T) {
	_, err := schema.New([]schema.Field{
		{Key: "DB_HOST", Type: schema.TypeString, Required: true},
		{Key: "PORT", Type: schema.TypeInt},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidate_RequiredMissing(t *testing.T) {
	s, _ := schema.New([]schema.Field{
		{Key: "API_KEY", Type: schema.TypeString, Required: true},
	})
	errs := s.Validate(map[string]string{})
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
}

func TestValidate_InvalidInt(t *testing.T) {
	s, _ := schema.New([]schema.Field{
		{Key: "PORT", Type: schema.TypeInt},
	})
	errs := s.Validate(map[string]string{"PORT": "not-a-number"})
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
}

func TestValidate_InvalidBool(t *testing.T) {
	s, _ := schema.New([]schema.Field{
		{Key: "ENABLED", Type: schema.TypeBool},
	})
	errs := s.Validate(map[string]string{"ENABLED": "yes"})
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
}

func TestValidate_InvalidURL(t *testing.T) {
	s, _ := schema.New([]schema.Field{
		{Key: "ENDPOINT", Type: schema.TypeURL},
	})
	errs := s.Validate(map[string]string{"ENDPOINT": "ftp://bad"})
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
}

func TestValidate_PatternMismatch(t *testing.T) {
	s, _ := schema.New([]schema.Field{
		{Key: "REGION", Pattern: `^us-[a-z]+-\d+$`},
	})
	errs := s.Validate(map[string]string{"REGION": "eu-west-1"})
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
}

func TestValidate_AllValid(t *testing.T) {
	s, _ := schema.New([]schema.Field{
		{Key: "DB_URL", Type: schema.TypeURL, Required: true},
		{Key: "WORKERS", Type: schema.TypeInt},
		{Key: "DEBUG", Type: schema.TypeBool},
	})
	errs := s.Validate(map[string]string{
		"DB_URL":  "https://db.example.com",
		"WORKERS": "4",
		"DEBUG":   "false",
	})
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
}
