package lock

import (
	"fmt"
	"os"
	"path/filepath"
)

// Options configures lock file creation.
type Options struct {
	// Dir is the directory where the lock file is created.
	// Defaults to os.TempDir() if empty.
	Dir string

	// Filename is the name of the lock file.
	// Defaults to "vaultpipe.lock" if empty.
	Filename string
}

// NewFromOptions creates a Locker using the provided options.
// Returns an error if the resolved directory cannot be determined.
func NewFromOptions(opts Options) (*Locker, error) {
	dir := opts.Dir
	if dir == "" {
		dir = os.TempDir()
	}

	name := opts.Filename
	if name == "" {
		name = "vaultpipe.lock"
	}

	abs, err := filepath.Abs(dir)
	if err != nil {
		return nil, fmt.Errorf("lock: resolve dir %q: %w", dir, err)
	}

	return New(filepath.Join(abs, name)), nil
}
