package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"sort"
	"strconv"
	"strings"
)

type fieldStats struct {
	Count int
	Sum   float64
	Min   float64
	Max   float64
}

func newFieldStats() *fieldStats {
	return &fieldStats{Min: math.MaxFloat64, Max: -math.MaxFloat64}
}

func (s *fieldStats) record(v float64) {
	s.Count++
	s.Sum += v
	if v < s.Min {
		s.Min = v
	}
	if v > s.Max {
		s.Max = v
	}
}

func (s *fieldStats) avg() float64 {
	if s.Count == 0 {
		return 0
	}
	return s.Sum / float64(s.Count)
}

// groupKey -> fieldName -> stats
type summaryResult map[string]map[string]*fieldStats

func collectSummary(r io.Reader, fields []string, groupBy string) summaryResult {
	result := make(summaryResult)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var obj map[string]interface{}
		if err := json.Unmarshal([]byte(line), &obj); err != nil {
			continue
		}
		groupKey := "_all"
		if groupBy != "" {
			if v, ok := obj[groupBy]; ok {
				groupKey = fmt.Sprintf("%v", v)
			} else {
				groupKey = "(missing)"
			}
		}
		if _, ok := result[groupKey]; !ok {
			result[groupKey] = make(map[string]*fieldStats)
		}
		for _, f := range fields {
			v, ok := obj[f]
			if !ok {
				continue
			}
			var fval float64
			switch tv := v.(type) {
			case float64:
				fval = tv
			case string:
				parsed, err := strconv.ParseFloat(tv, 64)
				if err != nil {
					continue
				}
				fval = parsed
			default:
				continue
			}
			if _, ok := result[groupKey][f]; !ok {
				result[groupKey][f] = newFieldStats()
			}
			result[groupKey][f].record(fval)
		}
	}
	return result
}

func printSummary(w io.Writer, results summaryResult, fields []string, groupBy string) {
	groups := make([]string, 0, len(results))
	for g := range results {
		groups = append(groups, g)
	}
	sort.Strings(groups)

	for _, g := range groups {
		if groupBy != "" {
			fmt.Fprintf(w, "[%s=%s]\n", groupBy, g)
		}
		for _, f := range fields {
			s, ok := results[g][f]
			if !ok {
				fmt.Fprintf(w, "  %-20s no data\n", f+":")
				continue
			}
			fmt.Fprintf(w, "  %-20s count=%-6d sum=%-12.4f min=%-12.4f max=%-12.4f avg=%.4f\n",
				f+":", s.Count, s.Sum, s.Min, s.Max, s.avg())
		}
	}
}
