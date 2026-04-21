package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunConvert_JSONToRaw(t *testing.T) {
	input := `{"level":"info","msg":"started"}` + "\n" +
		`{"level":"error","msg":"failed"}` + "\n"

	var out bytes.Buffer
	writer := NewWriter(&out, "raw", nil)
	scanner := strings.NewReader(input)
	_ = scanner

	lines := []string{
		`{"level":"info","msg":"started"}`,
		`{"level":"error","msg":"failed"}`,
	}
	for _, line := range lines {
		m, err := convertParseLine(line, "json")
		if err != nil {
			t.Fatalf("parse error: %v", err)
		}
		if err := writer.WriteLine(line, m); err != nil {
			t.Fatalf("write error: %v", err)
		}
	}

	result := out.String()
	if !strings.Contains(result, "started") {
		t.Errorf("expected 'started' in output, got: %q", result)
	}
	if !strings.Contains(result, "failed") {
		t.Errorf("expected 'failed' in output, got: %q", result)
	}
}

func TestRunConvert_JSONToCSV(t *testing.T) {
	var out bytes.Buffer
	fields := []string{"level", "msg"}
	writer := NewWriter(&out, "csv", fields)

	line := `{"level":"warn","msg":"disk full","host":"srv1"}`
	m, err := convertParseLine(line, "json")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if err := writer.WriteLine(line, m); err != nil {
		t.Fatalf("write error: %v", err)
	}

	result := out.String()
	if !strings.Contains(result, "warn") {
		t.Errorf("expected 'warn' in CSV output, got: %q", result)
	}
	if strings.Contains(result, "srv1") {
		t.Errorf("expected 'srv1' to be excluded from filtered CSV, got: %q", result)
	}
}

func TestRunConvert_SkipsEmptyLines(t *testing.T) {
	lines := []string{
		`{"msg":"ok"}`,
		"",
		"   ",
		`{"msg":"also ok"}`,
	}

	var out bytes.Buffer
	writer := NewWriter(&out, "raw", nil)
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		m, _ := convertParseLine(line, "json")
		_ = writer.WriteLine(line, m)
	}

	result := out.String()
	lineCount := strings.Count(strings.TrimSpace(result), "\n") + 1
	if lineCount != 2 {
		t.Errorf("expected 2 output lines, got %d: %q", lineCount, result)
	}
}
