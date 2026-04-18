package export

import (
	"encoding/json"
	"fmt"
	"os"\n	"sort"
	"strings"
)

// Format represents an export output format.
type Format string

const (
	FormatEnv  Format = "env"
	FormatJSON Format = "json"
	FormatYAML Format = "yaml"
)

// Exporter writes secrets to a file in a specified format.
type Exporter struct {
	format Format
}

// New creates a new Exporter for the given format.
func New(format Format) (*Exporter, error) {
	switch format {
	case FormatEnv, FormatJSON, FormatYAML:
		return &Exporter{format: format}, nil
	default:
		return nil, fmt.Errorf("unsupported export format: %q", format)
	}
}

// Export writes secrets to the given file path.
func (e *Exporter) Export(secrets map[string]string, path string) error {
	var data []byte
	var err error

	switch e.format {
	case FormatEnv:
		data = []byte(toEnv(secrets))
	case FormatJSON:
		data, err = toJSON(secrets)
	case FormatYAML:
		data = []byte(toYAML(secrets))
	}
	if err != nil {
		return fmt.Errorf("export marshal: %w", err)
	}
	return os.WriteFile(path, data, 0600)
}

func toEnv(secrets map[string]string) string {
	keys := sortedKeys(secrets)
	var sb strings.Builder
	for _, k := range keys {
		fmt.Fprintf(&sb, "%s=%s\n", k, secrets[k])
	}
	return sb.String()
}

func toJSON(secrets map[string]string) ([]byte, error) {
	return json.MarshalIndent(secrets, "", "  ")
}

func toYAML(secrets map[string]string) string {
	keys := sortedKeys(secrets)
	var sb strings.Builder
	for _, k := range keys {
		fmt.Fprintf(&sb, "%s: %q\n", k, secrets[k])
	}
	return sb.String()
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
