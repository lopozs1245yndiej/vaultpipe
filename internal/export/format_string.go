package export

import "fmt"

// ParseFormat converts a raw string to a Format, returning an error if unrecognised.
func ParseFormat(s string) (Format, error) {
	switch Format(s) {
	case FormatEnv:
		return FormatEnv, nil
	case FormatJSON:
		return FormatJSON, nil
	case FormatYAML:
		return FormatYAML, nil
	default:
		return "", fmt.Errorf("unknown format %q: must be one of env, json, yaml", s)
	}
}

// String returns the string representation of the Format.
func (f Format) String() string {
	return string(f)
}
