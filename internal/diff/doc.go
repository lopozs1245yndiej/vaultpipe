// Package diff compares an existing local .env file against a fresh set of
// secrets retrieved from HashiCorp Vault.
//
// It categorises every key into one of three change types:
//
//   - Added   – the key is present in Vault but absent from the local file.
//   - Removed – the key exists in the local file but is no longer in Vault.
//   - Changed – the key exists in both but its value has been updated in Vault.
//
// Usage:
//
//	changes, err := diff.Compare(".env", secretsMap)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, c := range changes {
//	    fmt.Printf("%s: %s\n", c.Type, c.Key)
//	}
package diff
