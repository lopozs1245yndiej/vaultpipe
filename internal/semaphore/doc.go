// Package semaphore implements a simple counting semaphore that can be used to
// bound the number of concurrent operations running at any given time.
//
// # Usage
//
//	s, err := semaphore.New(5, 2*time.Second)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	for _, path := range secretPaths {
//		if err := s.Acquire(ctx); err != nil {
//			return err
//		}
//		go func(p string) {
//			defer s.Release()
//			// fetch secret at p ...
//		}(path)
//	}
//
// The timeout passed to New applies per Acquire call. A zero timeout means the
// caller blocks until a slot is free or the context is cancelled.
package semaphore
