package cmd

import (
	"os"
	"strings"
	"testing"
)

func writeTempMergeFile(t *testing.T, lines []string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "merge-*.log")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	defer f.Close()
	for _, l := range lines {
		f.WriteString(l + "\n")
	}
	return f.Name()
}

func TestRunMerge_MissingFile(t *testing.T) {
	err := RunMerge([]string{"/nonexistent/a.log", "/nonexistent/b.log"}, &strings.Builder{})
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestRunMerge_StableSort(t *testing.T) {
	// Lines with identical timestamps should preserve file order (stable sort)
	ts := "2024-06-01T08:00:00Z"
	lines1 := []string{`{"time":"` + ts + `","msg":"from-file1"}`}
	lines2 := []string{`{"time":"` + ts + `","msg":"from-file2"}`}

	f1 := writeTempMergeFile(t, lines1)
	f2 := writeTempMergeFile(t, lines2)

	var sb strings.Builder
	if err := RunMerge([]string{f1, f2}, &sb); err != nil {
		t.Fatalf("RunMerge error: %v", err)
	}

	result := sb.String()
	idx1 := strings.Index(result, "from-file1")
	idx2 := strings.Index(result, "from-file2")
	if idx1 < 0 || idx2 < 0 {
		t.Fatalf("missing expected lines in output: %s", result)
	}
	if idx1 > idx2 {
		t.Errorf("expected file1 line before file2 line for equal timestamps")
	}
}

func TestRunMerge_MixedJSONAndRaw(t *testing.T) {
	lines1 := []string{`{"time":"2024-03-01T06:00:00Z","msg":"json-line"}`}
	lines2 := []string{"raw log line without json"}

	f1 := writeTempMergeFile(t, lines1)
	f2 := writeTempMergeFile(t, lines2)

	var sb strings.Builder
	if err := RunMerge([]string{f1, f2}, &sb); err != nil {
		t.Fatalf("RunMerge error: %v", err)
	}

	result := sb.String()
	if !strings.Contains(result, "json-line") || !strings.Contains(result, "raw log line") {
		t.Errorf("expected both lines in output, got: %s", result)
	}
}
