// Package ratelimit provides a simple token-bucket rate limiter for
// controlling the frequency of outbound Vault API requests made by vaultpipe.
//
// Usage:
//
//	l, err := ratelimit.New(10, time.Second)
//	if err != nil {
//		log.Fatal(err)
//	}
//	// Before each Vault call:
//	if err := l.Wait(ctx); err != nil {
//		return err
//	}
//
// Alternatively, use NewFromOptions for struct-based configuration:
//
//	l, err := ratelimit.NewFromOptions(ratelimit.Options{
//		Rate:   20,
//		Window: time.Minute,
//	})
package ratelimit
