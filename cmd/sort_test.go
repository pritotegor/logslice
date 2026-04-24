package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestExtractSortKey_NoField_UsesTimestamp(t *testing.T) {
	line := `{"ts":"2024-01-02T10:00:00Z","msg":"hello"}`
	key := extractSortKey(line, "")
	if key != "2024-01-02T10:00:00Z" {
		t.Errorf("expected timestamp key, got %q", key)
	}
}

func TestExtractSortKey_WithField(t *testing.T) {
	line := `{"level":"error","msg":"oops"}`
	key := extractSortKey(line, "level")
	if key != "error" {
		t.Errorf("expected 'error', got %q", key)
	}
}

func TestExtractSortKey_MissingField_ReturnsEmpty(t *testing.T) {
	line := `{"msg":"hello"}`
	key := extractSortKey(line, "level")
	if key != "" {
		t.Errorf("expected empty string, got %q", key)
	}
}

func TestExtractSortKey_NonJSON_ReturnsLine(t *testing.T) {
	line := "plain text log line"
	key := extractSortKey(line, "")
	if key != line {
		t.Errorf("expected line as key, got %q", key)
	}
}

func TestRunSort_ByField_Ascending(t *testing.T) {
	input := strings.Join([]string{
		`{"level":"warn","msg":"b"}`,
		`{"level":"error","msg":"a"}`,
		`{"level":"info","msg":"c"}`,
	}, "\n") + "\n"

	var out bytes.Buffer
	if err := RunSort(strings.NewReader(input), &out, "level", false); err != nil {
		t.Fatalf("RunSort error: %v", err)
	}

	lines := strings.Split(strings.TrimRight(out.String(), "\n"), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.Contains(lines[0], "error") {
		t.Errorf("expected first line to have level=error, got %q", lines[0])
	}
	if !strings.Contains(lines[2], "warn") {
		t.Errorf("expected last line to have level=warn, got %q", lines[2])
	}
}

func TestRunSort_ByField_Descending(t *testing.T) {
	input := strings.Join([]string{
		`{"level":"info","msg":"c"}`,
		`{"level":"error","msg":"a"}`,
		`{"level":"warn","msg":"b"}`,
	}, "\n") + "\n"

	var out bytes.Buffer
	if err := RunSort(strings.NewReader(input), &out, "level", true); err != nil {
		t.Fatalf("RunSort error: %v", err)
	}

	lines := strings.Split(strings.TrimRight(out.String(), "\n"), "\n")
	if !strings.Contains(lines[0], "warn") {
		t.Errorf("expected first line to have level=warn, got %q", lines[0])
	}
}

func TestRunSort_SkipsEmptyLines(t *testing.T) {
	input := `{"level":"info"}` + "\n\n" + `{"level":"error"}` + "\n"
	var out bytes.Buffer
	if err := RunSort(strings.NewReader(input), &out, "level", false); err != nil {
		t.Fatalf("RunSort error: %v", err)
	}
	lines := strings.Split(strings.TrimRight(out.String(), "\n"), "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 lines, got %d", len(lines))
	}
}
