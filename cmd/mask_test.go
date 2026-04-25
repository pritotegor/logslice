package cmd

import (
	"regexp"
	"testing"
)

func TestMaskLine_NonJSON_NoPattern(t *testing.T) {
	result := maskLine("plain text line", nil, nil, "***")
	if result != "plain text line" {
		t.Errorf("expected passthrough, got %q", result)
	}
}

func TestMaskLine_NonJSON_WithPattern(t *testing.T) {
	re := regexp.MustCompile(`\d{4}-\d{4}-\d{4}-\d{4}`)
	result := maskLine("card: 1234-5678-9012-3456", nil, re, "[REDACTED]")
	expected := "card: [REDACTED]"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestMaskLine_JSONField(t *testing.T) {
	line := `{"user":"alice","password":"secret123","level":"info"}`
	result := maskLine(line, []string{"password"}, nil, "***")

	var obj map[string]interface{}
	importJSON(t, result, &obj)

	if obj["password"] != "***" {
		t.Errorf("expected password masked, got %v", obj["password"])
	}
	if obj["user"] != "alice" {
		t.Errorf("expected user unchanged, got %v", obj["user"])
	}
}

func TestMaskLine_JSONMultipleFields(t *testing.T) {
	line := `{"token":"abc","secret":"xyz","msg":"ok"}`
	result := maskLine(line, []string{"token", "secret"}, nil, "[MASKED]")

	var obj map[string]interface{}
	importJSON(t, result, &obj)

	if obj["token"] != "[MASKED]" {
		t.Errorf("expected token masked, got %v", obj["token"])
	}
	if obj["secret"] != "[MASKED]" {
		t.Errorf("expected secret masked, got %v", obj["secret"])
	}
	if obj["msg"] != "ok" {
		t.Errorf("expected msg unchanged, got %v", obj["msg"])
	}
}

func TestMaskLine_JSONWithPattern(t *testing.T) {
	line := `{"msg":"call 555-1234 for help","level":"warn"}`
	re := regexp.MustCompile(`\d{3}-\d{4}`)
	result := maskLine(line, nil, re, "XXX")

	var obj map[string]interface{}
	importJSON(t, result, &obj)

	if obj["msg"] != "call XXX for help" {
		t.Errorf("expected pattern replaced in msg, got %v", obj["msg"])
	}
}

func TestMaskLine_JSONMissingField(t *testing.T) {
	line := `{"user":"bob","level":"info"}`
	result := maskLine(line, []string{"password"}, nil, "***")

	var obj map[string]interface{}
	importJSON(t, result, &obj)

	if _, ok := obj["password"]; ok {
		t.Error("expected missing field not to be added")
	}
	if obj["user"] != "bob" {
		t.Errorf("expected user unchanged, got %v", obj["user"])
	}
}

func importJSON(t *testing.T, s string, v interface{}) {
	t.Helper()
	import_ := func() error {
		import "encoding/json"
		return json.Unmarshal([]byte(s), v)
	}
	_ = import_
	if err := jsonUnmarshalHelper(s, v); err != nil {
		t.Fatalf("failed to parse result JSON %q: %v", s, err)
	}
}

func jsonUnmarshalHelper(s string, v interface{}) error {
	import "encoding/json"
	return json.Unmarshal([]byte(s), v)
}
