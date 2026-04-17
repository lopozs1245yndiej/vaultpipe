package watch

import "context"

// Syncable is the interface required by Watcher to perform a sync operation.
type Syncable interface {
	Run(ctx context.Context) error
}
