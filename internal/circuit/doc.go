// Package circuit provides a circuit breaker implementation for use with
// external service calls such as HashiCorp Vault API requests.
//
// The circuit breaker tracks consecutive failures and transitions through
// three states:
//
//   - Closed: normal operation; calls are allowed through.
//   - Open: failure threshold exceeded; calls are rejected with ErrOpen
//     until the reset timeout elapses.
//   - Half-Open: the reset timeout has elapsed; one call is allowed through
//     to probe recovery. A success closes the circuit; a failure reopens it.
//
// Example usage:
//
//	b, err := circuit.New(5, 30*time.Second)
//	if err != nil { ... }
//
//	if err := b.Allow(); err != nil {
//		// circuit is open, skip the call
//	}
//	if err := doVaultCall(); err != nil {
//		b.RecordFailure()
//	} else {
//		b.RecordSuccess()
//	}
package circuit
