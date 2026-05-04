package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func capturePivotOutput(r *strings.Reader, rowField, colField, valField string) (string, error) {
	var buf bytes.Buffer
	err := RunPivot(r, &buf, rowField, colField, valField)
	return buf.String(), err
}

func TestRunPivot_HeaderContainsColumns(t *testing.T) {
	input := strings.NewReader(
		`{"env":"prod","level":"error"}` + "\n" +
			`{"env":"prod","level":"info"}` + "\n" +
			`{"env":"dev","level":"warn"}` + "\n",
	)
	out, err := capturePivotOutput(input, "env", "level", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines (header + 2 rows), got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "env") {
		t.Errorf("header should start with row field, got: %s", lines[0])
	}
	if !strings.Contains(lines[0], "error") || !strings.Contains(lines[0], "info") {
		t.Errorf("header missing expected columns: %s", lines[0])
	}
}

func TestRunPivot_CountValues(t *testing.T) {
	input := strings.NewReader(
		`{"env":"prod","level":"error"}` + "\n" +
			`{"env":"prod","level":"error"}` + "\n" +
			`{"env":"prod","level":"info"}` + "\n",
	)
	out, err := capturePivotOutput(input, "env", "level", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "2") {
		t.Errorf("expected count of 2 in output, got: %s", out)
	}
}

func TestRunPivot_MissingColAndRowFlags(t *testing.T) {
	input := strings.NewReader(`{"env":"prod","level":"error"}` + "\n")
	var buf bytes.Buffer
	err := RunPivot(input, &buf, "", "level", "")
	if err == nil {
		t.Error("expected error when row field is empty")
	}
}

func TestRunPivot_EmptyInput(t *testing.T) {
	input := strings.NewReader("")
	out, err := capturePivotOutput(input, "env", "level", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	// Only header line with just the row field (no col values)
	if len(lines) != 1 {
		t.Errorf("expected 1 header line for empty input, got %d", len(lines))
	}
}
