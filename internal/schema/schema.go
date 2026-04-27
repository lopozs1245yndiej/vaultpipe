package schema

import (
	"errors"
	"fmt"
	"regexp"
)

// FieldType represents the expected type of a secret value.
type FieldType string

const (
	TypeString FieldType = "string"
	TypeInt    FieldType = "int"
	TypeBool   FieldType = "bool"
	TypeURL    FieldType = "url"
)

// Field defines a single key's schema constraints.
type Field struct {
	Key      string
	Type     FieldType
	Required bool
	Pattern  string // optional regex pattern
}

// Schema holds the set of field definitions to validate secrets against.
type Schema struct {
	fields map[string]Field
}

// ValidationError describes a single schema violation.
type ValidationError struct {
	Key     string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("schema: key %q: %s", e.Key, e.Message)
}

// New creates a Schema from a slice of Field definitions.
// Returns an error if any field has an empty key or invalid pattern.
func New(fields []Field) (*Schema, error) {
	m := make(map[string]Field, len(fields))
	for _, f := range fields {
		if f.Key == "" {
			return nil, errors.New("schema: field key must not be empty")
		}
		if f.Pattern != "" {
			if _, err := regexp.Compile(f.Pattern); err != nil {
				return nil, fmt.Errorf("schema: field %q has invalid pattern: %w", f.Key, err)
			}
		}
		m[f.Key] = f
	}
	return &Schema{fields: m}, nil
}

// Validate checks secrets against the schema and returns all violations.
func (s *Schema) Validate(secrets map[string]string) []ValidationError {
	var errs []ValidationError

	for key, field := range s.fields {
		val, exists := secrets[key]
		if !exists || val == "" {
			if field.Required {
				errs = append(errs, ValidationError{Key: key, Message: "required key is missing or empty"})
			}
			continue
		}
		if err := validateType(val, field.Type); err != nil {
			errs = append(errs, ValidationError{Key: key, Message: err.Error()})
		}
		if field.Pattern != "" {
			matched, _ := regexp.MatchString(field.Pattern, val)
			if !matched {
				errs = append(errs, ValidationError{Key: key, Message: fmt.Sprintf("value does not match pattern %q", field.Pattern)})
			}
		}
	}
	return errs
}

func validateType(val string, t FieldType) error {
	switch t {
	case TypeInt:
		if matched, _ := regexp.MatchString(`^-?\d+$`, val); !matched {
			return fmt.Errorf("expected int, got %q", val)
		}
	case TypeBool:
		if matched, _ := regexp.MatchString(`^(true|false|1|0)$`, val); !matched {
			return fmt.Errorf("expected bool, got %q", val)
		}
	case TypeURL:
		if matched, _ := regexp.MatchString(`^https?://`, val); !matched {
			return fmt.Errorf("expected URL starting with http:// or https://, got %q", val)
		}
	}
	return nil
}
