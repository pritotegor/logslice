package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

// FieldStats holds aggregated statistics for a single field.
type FieldStats struct {
	Count  int
	Values map[string]int
}

// LogStats holds statistics collected from log lines.
type LogStats struct {
	TotalLines  int
	ParsedLines int
	FieldCounts map[string]*FieldStats
}

// NewLogStats creates an initialized LogStats.
func NewLogStats() *LogStats {
	return &LogStats{
		FieldCounts: make(map[string]*FieldStats),
	}
}

// Record updates stats with the parsed fields from a log line.
func (s *LogStats) Record(fields map[string]interface{}) {
	s.TotalLines++
	if fields == nil {
		return
	}
	s.ParsedLines++
	for k, v := range fields {
		if _, ok := s.FieldCounts[k]; !ok {
			s.FieldCounts[k] = &FieldStats{Values: make(map[string]int)}
		}
		s.FieldCounts[k].Count++
		val := fmt.Sprintf("%v", v)
		s.FieldCounts[k].Values[val]++
	}
}

// Print writes a human-readable summary to stdout.
func (s *LogStats) Print() {
	fmt.Printf("Total lines : %d\n", s.TotalLines)
	fmt.Printf("Parsed (JSON): %d\n", s.ParsedLines)
	fmt.Printf("Fields found : %d\n", len(s.FieldCounts))

	keys := make([]string, 0, len(s.FieldCounts))
	for k := range s.FieldCounts {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fs := s.FieldCounts[k]
		fmt.Printf("  %-20s occurrences=%-6d unique_values=%d\n",
			k, fs.Count, len(fs.Values))
	}
}

var statsCmd = &cobra.Command{
	Use:   "stats [file]",
	Short: "Print field statistics from a structured log file",
	Args:  cobra.MaximumNArgs(1),
	RunE:  RunStats,
}

func init() {
	rootCmd.AddCommand(statsCmd)
	statsCmd.Flags().StringSlice("fields", nil, "Only report stats for these fields (comma-separated)")
}

func RunStats(cmd *cobra.Command, args []string) error {
	reader, err := openInput(args)
	if err != nil {
		return err
	}

	filterFields, _ := cmd.Flags().GetStringSlice("fields")
	wantField := map[string]bool{}
	for _, f := range filterFields {
		wantField[strings.TrimSpace(f)] = true
	}

	stats := NewLogStats()

	for reader.Scan() {
		line := reader.Text()
		parsed := ParseLine(line)
		if len(wantField) > 0 && parsed != nil {
			filtered := map[string]interface{}{}
			for k, v := range parsed {
				if wantField[k] {
					filtered[k] = v
				}
			}
			stats.Record(filtered)
		} else {
			stats.Record(parsed)
		}
	}

	stats.Print()
	return reader.Err()
}
