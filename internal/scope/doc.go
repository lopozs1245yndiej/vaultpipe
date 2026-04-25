// Package scope provides namespace-scoped filtering and key remapping for
// secrets loaded from HashiCorp Vault.
//
// # Overview
//
// When secrets are stored under a namespaced path in Vault (e.g. APP_DB_HOST,
// APP_DB_PORT) it is often desirable to strip the namespace prefix before
// writing the keys to a .env file so that consuming applications receive
// plain keys (DB_HOST, DB_PORT).
//
// # Usage
//
//	s, err := scope.New("APP_",
//	    scope.WithReplacePrefix("SVC_"),
//	    scope.WithCaseSensitive(false),
//	)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	filtered := s.Apply(rawSecrets)
package scope
