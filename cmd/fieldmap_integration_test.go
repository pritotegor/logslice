package cmd

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func captureFieldmapOutput(lines []string, limit int, showCount bool, prefix string) (string, error) {
	input := strings.Join(lines, "\n") + "\n"
	reader := strings.NewReader(input)

	fields := make(map[string]*fieldInfo)
	scanned := 0
	scanner := newLineScanner(reader)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var obj map[string]interface{}
		if err := json.Unmarshal([]byte(line), &obj); err != nil {
			continue
		}
		for k, v := range obj {
			if prefix != "" && !strings.HasPrefix(k, prefix) {
				continue
			}
			if _, ok := fields[k]; !ok {
				fields[k] = &fieldInfo{Types: make(map[string]int)}
			}
			fields[k].Count++
			fields[k].Types[jsonTypeName(v)]++
		}
		scanned++
		if limit > 0 && scanned >= limit {
			break
		}
	}

	var buf bytes.Buffer
	writeFieldmapOutput(&buf, fields, showCount)
	return buf.String(), scanner.Err()
}

func TestRunFieldmap_BasicFields(t *testing.T) {
	lines := []string{
		`{"level":"info","msg":"started","pid":42}`,
		`{"level":"warn","msg":"slow","pid":99}`,
	}
	out, err := captureFieldmapOutput(lines, 0, false, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, field := range []string{"level", "msg", "pid"} {
		if !strings.Contains(out, field) {
			t.Errorf("expected field %q in output, got:\n%s", field, out)
		}
	}
}

func TestRunFieldmap_ShowCount(t *testing.T) {
	lines := []string{
		`{"level":"info"}`,
		`{"level":"warn"}`,
		`{"level":"error"}`,
	}
	out, err := captureFieldmapOutput(lines, 0, true, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "occurrences=3") {
		t.Errorf("expected occurrences=3 in output, got:\n%s", out)
	}
}

func TestRunFieldmap_PrefixFilter(t *testing.T) {
	lines := []string{
		`{"http_status":200,"http_method":"GET","level":"info"}`,
	}
	out, err := captureFieldmapOutput(lines, 0, false, "http_")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, "level") {
		t.Errorf("expected level to be filtered out, got:\n%s", out)
	}
	if !strings.Contains(out, "http_status") {
		t.Errorf("expected http_status in output, got:\n%s", out)
	}
}

func TestRunFieldmap_SkipsNonJSON(t *testing.T) {
	lines := []string{
		"not json at all",
		`{"level":"info"}`,
	}
	out, err := captureFieldmapOutput(lines, 0, false, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "level") {
		t.Errorf("expected level field in output, got:\n%s", out)
	}
}

func TestRunFieldmap_LimitLines(t *testing.T) {
	lines := []string{
		`{"a":1}`,
		`{"b":2}`,
		`{"c":3}`,
	}
	out, err := captureFieldmapOutput(lines, 1, false, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, "b") || strings.Contains(out, "c") {
		t.Errorf("expected only field 'a' due to limit=1, got:\n%s", out)
	}
	if !strings.Contains(out, "a") {
		t.Errorf("expected field 'a' in output, got:\n%s", out)
	}
}
