package cmd

import (
	"strings"
	"testing"
)

func TestHighlightLevel_ColorDisabled(t *testing.T) {
	result := HighlightLevel("error", false)
	if result != "error" {
		t.Errorf("expected 'error', got %q", result)
	}
}

func TestHighlightLevel_KnownLevels(t *testing.T) {
	cases := []struct {
		input    string
		expColor string
	}{
		{"error", ColorRed},
		{"warn", ColorYellow},
		{"info", ColorGreen},
		{"debug", ColorCyan},
		{"ERR", ColorRed},
		{"WARNING", ColorYellow},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			result := HighlightLevel(tc.input, true)
			if !strings.Contains(result, tc.expColor) {
				t.Errorf("expected color %q in result %q", tc.expColor, result)
			}
			if !strings.Contains(result, ColorReset) {
				t.Errorf("expected reset code in result %q", result)
			}
		})
	}
}

func TestHighlightLevel_UnknownLevel(t *testing.T) {
	result := HighlightLevel("trace", true)
	if result != "trace" {
		t.Errorf("expected unmodified 'trace', got %q", result)
	}
}

func TestHighlightField_ColorEnabled(t *testing.T) {
	result := HighlightField("timestamp", true)
	if !strings.Contains(result, ColorCyan) {
		t.Errorf("expected cyan in field highlight, got %q", result)
	}
}

func TestHighlightField_ColorDisabled(t *testing.T) {
	result := HighlightField("timestamp", false)
	if result != "timestamp" {
		t.Errorf("expected 'timestamp', got %q", result)
	}
}

func TestHighlightMatch_ReplacesAll(t *testing.T) {
	result := HighlightMatch("foo bar foo", "foo", true)
	count := strings.Count(result, ColorYellow)
	if count != 2 {
		t.Errorf("expected 2 highlighted matches, got %d in %q", count, result)
	}
}

func TestHighlightMatch_EmptyMatch(t *testing.T) {
	result := HighlightMatch("some text", "", true)
	if result != "some text" {
		t.Errorf("expected unchanged text, got %q", result)
	}
}
