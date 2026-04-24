package cmd

import (
	"encoding/json"
	"testing"
)

func TestFlattenMap_SimpleFields(t *testing.T) {
	input := map[string]interface{}{
		"level": "info",
		"msg":   "hello",
	}
	out := make(map[string]interface{})
	flattenMap(input, "", ".", out)
	if out["level"] != "info" {
		t.Errorf("expected level=info, got %v", out["level"])
	}
	if out["msg"] != "hello" {
		t.Errorf("expected msg=hello, got %v", out["msg"])
	}
}

func TestFlattenMap_NestedFields(t *testing.T) {
	input := map[string]interface{}{
		"http": map[string]interface{}{
			"method": "GET",
			"status": float64(200),
		},
	}
	out := make(map[string]interface{})
	flattenMap(input, "", ".", out)
	if out["http.method"] != "GET" {
		t.Errorf("expected http.method=GET, got %v", out["http.method"])
	}
	if out["http.status"] != float64(200) {
		t.Errorf("expected http.status=200, got %v", out["http.status"])
	}
}

func TestFlattenMap_DeeplyNested(t *testing.T) {
	input := map[string]interface{}{
		"a": map[string]interface{}{
			"b": map[string]interface{}{
				"c": "deep",
			},
		},
	}
	out := make(map[string]interface{})
	flattenMap(input, "", ".", out)
	if out["a.b.c"] != "deep" {
		t.Errorf("expected a.b.c=deep, got %v", out["a.b.c"])
	}
}

func TestFlattenMap_WithPrefix(t *testing.T) {
	input := map[string]interface{}{"key": "val"}
	out := make(map[string]interface{})
	flattenMap(input, "root", ".", out)
	if out["root.key"] != "val" {
		t.Errorf("expected root.key=val, got %v", out["root.key"])
	}
}

func TestFlattenMap_CustomSeparator(t *testing.T) {
	input := map[string]interface{}{
		"x": map[string]interface{}{"y": "z"},
	}
	out := make(map[string]interface{})
	flattenMap(input, "", "_", out)
	if out["x_y"] != "z" {
		t.Errorf("expected x_y=z, got %v", out["x_y"])
	}
}

func TestFlattenMap_RoundTrip(t *testing.T) {
	raw := `{"a":{"b":1,"c":"hello"},"d":true}`
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &obj); err != nil {
		t.Fatal(err)
	}
	out := make(map[string]interface{})
	flattenMap(obj, "", ".", out)
	if len(out) != 3 {
		t.Errorf("expected 3 flat keys, got %d: %v", len(out), out)
	}
}
