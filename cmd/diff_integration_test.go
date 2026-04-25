package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempDiffFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "log.jsonl")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestRunDiff_IdenticalFiles(t *testing.T) {
	content := "{\"level\":\"info\",\"msg\":\"started\"}\n{\"level\":\"warn\",\"msg\":\"retry\"}\n"
	p1 := writeTempDiffFile(t, content)
	p2 := writeTempDiffFile(t, content)

	var out strings.Builder
	err := diffReaders(&out, mustOpen(t, p1), mustOpen(t, p2), nil, false)
	if err != nil {
		t.Fatal(err)
	}
	if out.Len() != 0 {
		t.Errorf("expected no output for identical files, got: %s", out.String())
	}
}

func TestRunDiff_ChangedField(t *testing.T) {
	c1 := "{\"level\":\"info\",\"code\":200}\n"
	c2 := "{\"level\":\"error\",\"code\":500}\n"
	p1 := writeTempDiffFile(t, c1)
	p2 := writeTempDiffFile(t, c2)

	var out strings.Builder
	err := diffReaders(&out, mustOpen(t, p1), mustOpen(t, p2), nil, false)
	if err != nil {
		t.Fatal(err)
	}
	result := out.String()
	if !strings.Contains(result, "level") || !strings.Contains(result, "code") {
		t.Errorf("expected diffs for level and code, got: %s", result)
	}
}

func TestRunDiff_FieldFilter(t *testing.T) {
	c1 := "{\"level\":\"info\",\"code\":200}\n"
	c2 := "{\"level\":\"error\",\"code\":200}\n"
	p1 := writeTempDiffFile(t, c1)
	p2 := writeTempDiffFile(t, c2)

	var out strings.Builder
	err := diffReaders(&out, mustOpen(t, p1), mustOpen(t, p2), []string{"code"}, false)
	if err != nil {
		t.Fatal(err)
	}
	if out.Len() != 0 {
		t.Errorf("expected no diffs when filtering on equal field, got: %s", out.String())
	}
}

func TestRunDiff_EmptyFiles(t *testing.T) {
	p1 := writeTempDiffFile(t, "")
	p2 := writeTempDiffFile(t, "")

	var out strings.Builder
	err := diffReaders(&out, mustOpen(t, p1), mustOpen(t, p2), nil, false)
	if err != nil {
		t.Fatal(err)
	}
	if out.Len() != 0 {
		t.Errorf("expected no output for empty files")
	}
}

func mustOpen(t *testing.T, path string) *os.File {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { f.Close() })
	return f
}
