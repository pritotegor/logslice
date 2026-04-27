package cmd

import (
	"testing"
)

func TestParseAnnotatePairs_Valid(t *testing.T) {
	pairs, err := parseAnnotatePairs([]string{"env=prod", "region=us-east-1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pairs["env"] != "prod" {
		t.Errorf("expected env=prod, got %q", pairs["env"])
	}
	if pairs["region"] != "us-east-1" {
		t.Errorf("expected region=us-east-1, got %q", pairs["region"])
	}
}

func TestParseAnnotatePairs_Invalid(t *testing.T) {
	_, err := parseAnnotatePairs([]string{"noequals"})
	if err == nil {
		t.Fatal("expected error for missing '='")
	}
}

func TestParseAnnotatePairs_EmptyKey(t *testing.T) {
	_, err := parseAnnotatePairs([]string{"=value"})
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestAnnotateLine_InjectsField(t *testing.T) {
	line := `{"msg":"hello"}`
	pairs := map[string]string{"env": "prod"}
	out, err := annotateLine(line, pairs, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsField(out, "env", "prod") {
		t.Errorf("expected env=prod in output: %s", out)
	}
}

func TestAnnotateLine_NoOverwrite(t *testing.T) {
	line := `{"env":"staging"}`
	pairs := map[string]string{"env": "prod"}
	out, err := annotateLine(line, pairs, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsField(out, "env", "staging") {
		t.Errorf("expected env to remain staging: %s", out)
	}
}

func TestAnnotateLine_WithOverwrite(t *testing.T) {
	line := `{"env":"staging"}`
	pairs := map[string]string{"env": "prod"}
	out, err := annotateLine(line, pairs, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsField(out, "env", "prod") {
		t.Errorf("expected env to be overwritten to prod: %s", out)
	}
}

func TestAnnotateLine_NonJSON_ReturnsError(t *testing.T) {
	_, err := annotateLine("not json", map[string]string{"k": "v"}, false)
	if err == nil {
		t.Fatal("expected error for non-JSON input")
	}
}

// containsField is a small helper to check key/value presence in a JSON string.
func containsField(jsonStr, key, val string) bool {
	expected := `"` + key + `":"` + val + `"`
	return len(jsonStr) > 0 && (len(expected) == 0 || contains(jsonStr, expected))
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsStr(s, sub))
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
