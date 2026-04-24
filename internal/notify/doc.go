// Package notify implements webhook-based notification for vaultpipe sync events.
//
// After a successful secret sync, callers can dispatch an Event to an HTTP
// endpoint of their choice. This is useful for triggering downstream
// processes (e.g. restarting a service, posting to Slack) whenever secrets
// are updated in the local environment.
//
// Basic usage:
//
//	 n, err := notify.New("https://hooks.example.com/vault", map[string]string{
//	     "Authorization": "Bearer " + token,
//	 })
//	 if err != nil { ... }
//
//	 err = n.Send(ctx, notify.Event{
//	     SecretPath:  "secret/data/myapp",
//	     KeysChanged: []string{"DB_PASSWORD"},
//	 })
package notify
