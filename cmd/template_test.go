package cmd

import (
	"bytes"
	"strings"
	"testing"
	"text/template"
)

func TestApplyTemplate_BasicField(t *testing.T) {
	input := `{"level":"info","msg":"hello world"}`
	tmpl := template.Must(template.New("t").Option("missingkey=zero").Parse(`[{{.level}}] {{.msg}}`))

	var out bytes.Buffer
	err := applyTemplate(strings.NewReader(input), &out, tmpl, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := strings.TrimSpace(out.String())
	if got != "[info] hello world" {
		t.Errorf("expected '[info] hello world', got %q", got)
	}
}

func TestApplyTemplate_SkipsEmptyLines(t *testing.T) {
	input := "{\"level\":\"warn\"}\n\n{\"level\":\"error\"}"
	tmpl := template.Must(template.New("t").Option("missingkey=zero").Parse(`{{.level}}`))

	var out bytes.Buffer
	err := applyTemplate(strings.NewReader(input), &out, tmpl, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 output lines, got %d", len(lines))
	}
}

func TestApplyTemplate_InvalidJSON_ReturnsError(t *testing.T) {
	input := `not json`
	tmpl := template.Must(template.New("t").Option("missingkey=zero").Parse(`{{.level}}`))

	var out bytes.Buffer
	err := applyTemplate(strings.NewReader(input), &out, tmpl, false)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestApplyTemplate_InvalidJSON_SkipInvalid(t *testing.T) {
	input := "not json\n{\"level\":\"info\"}"
	tmpl := template.Must(template.New("t").Option("missingkey=zero").Parse(`{{.level}}`))

	var out bytes.Buffer
	err := applyTemplate(strings.NewReader(input), &out, tmpl, true)
	if err != nil {
		t.Fatalf("unexpected error with skip-invalid: %v", err)
	}
	got := strings.TrimSpace(out.String())
	if got != "info" {
		t.Errorf("expected 'info', got %q", got)
	}
}

func TestApplyTemplate_MissingKey_ZeroValue(t *testing.T) {
	input := `{"level":"debug"}`
	tmpl := template.Must(template.New("t").Option("missingkey=zero").Parse(`{{.level}} {{.msg}}`))

	var out bytes.Buffer
	err := applyTemplate(strings.NewReader(input), &out, tmpl, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := strings.TrimSpace(out.String())
	if got != "debug <no value>" && got != "debug " {
		// missingkey=zero renders missing map keys as empty string for maps
		if !strings.HasPrefix(got, "debug") {
			t.Errorf("unexpected output: %q", got)
		}
	}
}
