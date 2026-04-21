package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestDedupeLines_NoDuplicates(t *testing.T) {
	input := "line1\nline2\nline3\n"
	var out bytes.Buffer
	err := dedupeLines(strings.NewReader(input), &out, "", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := strings.TrimSpace(out.String())
	lines := strings.Split(got, "\n")
	if len(lines) != 3 {
		t.Errorf("expected 3 lines, got %d", len(lines))
	}
}

func TestDedupeLines_WithDuplicates(t *testing.T) {
	input := "line1\nline2\nline1\nline3\nline2\n"
	var out bytes.Buffer
	err := dedupeLines(strings.NewReader(input), &out, "", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := strings.TrimSpace(out.String())
	lines := strings.Split(got, "\n")
	if len(lines) != 3 {
		t.Errorf("expected 3 unique lines, got %d: %v", len(lines), lines)
	}
}

func TestDedupeLines_ByField(t *testing.T) {
	input := `{"msg":"hello","id":1}
{"msg":"world","id":2}
{"msg":"again","id":1}
`
	var out bytes.Buffer
	err := dedupeLines(strings.NewReader(input), &out, "id", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := strings.TrimSpace(out.String())
	lines := strings.Split(got, "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 lines after field dedupe, got %d: %v", len(lines), lines)
	}
}

func TestDedupeLines_WindowEvicts(t *testing.T) {
	// With window=2, line1 should be evicted after 2 more lines, so it appears again
	input := "line1\nline2\nline3\nline1\n"
	var out bytes.Buffer
	err := dedupeLines(strings.NewReader(input), &out, "", 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := strings.TrimSpace(out.String())
	lines := strings.Split(got, "\n")
	if len(lines) != 4 {
		t.Errorf("expected 4 lines with window eviction, got %d: %v", len(lines), lines)
	}
}

func TestDedupeLines_SkipsEmptyLines(t *testing.T) {
	input := "line1\n\n   \nline2\n"
	var out bytes.Buffer
	err := dedupeLines(strings.NewReader(input), &out, "", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := strings.TrimSpace(out.String())
	lines := strings.Split(got, "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 non-empty lines, got %d", len(lines))
	}
}

func TestDedupeKey_FullLine(t *testing.T) {
	k1 := dedupeKey("hello world", "")
	k2 := dedupeKey("hello world", "")
	k3 := dedupeKey("different", "")
	if k1 != k2 {
		t.Error("same line should produce same key")
	}
	if k1 == k3 {
		t.Error("different lines should produce different keys")
	}
}

func TestDedupeKey_ByField_FallsBackOnMissingField(t *testing.T) {
	line := `{"msg":"hello"}`
	k := dedupeKey(line, "nonexistent")
	if k == "" {
		t.Error("expected non-empty fallback key")
	}
}
