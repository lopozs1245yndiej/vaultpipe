package rollback

import "fmt"

// Options configures a Rollbacker.
type Options struct {
	// BackupDir is the directory where *.env.bak files are stored.
	BackupDir string

	// MaxRestoreAge, if non-zero, causes Restore to reject backups older
	// than MaxRestoreAge. Zero means no age limit.
	MaxRestoreAge int // seconds
}

// NewFromOptions constructs a Rollbacker from an Options struct.
func NewFromOptions(opts Options) (*Rollbacker, error) {
	if opts.BackupDir == "" {
		return nil, fmt.Errorf("rollback: BackupDir is required")
	}
	r, err := New(opts.BackupDir)
	if err != nil {
		return nil, err
	}
	return r, nil
}
