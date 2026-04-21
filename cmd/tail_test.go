package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestReadLastN_FewerLinesThanN(t *testing.T) {
	input := "line1\nline2\nline3\n"
	lines, err := readLastN(strings.NewReader(input), 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 3 {
		t.Errorf("expected 3 lines, got %d", len(lines))
	}
}

func TestReadLastN_ExactN(t *testing.T) {
	input := "line1\nline2\nline3\n"
	lines, err := readLastN(strings.NewReader(input), 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 3 {
		t.Errorf("expected 3 lines, got %d", len(lines))
	}
}

func TestReadLastN_MoreLinesThanN(t *testing.T) {
	input := "a\nb\nc\nd\ne\n"
	lines, err := readLastN(strings.NewReader(input), 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 3 {
		t.Errorf("expected 3 lines, got %d", len(lines))
	}
	if lines[0] != "c" || lines[1] != "d" || lines[2] != "e" {
		t.Errorf("unexpected lines: %v", lines)
	}
}

func TestReadLastN_EmptyInput(t *testing.T) {
	lines, err := readLastN(strings.NewReader(""), 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 0 {
		t.Errorf("expected 0 lines, got %d", len(lines))
	}
}

func TestRunTail_BasicOutput(t *testing.T) {
	// Simulate RunTail by reading last 2 lines and writing to buffer.
	input := "first\nsecond\nthird\n"
	lines, err := readLastN(strings.NewReader(input), 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var buf bytes.Buffer
	for _, l := range lines {
		buf.WriteString(l + "\n")
	}

	got := buf.String()
	if !strings.Contains(got, "second") || !strings.Contains(got, "third") {
		t.Errorf("unexpected output: %q", got)
	}
	if strings.Contains(got, "first") {
		t.Errorf("should not contain first line: %q", got)
	}
}
