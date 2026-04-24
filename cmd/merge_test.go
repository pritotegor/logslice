package cmd

import (
	"strings"
	"testing"
)

func TestExtractMergeKey_NonJSON(t *testing.T) {
	line := "plain log line"
	key := extractMergeKey(line)
	if key != line {
		t.Errorf("expected raw line as key, got %q", key)
	}
}

func TestExtractMergeKey_ValidJSON(t *testing.T) {
	line := `{"time":"2024-01-02T10:00:00Z","msg":"hello"}`
	key := extractMergeKey(line)
	if !strings.HasPrefix(key, "2024-01-02") {
		t.Errorf("expected timestamp-based key, got %q", key)
	}
}

func TestExtractMergeKey_JSONNoTimestamp(t *testing.T) {
	line := `{"level":"info","msg":"no time"}`
	key := extractMergeKey(line)
	if key != line {
		t.Errorf("expected raw line fallback, got %q", key)
	}
}

func TestRunMerge_OrdersByTimestamp(t *testing.T) {
	lines1 := []string{
		`{"time":"2024-01-02T10:00:00Z","msg":"first"}`,
		`{"time":"2024-01-02T12:00:00Z","msg":"third"}`,
	}
	lines2 := []string{
		`{"time":"2024-01-02T11:00:00Z","msg":"second"}`,
		`{"time":"2024-01-02T13:00:00Z","msg":"fourth"}`,
	}

	f1 := writeTempMergeFile(t, lines1)
	f2 := writeTempMergeFile(t, lines2)

	var sb strings.Builder
	if err := RunMerge([]string{f1, f2}, &sb); err != nil {
		t.Fatalf("RunMerge error: %v", err)
	}

	result := strings.Split(strings.TrimSpace(sb.String()), "\n")
	expected := []string{"first", "second", "third", "fourth"}
	for i, exp := range expected {
		if !strings.Contains(result[i], exp) {
			t.Errorf("line %d: expected %q in %q", i, exp, result[i])
		}
	}
}

func TestRunMerge_SkipsEmptyLines(t *testing.T) {
	lines1 := []string{`{"time":"2024-01-01T09:00:00Z","msg":"a"}`, ""}
	lines2 := []string{"", `{"time":"2024-01-01T10:00:00Z","msg":"b"}`}

	f1 := writeTempMergeFile(t, lines1)
	f2 := writeTempMergeFile(t, lines2)

	var sb strings.Builder
	if err := RunMerge([]string{f1, f2}, &sb); err != nil {
		t.Fatalf("RunMerge error: %v", err)
	}

	count := strings.Count(strings.TrimSpace(sb.String()), "\n") + 1
	if count != 2 {
		t.Errorf("expected 2 output lines, got %d", count)
	}
}
