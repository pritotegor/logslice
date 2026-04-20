package cmd

import (
	"testing"
	"time"
)

func TestParseLine_NonJSON(t *testing.T) {
	pl, err := ParseLine("plain text log line", "time")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pl.IsJSON {
		t.Error("expected IsJSON=false for plain text")
	}
	if pl.Timestamp != nil {
		t.Error("expected nil timestamp for plain text")
	}
}

func TestParseLine_ValidJSON_WithTimestamp(t *testing.T) {
	raw := `{"time":"2024-03-15T10:00:00Z","level":"info","msg":"started"}`
	pl, err := ParseLine(raw, "time")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !pl.IsJSON {
		t.Error("expected IsJSON=true")
	}
	if pl.Timestamp == nil {
		t.Fatal("expected non-nil timestamp")
	}
	expected := time.Date(2024, 3, 15, 10, 0, 0, 0, time.UTC)
	if !pl.Timestamp.Equal(expected) {
		t.Errorf("got %v, want %v", pl.Timestamp, expected)
	}
}

func TestParseLine_ValidJSON_NoTimestampField(t *testing.T) {
	raw := `{"level":"warn","msg":"disk full"}`
	pl, err := ParseLine(raw, "time")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pl.Timestamp != nil {
		t.Error("expected nil timestamp when field absent")
	}
}

func TestParseLine_InvalidJSON(t *testing.T) {
	_, err := ParseLine(`{not valid json}`, "time")
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestParseTimestamp_MultipleFormats(t *testing.T) {
	cases := []struct {
		input string
	}{
		{"2024-03-15T10:00:00Z"},
		{"2024-03-15T10:00:00.000Z"},
		{"2024-03-15 10:00:00"},
		{"2024/03/15 10:00:00"},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			_, err := ParseTimestamp(tc.input)
			if err != nil {
				t.Errorf("failed to parse %q: %v", tc.input, err)
			}
		})
	}
}

func TestParseTimestamp_Invalid(t *testing.T) {
	_, err := ParseTimestamp("not-a-date")
	if err == nil {
		t.Error("expected error for invalid timestamp")
	}
}
