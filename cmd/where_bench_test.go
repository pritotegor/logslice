package cmd

import (
	"fmt"
	"strings"
	"testing"
)

func buildWhereLines(n int) string {
	var sb strings.Builder
	levels := []string{"info", "warn", "error", "debug"}
	for i := 0; i < n; i++ {
		lvl := levels[i%len(levels)]
		status := 200 + (i%4)*100
		fmt.Fprintf(&sb, `{"level":%q,"status":%d,"msg":"log line %d"}`+"\n", lvl, status, i)
	}
	return sb.String()
}

func BenchmarkWhereLineMatches_StringEq_1000(b *testing.B) {
	lines := strings.Split(strings.TrimSpace(buildWhereLines(1000)), "\n")
	preds := []wherePredicate{{field: "level", op: "=", operand: "error"}}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, l := range lines {
			whereLineMatches(l, preds)
		}
	}
}

func BenchmarkWhereLineMatches_NumericGte_1000(b *testing.B) {
	lines := strings.Split(strings.TrimSpace(buildWhereLines(1000)), "\n")
	preds := []wherePredicate{{field: "status", op: ">=", operand: "400"}}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, l := range lines {
			whereLineMatches(l, preds)
		}
	}
}

func BenchmarkWhereLineMatches_MultiPredicate_1000(b *testing.B) {
	lines := strings.Split(strings.TrimSpace(buildWhereLines(1000)), "\n")
	preds := []wherePredicate{
		{field: "level", op: "=", operand: "error"},
		{field: "status", op: ">=", operand: "500"},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, l := range lines {
			whereLineMatches(l, preds)
		}
	}
}

func BenchmarkParseWhereExpressions(b *testing.B) {
	exprs := []string{"level=error", "status>=400", "msg contains timeout"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = parseWhereExpressions(exprs)
	}
}

func BenchmarkRunWhereOnReader_10000(b *testing.B) {
	input := buildWhereLines(10000)
	exprs := []string{"level=error", "status>=400"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := strings.NewReader(input)
		_ = runWhereOnReader(r, exprs, &strings.Builder{})
	}
}
