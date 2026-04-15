// Package diff provides utilities for comparing existing .env file contents
// with newly fetched secrets from Vault, reporting additions, removals, and changes.
package diff

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ChangeType represents the kind of change detected for a key.
type ChangeType string

const (
	Added   ChangeType = "added"
	Removed ChangeType = "removed"
	Changed ChangeType = "changed"
)

// Change represents a single secret key change.
type Change struct {
	Key  string
	Type ChangeType
}

// Compare reads an existing env file and compares it against incoming secrets.
// It returns a slice of Change entries describing what would be modified.
func Compare(envFilePath string, incoming map[string]string) ([]Change, error) {
	existing, err := parseEnvFile(envFilePath)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("diff: reading existing file: %w", err)
	}

	var changes []Change

	for k, v := range incoming {
		if oldVal, ok := existing[k]; !ok {
			changes = append(changes, Change{Key: k, Type: Added})
		} else if oldVal != v {
			changes = append(changes, Change{Key: k, Type: Changed})
		}
	}

	for k := range existing {
		if _, ok := incoming[k]; !ok {
			changes = append(changes, Change{Key: k, Type: Removed})
		}
	}

	return changes, nil
}

// parseEnvFile reads a KEY=VALUE formatted file into a map.
func parseEnvFile(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	result := make(map[string]string)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		result[parts[0]] = strings.Trim(parts[1], "\"")
	}
	return result, scanner.Err()
}
