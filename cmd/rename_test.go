package cmd

import (
	"testing"
)

func TestParseRenameMappings_Valid(t *testing.T) {
	mappings, err := parseRenameMappings([]string{"old=new", "foo=bar"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mappings["old"] != "new" {
		t.Errorf("expected old->new, got %q", mappings["old"])
	}
	if mappings["foo"] != "bar" {
		t.Errorf("expected foo->bar, got %q", mappings["foo"])
	}
}

func TestParseRenameMappings_Invalid(t *testing.T) {
	cases := []string{"noequals", "=nokey", "noval="}
	for _, c := range cases {
		_, err := parseRenameMappings([]string{c})
		if err == nil {
			t.Errorf("expected error for mapping %q", c)
		}
	}
}

func TestRenameLineFields_ValidJSON(t *testing.T) {
	line := `{"level":"info","msg":"hello"}`
	mappings := map[string]string{"level": "severity", "msg": "message"}
	result := renameLineFields(line, mappings)

	var obj map[string]interface{}
	if err := parseJSONObj(result, &obj); err != nil {
		t.Fatalf("result is not valid JSON: %v", err)
	}
	if _, ok := obj["level"]; ok {
		t.Error("old key 'level' should not exist")
	}
	if obj["severity"] != "info" {
		t.Errorf("expected severity=info, got %v", obj["severity"])
	}
	if obj["message"] != "hello" {
		t.Errorf("expected message=hello, got %v", obj["message"])
	}
}

func TestRenameLineFields_NonJSON(t *testing.T) {
	line := "plain text log line"
	mappings := map[string]string{"foo": "bar"}
	result := renameLineFields(line, mappings)
	if result != line {
		t.Errorf("expected passthrough, got %q", result)
	}
}

func TestRenameLineFields_MissingKey(t *testing.T) {
	line := `{"level":"warn"}`
	mappings := map[string]string{"nonexistent": "other"}
	result := renameLineFields(line, mappings)

	var obj map[string]interface{}
	if err := parseJSONObj(result, &obj); err != nil {
		t.Fatalf("result is not valid JSON: %v", err)
	}
	if obj["level"] != "warn" {
		t.Errorf("expected level=warn to be preserved, got %v", obj["level"])
	}
}

func TestRenameLineFields_NoMappings(t *testing.T) {
	line := `{"level":"debug","count":3}`
	mappings := map[string]string{}
	result := renameLineFields(line, mappings)
	var obj map[string]interface{}
	if err := parseJSONObj(result, &obj); err != nil {
		t.Fatalf("result is not valid JSON: %v", err)
	}
	if obj["level"] != "debug" {
		t.Errorf("expected level=debug, got %v", obj["level"])
	}
}
