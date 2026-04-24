package cmd

import (
	"fmt"
	"strings"
	"testing"
)

func buildSummarizeLines(n int, withGroup bool) string {
	var sb strings.Builder
	envs := []string{"prod", "dev", "staging"}
	for i := 0; i < n; i++ {
		if withGroup {
			fmt.Fprintf(&sb, `{"env":%q,"latency":%d,"size":%d}`+"\n",
				envs[i%len(envs)], i%500, i%4096)
		} else {
			fmt.Fprintf(&sb, `{"latency":%d,"size":%d}`+"\n", i%500, i%4096)
		}
	}
	return sb.String()
}

func BenchmarkCollectSummary_NoGroup_1000(b *testing.B) {
	data := buildSummarizeLines(1000, false)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		collectSummary(strings.NewReader(data), []string{"latency", "size"}, "")
	}
}

func BenchmarkCollectSummary_NoGroup_10000(b *testing.B) {
	data := buildSummarizeLines(10000, false)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		collectSummary(strings.NewReader(data), []string{"latency", "size"}, "")
	}
}

func BenchmarkCollectSummary_GroupBy_1000(b *testing.B) {
	data := buildSummarizeLines(1000, true)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		collectSummary(strings.NewReader(data), []string{"latency", "size"}, "env")
	}
}

func BenchmarkFieldStats_Record(b *testing.B) {
	s := newFieldStats()
	for i := 0; i < b.N; i++ {
		s.record(float64(i))
	}
}
