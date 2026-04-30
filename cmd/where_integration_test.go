package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func captureWhereOutput(input string, exprs []string) (string, error) {
	reader := strings.NewReader(input)
	var buf bytes.Buffer
	err := runWhereOnReader(reader, exprs, &buf)
	return buf.String(), err
}

func TestRunWhere_FilterByStringField(t *testing.T) {
	input := strings.Join([]string{
		`{"level":"error","msg":"disk full"}`,
		`{"level":"info","msg":"started"}`,
		`{"level":"error","msg":"oom"}`,
	}, "\n") + "\n"

	out, err := captureWhereOutput(input, []string{"level=error"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 lines, got %d", len(lines))
	}
	for _, l := range lines {
		if !strings.Contains(l, `"error"`) {
			t.Errorf("unexpected line: %s", l)
		}
	}
}

func TestRunWhere_FilterByNumericField(t *testing.T) {
	input := strings.Join([]string{
		`{"status":200,"path":"/ok"}`,
		`{"status":404,"path":"/missing"}`,
		`{"status":500,"path":"/error"}`,
	}, "\n") + "\n"

	out, err := captureWhereOutput(input, []string{"status>=400"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 lines, got %d", len(lines))
	}
}

func TestRunWhere_MultipleExpressionsAND(t *testing.T) {
	input := strings.Join([]string{
		`{"level":"error","status":500}`,
		`{"level":"error","status":200}`,
		`{"level":"info","status":500}`,
	}, "\n") + "\n"

	out, err := captureWhereOutput(input, []string{"level=error", "status>=500"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 1 {
		t.Errorf("expected 1 line, got %d: %v", len(lines), lines)
	}
}

func TestRunWhere_SkipsEmptyLines(t *testing.T) {
	input := `{"level":"error"}` + "\n\n" + `{"level":"error"}` + "\n"
	out, err := captureWhereOutput(input, []string{"level=error"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 lines, got %d", len(lines))
	}
}

func TestRunWhere_NonJSONLinesExcluded(t *testing.T) {
	input := strings.Join([]string{
		`plain text line`,
		`{"level":"error"}`,
	}, "\n") + "\n"

	out, err := captureWhereOutput(input, []string{"level=error"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 1 {
		t.Errorf("expected 1 line, got %d", len(lines))
	}
}

func TestRunWhere_InvalidExpression(t *testing.T) {
	input := `{"level":"error"}` + "\n"
	_, err := captureWhereOutput(input, []string{"badexpr"})
	if err == nil {
		t.Error("expected error for invalid expression")
	}
}
