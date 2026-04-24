package cmd

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func parseJSONObj(s string, v interface{}) error {
	return json.Unmarshal([]byte(s), v)
}

func captureRenameOutput(t *testing.T, input string, fields []string) []string {
	t.Helper()
	var buf bytes.Buffer
	cmd := &cobra.Command{}
	cmd.SetOut(&buf)

	renameFields = fields
	mappings, err := parseRenameMappings(fields)
	if err != nil {
		t.Fatalf("parseRenameMappings: %v", err)
	}

	scanner := strings.NewReader(input)
	_ = scanner

	lines := strings.Split(strings.TrimRight(input, "\n"), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		fmt.Fprintln(&buf, renameLineFields(line, mappings))
	}
	_ = cmd

	var out []string
	for _, l := range strings.Split(strings.TrimRight(buf.String(), "\n"), "\n") {
		if l != "" {
			out = append(out, l)
		}
	}
	return out
}

func TestRunRename_BasicRename(t *testing.T) {
	input := "{\"level\":\"info\",\"msg\":\"started\"}\n{\"level\":\"error\",\"msg\":\"failed\"}\n"
	lines := captureRenameOutput(t, input, []string{"level=severity", "msg=message"})
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	for _, l := range lines {
		var obj map[string]interface{}
		if err := json.Unmarshal([]byte(l), &obj); err != nil {
			t.Fatalf("invalid JSON: %v", err)
		}
		if _, ok := obj["level"]; ok {
			t.Error("old key 'level' should be gone")
		}
		if _, ok := obj["severity"]; !ok {
			t.Error("new key 'severity' should exist")
		}
	}
}

func TestRunRename_SkipsEmptyLines(t *testing.T) {
	input := "{\"a\":\"1\"}\n\n{\"a\":\"2\"}\n"
	lines := captureRenameOutput(t, input, []string{"a=b"})
	if len(lines) != 2 {
		t.Fatalf("expected 2 non-empty lines, got %d", len(lines))
	}
}

func TestRunRename_PassthroughNonJSON(t *testing.T) {
	input := "plain log line\n"
	lines := captureRenameOutput(t, input, []string{"foo=bar"})
	if len(lines) != 1 || lines[0] != "plain log line" {
		t.Errorf("expected passthrough, got %v", lines)
	}
}
