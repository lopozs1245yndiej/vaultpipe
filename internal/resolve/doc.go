// Package resolve implements variable interpolation for secret maps.
//
// It supports ${KEY} syntax, allowing one secret value to reference
// another key in the same map. This is useful when constructing
// compound values such as connection strings from individual parts.
//
// Example:
//
//	secrets := map[string]string{
//	    "DB_HOST": "localhost",
//	    "DB_PORT": "5432",
//	    "DB_DSN":  "postgres://${DB_HOST}:${DB_PORT}/mydb",
//	}
//
//	r, _ := resolve.New(10)
//	resolved, err := r.Resolve(secrets)
//	// resolved["DB_DSN"] == "postgres://localhost:5432/mydb"
//
// Circular references are detected via a configurable maximum depth.
// References to undefined keys produce an error.
package resolve
