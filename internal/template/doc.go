// Package template provides a Renderer that generates .env file output
// from a user-supplied Go text/template string.
//
// Templates receive a map[string]string of secret key/value pairs as their
// data context. A built-in "upper" function is available to convert keys to
// uppercase.
//
// Example template:
//
//	{{range $k,$v := .}}{{upper $k}}={{$v}}
//	{{end}}
//
// Use New to parse an inline template string, or NewFromFile to load one
// from disk. Call Render to obtain the rendered bytes, or RenderToFile to
// write directly to a path (mode 0600).
package template
