package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func captureMaskOutput(t *testing.T, input string, fields []string, pattern, replacement string) []string {
	t.Helper()

	prev := maskFields
	prevPat := maskPattern
	prevRep := maskReplacement
	defer func() {
		maskFields = prev
		maskPattern = prevPat
		maskReplacement = prevRep
	}()

	maskFields = fields
	maskPattern = pattern
	maskReplacement = replacement

	reader := strings.NewReader(input)
	var re *regexp.Regexp
	var err error
	if pattern != "" {
		re, err = regexp.Compile(pattern)
		if err != nil {
			t.Fatalf("invalid pattern: %v", err)
		}
	}

	var buf bytes.Buffer
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		buf.WriteString(maskLine(line, fields, re, replacement) + "\n")
	}

	var lines []string
	for _, l := range strings.Split(strings.TrimRight(buf.String(), "\n"), "\n") {
		if l != "" {
			lines = append(lines, l)
		}
	}
	return lines
}

func TestRunMask_FieldRedaction(t *testing.T) {
	input := `{"user":"alice","password":"hunter2","level":"info"}
{"user":"bob","password":"letmein","level":"warn"}
`
	lines := captureMaskOutput(t, input, []string{"password"}, "", "***")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	for _, l := range lines {
		var obj map[string]interface{}
		if err := json.Unmarshal([]byte(l), &obj); err != nil {
			t.Fatalf("invalid JSON: %v", err)
		}
		if obj["password"] != "***" {
			t.Errorf("expected password masked, got %v", obj["password"])
		}
	}
}

func TestRunMask_PatternRedaction(t *testing.T) {
	input := `{"msg":"ssn is 123-45-6789","level":"info"}
{"msg":"no ssn here","level":"debug"}
`
	lines := captureMaskOutput(t, input, nil, `\d{3}-\d{2}-\d{4}`, "[SSN]")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(lines[0]), &obj); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if obj["msg"] != "ssn is [SSN]" {
		t.Errorf("expected SSN masked, got %v", obj["msg"])
	}
}

func TestRunMask_SkipsEmptyLines(t *testing.T) {
	input := `{"user":"alice","token":"abc"}

{"user":"bob","token":"xyz"}
`
	lines := captureMaskOutput(t, input, []string{"token"}, "", "***")
	if len(lines) != 2 {
		t.Errorf("expected 2 non-empty lines, got %d", len(lines))
	}
}
