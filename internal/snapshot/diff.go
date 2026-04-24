package snapshot

// DiffResult holds the changes between two snapshots.
type DiffResult struct {
	Added   map[string]string
	Removed map[string]string
	Changed map[string][2]string // key -> [old, new]
}

// Diff compares two snapshots and returns the differences.
func Diff(old, new *Snapshot) DiffResult {
	result := DiffResult{
		Added:   make(map[string]string),
		Removed: make(map[string]string),
		Changed: make(map[string][2]string),
	}

	for k, newVal := range new.Secrets {
		oldVal, exists := old.Secrets[k]
		if !exists {
			result.Added[k] = newVal
		} else if oldVal != newVal {
			result.Changed[k] = [2]string{oldVal, newVal}
		}
	}

	for k, oldVal := range old.Secrets {
		if _, exists := new.Secrets[k]; !exists {
			result.Removed[k] = oldVal
		}
	}

	return result
}

// HasChanges returns true if there are any differences.
func (d DiffResult) HasChanges() bool {
	return len(d.Added) > 0 || len(d.Removed) > 0 || len(d.Changed) > 0
}
