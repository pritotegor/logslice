package cmd

import (
	"testing"
	"time"
)

func TestValidateFormat_Valid(t *testing.T) {
	for _, f := range SupportedFormats {
		if err := ValidateFormat(f); err != nil {
			t.Errorf("expected format %q to be valid, got error: %v", f, err)
		}
	}
}

func TestValidateFormat_Invalid(t *testing.T) {
	err := ValidateFormat("xml")
	if err == nil {
		t.Error("expected error for unsupported format 'xml', got nil")
	}
}

func TestResolveTimeFormat_KnownAlias(t *testing.T) {
	layout := ResolveTimeFormat("rfc3339")
	if layout != time.RFC3339 {
		t.Errorf("expected RFC3339 layout, got %q", layout)
	}
}

func TestResolveTimeFormat_CaseInsensitive(t *testing.T) {
	layout := ResolveTimeFormat("RFC3339")
	if layout != time.RFC3339 {
		t.Errorf("expected RFC3339 layout for uppercase alias, got %q", layout)
	}
}

func TestResolveTimeFormat_UnknownAlias(t *testing.T) {
	custom := "2006/01/02"
	layout := ResolveTimeFormat(custom)
	if layout != custom {
		t.Errorf("expected raw layout %q to be returned as-is, got %q", custom, layout)
	}
}

func TestFormatTimestamp(t *testing.T) {
	ts := time.Date(2024, 6, 15, 10, 30, 0, 0, time.UTC)
	result := FormatTimestamp(ts, "date")
	if result != "2024-06-15" {
		t.Errorf("expected '2024-06-15', got %q", result)
	}
}

func TestParseTimeArg_Valid(t *testing.T) {
	_, err := ParseTimeArg("2024-06-15T10:30:00Z", "rfc3339")
	if err != nil {
		t.Errorf("unexpected error parsing valid time: %v", err)
	}
}

func TestParseTimeArg_Invalid(t *testing.T) {
	_, err := ParseTimeArg("not-a-time", "rfc3339")
	if err == nil {
		t.Error("expected error for invalid time string, got nil")
	}
}
