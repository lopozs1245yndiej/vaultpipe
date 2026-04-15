package envwriter

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// Writer handles writing secrets to .env files.
type Writer struct {
	filePath  string
	namespace string
}

// NewWriter creates a new Writer for the given file path.
// If namespace is non-empty, only keys matching the namespace prefix are written.
func NewWriter(filePath, namespace string) *Writer {
	return &Writer{
		filePath:  filePath,
		namespace: namespace,
	}
}

// Write serialises the provided secrets map into .env format and writes it
// to the configured file path, creating or truncating the file as needed.
func (w *Writer) Write(secrets map[string]string) error {
	filtered := w.filter(secrets)

	lines := make([]string, 0, len(filtered))
	for k, v := range filtered {
		lines = append(lines, fmt.Sprintf("%s=%s", k, escapeValue(v)))
	}
	sort.Strings(lines)

	content := strings.Join(lines, "\n")
	if len(lines) > 0 {
		content += "\n"
	}

	return os.WriteFile(w.filePath, []byte(content), 0o600)
}

// filter returns only the entries whose keys match the namespace prefix.
// If namespace is empty all entries are returned unchanged.
func (w *Writer) filter(secrets map[string]string) map[string]string {
	if w.namespace == "" {
		return secrets
	}

	prefix := strings.ToUpper(w.namespace) + "_"
	out := make(map[string]string)
	for k, v := range secrets {
		if strings.HasPrefix(strings.ToUpper(k), prefix) {
			out[k] = v
		}
	}
	return out
}

// escapeValue wraps values that contain spaces or special characters in
// double-quotes and escapes any embedded double-quotes.
func escapeValue(v string) string {
	if strings.ContainsAny(v, " \t\n\r#") {
		v = strings.ReplaceAll(v, `"`, `\"`)
		return `"` + v + `"`
	}
	return v
}
