package cmd

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func captureAnnotateOutput(t *testing.T, input string, fields []string, overwrite bool) string {
	t.Helper()
	r := strings.NewReader(input)
	var buf bytes.Buffer
	if err := RunAnnotate(r, &buf, fields, overwrite); err != nil {
		t.Fatalf("RunAnnotate error: %v", err)
	}
	return buf.String()
}

func TestRunAnnotate_InjectsNewField(t *testing.T) {
	input := `{"msg":"started"}` + "\n"
	out := captureAnnotateOutput(t, input, []string{"env=prod"}, false)
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(lines[0]), &obj); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if obj["env"] != "prod" {
		t.Errorf("expected env=prod, got %v", obj["env"])
	}
	if obj["msg"] != "started" {
		t.Errorf("expected msg=started to be preserved")
	}
}

func TestRunAnnotate_MultipleFields(t *testing.T) {
	input := `{"level":"info"}` + "\n"
	out := captureAnnotateOutput(t, input, []string{"env=prod", "region=eu-west-1"}, false)
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(strings.TrimSpace(out)), &obj); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if obj["env"] != "prod" || obj["region"] != "eu-west-1" {
		t.Errorf("missing injected fields: %v", obj)
	}
}

func TestRunAnnotate_SkipsEmptyLines(t *testing.T) {
	input := `{"msg":"a"}` + "\n\n" + `{"msg":"b"}` + "\n"
	out := captureAnnotateOutput(t, input, []string{"env=test"}, false)
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 lines, got %d: %v", len(lines), lines)
	}
}

func TestRunAnnotate_NonJSONPassthrough(t *testing.T) {
	input := "plain text line\n"
	out := captureAnnotateOutput(t, input, []string{"env=prod"}, false)
	if strings.TrimSpace(out) != "plain text line" {
		t.Errorf("expected passthrough of non-JSON, got: %q", out)
	}
}

func TestRunAnnotate_OverwriteExisting(t *testing.T) {
	input := `{"env":"staging"}` + "\n"
	out := captureAnnotateOutput(t, input, []string{"env=prod"}, true)
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(strings.TrimSpace(out)), &obj); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if obj["env"] != "prod" {
		t.Errorf("expected env to be overwritten to prod, got %v", obj["env"])
	}
}
