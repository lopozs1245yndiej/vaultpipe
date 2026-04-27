// Package batch provides concurrent processing of multiple secret maps.
//
// A Batcher accepts a slice of secret maps and a ProcessFunc, then fans the
// work out across a bounded pool of goroutines. Each item is processed
// independently; partial failures are captured per-item rather than aborting
// the entire run.
//
// Basic usage:
//
//	b, err := batch.New(4, func(ctx context.Context, secrets map[string]string) error {
//		// write secrets to a destination
//		return writer.Write(secrets)
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	results, err := b.Run(ctx, allSecretMaps)
//	for _, r := range batch.Errors(results) {
//		log.Printf("item %d failed: %v", r.Index, r.Err)
//	}
package batch
