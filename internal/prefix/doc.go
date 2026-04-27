// Package prefix provides key-prefix management for secret maps.
//
// It supports three operations:
//
//   - Add: prepend a namespace prefix to every key in a secrets map.
//   - Strip: remove a known prefix from matching keys, passing others through.
//   - Replace: atomically swap one prefix for another on matching keys.
//
// Example usage:
//
//	p, err := prefix.New("APP")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// {"DB_HOST": "localhost"} → {"APP_DB_HOST": "localhost"}
//	prefixed := p.Add(secrets)
//
//	// {"APP_DB_HOST": "localhost"} → {"DB_HOST": "localhost"}
//	stripped := p.Strip(prefixed)
package prefix
