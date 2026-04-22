package cmd

import (
	"regexp"
	"testing"
)

func TestGrepMatch_FullLine_Matches(t *testing.T) {
	re := regexp.MustCompile(`error`)
	line := `{"level":"error","msg":"something failed"}`
	if !grepMatch(re, line, "") {
		t.Error("expected match on full line")
	}
}

func TestGrepMatch_FullLine_NoMatch(t *testing.T) {
	re := regexp.MustCompile(`panic`)
	line := `{"level":"info","msg":"all good"}`
	if grepMatch(re, line, "") {
		t.Error("expected no match")
	}
}

func TestGrepMatch_Field_Matches(t *testing.T) {
	re := regexp.MustCompile(`failed`)
	line := `{"level":"error","msg":"something failed"}`
	if !grepMatch(re, line, "msg") {
		t.Error("expected match on field 'msg'")
	}
}

func TestGrepMatch_Field_NoMatch(t *testing.T) {
	re := regexp.MustCompile(`failed`)
	line := `{"level":"error","msg":"all clear"}`
	if grepMatch(re, line, "msg") {
		t.Error("expected no match on field 'msg'")
	}
}

func TestGrepMatch_Field_Missing(t *testing.T) {
	re := regexp.MustCompile(`anything`)
	line := `{"level":"info"}`
	if grepMatch(re, line, "msg") {
		t.Error("expected no match when field is absent")
	}
}

func TestGrepMatch_Field_NonJSON(t *testing.T) {
	re := regexp.MustCompile(`hello`)
	line := `plain text hello world`
	if grepMatch(re, line, "msg") {
		t.Error("expected no match for non-JSON line with field filter")
	}
}

func TestGrepMatch_CaseInsensitive(t *testing.T) {
	re := regexp.MustCompile(`(?i)ERROR`)
	line := `{"level":"error","msg":"oops"}`
	if !grepMatch(re, line, "") {
		t.Error("expected case-insensitive match")
	}
}

func TestGrepMatch_NumericField(t *testing.T) {
	re := regexp.MustCompile(`^42`)
	line := `{"code":42,"msg":"ok"}`
	if !grepMatch(re, line, "code") {
		t.Error("expected match on numeric field value")
	}
}
