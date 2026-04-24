package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"
)

func captureFlattenOutput(t *testing.T, input string) []string {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w

	cmd := flattenCmd
	cmd.ResetFlags()
	cmd.Flags().StringP("prefix", "p", "", "")
	cmd.Flags().StringP("separator", "s", ".", "")

	origStdin := os.Stdin
	pr, pw, _ := os.Pipe()
	os.Stdin = pr
	_, _ = io.WriteString(pw, input)
	pw.Close()

	_ = RunFlatten(cmd, []string{})

	w.Close()
	os.Stdout = old
	os.Stdin = origStdin

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	r.Close()

	var lines []string
	scanner := bufio.NewScanner(&buf)
	for scanner.Scan() {
		l := strings.TrimSpace(scanner.Text())
		if l != "" {
			lines = append(lines, l)
		}
	}
	return lines
}

func TestRunFlatten_NestedJSON(t *testing.T) {
	input := `{"http":{"method":"POST","status":201},"level":"info"}` + "\n"
	lines := captureFlattenOutput(t, input)
	if len(lines) != 1 {
		t.Fatalf("expected 1 output line, got %d", len(lines))
	}
	var out map[string]interface{}
	if err := json.Unmarshal([]byte(lines[0]), &out); err != nil {
		t.Fatalf("output not valid JSON: %v", err)
	}
	if out["http.method"] != "POST" {
		t.Errorf("expected http.method=POST, got %v", out["http.method"])
	}
	if out["level"] != "info" {
		t.Errorf("expected level=info, got %v", out["level"])
	}
}

func TestRunFlatten_NonJSONPassthrough(t *testing.T) {
	input := "plain log line\n"
	lines := captureFlattenOutput(t, input)
	if len(lines) != 1 || lines[0] != "plain log line" {
		t.Errorf("expected passthrough of raw line, got %v", lines)
	}
}

func TestRunFlatten_SkipsEmptyLines(t *testing.T) {
	input := "\n\n{\"a\":{\"b\":\"c\"}}\n\n"
	lines := captureFlattenOutput(t, input)
	if len(lines) != 1 {
		t.Errorf("expected 1 non-empty line, got %d: %v", len(lines), lines)
	}
}
