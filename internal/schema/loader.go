package schema

import (
	"encoding/json"
	"fmt"
	"os"
)

// FileField is the JSON-serialisable representation of a Field.
type FileField struct {
	Key      string `json:"key"`
	Type     string `json:"type"`
	Required bool   `json:"required"`
	Pattern  string `json:"pattern,omitempty"`
}

// LoadFromFile reads a JSON schema definition file and returns a Schema.
// The file must contain a JSON array of field objects.
func LoadFromFile(path string) (*Schema, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("schema: read file %q: %w", path, err)
	}

	var raw []FileField
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("schema: parse file %q: %w", path, err)
	}

	fields := make([]Field, 0, len(raw))
	for _, r := range raw {
		ft := FieldType(r.Type)
		if ft == "" {
			ft = TypeString
		}
		fields = append(fields, Field{
			Key:      r.Key,
			Type:     ft,
			Required: r.Required,
			Pattern:  r.Pattern,
		})
	}
	return New(fields)
}

// MustLoad is like LoadFromFile but panics on error.
func MustLoad(path string) *Schema {
	s, err := LoadFromFile(path)
	if err != nil {
		panic(err)
	}
	return s
}
