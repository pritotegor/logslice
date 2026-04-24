package cmd

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func buildSortLines(n int, field string) string {
	var sb strings.Builder
	for i := n; i > 0; i-- {
		fmt.Fprintf(&sb, `{%q:"%05d","msg":"line %d"}`+"\n", field, i, i)
	}
	return sb.String()
}

func BenchmarkRunSort_ByField_1000(b *testing.B) {
	input := buildSortLines(1000, "level")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var out bytes.Buffer
		if err := RunSort(strings.NewReader(input), &out, "level", false); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRunSort_ByField_10000(b *testing.B) {
	input := buildSortLines(10000, "level")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var out bytes.Buffer
		if err := RunSort(strings.NewReader(input), &out, "level", false); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkExtractSortKey_WithField(b *testing.B) {
	line := `{"level":"error","ts":"2024-01-01T00:00:00Z","msg":"something happened"}`
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		extractSortKey(line, "level")
	}
}

func BenchmarkExtractSortKey_Timestamp(b *testing.B) {
	line := `{"ts":"2024-01-01T00:00:00Z","msg":"something happened"}`
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		extractSortKey(line, "")
	}
}
