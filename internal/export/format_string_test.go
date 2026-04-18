package export

import "testing"

func TestParseFormat_Valid(t *testing.T) {
	cases := []struct {
		input    string
		expected Format
	}{
		{"env", FormatEnv},
		{"json", FormatJSON},
		{"yaml", FormatYAML},
	}
	for _, tc := range cases {
		f, err := ParseFormat(tc.input)
		if err != nil {
			t.Errorf("unexpected error for %q: %v", tc.input, err)
		}
		if f != tc.expected {
			t.Errorf("expected %v, got %v", tc.expected, f)
		}
	}
}

func TestParseFormat_Invalid(t *testing.T) {
	_, err := ParseFormat("xml")
	if err == nil {
		t.Fatal("expected error for invalid format")
	}
}

func TestFormat_String(t *testing.T) {
	if FormatJSON.String() != "json" {
		t.Errorf("expected 'json', got %q", FormatJSON.String())
	}
}
