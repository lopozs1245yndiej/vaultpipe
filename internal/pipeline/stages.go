package pipeline

import (
	"context"
	"strings"
)

// UppercaseValues returns a Stage that converts all secret values to uppercase.
func UppercaseValues() Stage {
	return func(_ context.Context, secrets map[string]string) (map[string]string, error) {
		out := make(map[string]string, len(secrets))
		for k, v := range secrets {
			out[k] = strings.ToUpper(v)
		}
		return out, nil
	}
}

// PrefixKeys returns a Stage that prepends prefix to every key.
func PrefixKeys(prefix string) Stage {
	return func(_ context.Context, secrets map[string]string) (map[string]string, error) {
		out := make(map[string]string, len(secrets))
		for k, v := range secrets {
			out[prefix+k] = v
		}
		return out, nil
	}
}

// FilterKeys returns a Stage that retains only keys matching the given prefix.
func FilterKeys(prefix string) Stage {
	return func(_ context.Context, secrets map[string]string) (map[string]string, error) {
		out := make(map[string]string)
		for k, v := range secrets {
			if strings.HasPrefix(k, prefix) {
				out[k] = v
			}
		}
		return out, nil
	}
}

// TrimValueSpace returns a Stage that trims leading and trailing whitespace
// from every secret value.
func TrimValueSpace() Stage {
	return func(_ context.Context, secrets map[string]string) (map[string]string, error) {
		out := make(map[string]string, len(secrets))
		for k, v := range secrets {
			out[k] = strings.TrimSpace(v)
		}
		return out, nil
	}
}
