package cmd

import (
	"testing"
)

func TestParseWhereExpr_Equals(t *testing.T) {
	p, err := parseWhereExpr("level=error")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.field != "level" || p.op != "=" || p.operand != "error" {
		t.Errorf("unexpected predicate: %+v", p)
	}
}

func TestParseWhereExpr_GreaterThanOrEqual(t *testing.T) {
	p, err := parseWhereExpr("status>=400")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.field != "status" || p.op != ">=" || p.operand != "400" {
		t.Errorf("unexpected predicate: %+v", p)
	}
}

func TestParseWhereExpr_Contains(t *testing.T) {
	p, err := parseWhereExpr("msg contains timeout")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.field != "msg" || p.op != "contains" || p.operand != "timeout" {
		t.Errorf("unexpected predicate: %+v", p)
	}
}

func TestParseWhereExpr_Invalid(t *testing.T) {
	_, err := parseWhereExpr("nodivider")
	if err == nil {
		t.Error("expected error for expression without operator")
	}
}

func TestWhereLineMatches_AllPredicatesMatch(t *testing.T) {
	line := `{"level":"error","status":500,"msg":"timeout occurred"}`
	preds := []wherePredicate{
		{field: "level", op: "=", operand: "error"},
		{field: "status", op: ">=", operand: "400"},
		{field: "msg", op: "contains", operand: "timeout"},
	}
	if !whereLineMatches(line, preds) {
		t.Error("expected line to match all predicates")
	}
}

func TestWhereLineMatches_OneFails(t *testing.T) {
	line := `{"level":"info","status":200}`
	preds := []wherePredicate{
		{field: "level", op: "=", operand: "error"},
	}
	if whereLineMatches(line, preds) {
		t.Error("expected line not to match")
	}
}

func TestWhereLineMatches_NonJSON(t *testing.T) {
	line := "plain text log line"
	preds := []wherePredicate{{field: "level", op: "=", operand: "error"}}
	if whereLineMatches(line, preds) {
		t.Error("non-JSON should not match")
	}
}

func TestWhereLineMatches_MissingField(t *testing.T) {
	line := `{"msg":"hello"}`
	preds := []wherePredicate{{field: "status", op: ">", operand: "0"}}
	if whereLineMatches(line, preds) {
		t.Error("missing field should not match")
	}
}

func TestWhereLineMatches_NotEquals(t *testing.T) {
	line := `{"level":"warn"}`
	preds := []wherePredicate{{field: "level", op: "!=", operand: "error"}}
	if !whereLineMatches(line, preds) {
		t.Error("expected != to match")
	}
}

func TestWhereLineMatches_StartsWith(t *testing.T) {
	line := `{"path":"/api/v2/users"}`
	preds := []wherePredicate{{field: "path", op: "startswith", operand: "/api"}}
	if !whereLineMatches(line, preds) {
		t.Error("expected startswith to match")
	}
}

func TestWhereLineMatches_EndsWith(t *testing.T) {
	line := `{"file":"app.log"}`
	preds := []wherePredicate{{field: "file", op: "endswith", operand: ".log"}}
	if !whereLineMatches(line, preds) {
		t.Error("expected endswith to match")
	}
}

func TestCompareNumeric_LessThan(t *testing.T) {
	if !compareNumeric("3", "<", "5") {
		t.Error("3 < 5 should be true")
	}
	if compareNumeric("5", "<", "3") {
		t.Error("5 < 3 should be false")
	}
}

func TestCompareNumeric_NonNumeric(t *testing.T) {
	if compareNumeric("abc", ">", "1") {
		t.Error("non-numeric comparison should return false")
	}
}
