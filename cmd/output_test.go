package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestWriteLine_Raw(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf, FormatRaw, nil)
	if err := w.WriteLine(`{"level":"info","msg":"hello"}`); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "hello") {
		t.Errorf("expected raw line in output, got: %s", buf.String())
	}
}

func TestWriteLine_JSON_AllFields(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf, FormatJSON, nil)
	w.WriteLine(`{"level":"error","msg":"boom"}`)
	out := buf.String()
	if !strings.Contains(out, "error") || !strings.Contains(out, "boom") {
		t.Errorf("unexpected JSON output: %s", out)
	}
}

func TestWriteLine_JSON_FilteredFields(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf, FormatJSON, []string{"msg"})
	w.WriteLine(`{"level":"warn","msg":"watch out"}`)
	out := buf.String()
	if strings.Contains(out, "warn") {
		t.Errorf("filtered field 'level' should not appear: %s", out)
	}
	if !strings.Contains(out, "watch out") {
		t.Errorf("expected 'msg' field in output: %s", out)
	}
}

func TestWriteLine_CSV(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf, FormatCSV, []string{"level", "msg"})
	w.WriteLine(`{"level":"debug","msg":"trace me"}`)
	out := strings.TrimSpace(buf.String())
	if out != "debug,trace me" {
		t.Errorf("expected CSV 'debug,trace me', got: %s", out)
	}
}

func TestWriteLine_JSON_InvalidInput(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf, FormatJSON, nil)
	w.WriteLine("not json at all")
	out := buf.String()
	if !strings.Contains(out, "message") {
		t.Errorf("expected fallback message field, got: %s", out)
	}
}
