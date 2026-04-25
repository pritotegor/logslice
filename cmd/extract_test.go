package cmd

import (
	"testing"
)

func TestExtractFields_JSONMode_AllPresent(t *testing.T) {
	line := `{"level":"info","msg":"started","ts":"2024-01-01T00:00:00Z"}`
	out, err := extractFields(line, []string{"level", "msg"}, false, "\t")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out == "" {
		t.Fatal("expected non-empty output")
	}
	var obj map[string]interface{}
	if err := parseJSONLine(out, &obj); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if obj["level"] != "info" || obj["msg"] != "started" {
		t.Errorf("unexpected fields: %v", obj)
	}
	if _, ok := obj["ts"]; ok {
		t.Error("ts field should not be present")
	}
}

func TestExtractFields_JSONMode_MissingField(t *testing.T) {
	line := `{"level":"warn"}`
	out, err := extractFields(line, []string{"level", "missing"}, false, "\t")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var obj map[string]interface{}
	if err := parseJSONLine(out, &obj); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if _, ok := obj["missing"]; ok {
		t.Error("missing field should be absent from output")
	}
}

func TestExtractFields_RawMode(t *testing.T) {
	line := `{"level":"error","code":500}`
	out, err := extractFields(line, []string{"level", "code"}, true, "|")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "error|500" {
		t.Errorf("expected 'error|500', got %q", out)
	}
}

func TestExtractFields_RawMode_MissingField(t *testing.T) {
	line := `{"level":"debug"}`
	out, err := extractFields(line, []string{"level", "absent"}, true, ",")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "debug," {
		t.Errorf("expected 'debug,', got %q", out)
	}
}

func TestExtractFields_InvalidJSON(t *testing.T) {
	_, err := extractFields("not-json", []string{"level"}, false, "\t")
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func parseJSONLine(s string, v interface{}) error {
	import_json := func(data []byte, out interface{}) error {
		return jsonUnmarshal(data, out)
	}
	return import_json([]byte(s), v)
}

func jsonUnmarshal(data []byte, v interface{}) error {
	switch out := v.(type) {
	case *map[string]interface{}:
		var m map[string]interface{}
		if err := jsonDecode(data, &m); err != nil {
			return err
		}
		*out = m
		return nil
	}
	return jsonDecode(data, v)
}

func jsonDecode(data []byte, v interface{}) error {
	import (
		"encoding/json"
	)
	return json.Unmarshal(data, v)
}
