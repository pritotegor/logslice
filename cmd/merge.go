package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/spf13/cobra"
)

var mergeCmd = &cobra.Command{
	Use:   "merge [files...]",
	Short: "Merge multiple log files sorted by timestamp",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunMerge(args, os.Stdout)
	},
}

func init() {
	rootCmd.AddCommand(mergeCmd)
}

type mergeEntry struct {
	line string
	key  string
}

// RunMerge reads multiple log files and merges them in timestamp order.
func RunMerge(paths []string, out io.Writer) error {
	var entries []mergeEntry

	for _, path := range paths {
		f, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("open %s: %w", path, err)
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			if line == "" {
				continue
			}
			entries = append(entries, mergeEntry{
				line: line,
				key:  extractMergeKey(line),
			})
		}
		if err := scanner.Err(); err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}
	}

	sort.SliceStable(entries, func(i, j int) bool {
		return entries[i].key < entries[j].key
	})

	w := bufio.NewWriter(out)
	for _, e := range entries {
		fmt.Fprintln(w, e.line)
	}
	return w.Flush()
}

// extractMergeKey returns a sortable key for a log line.
// For JSON lines, it uses the parsed timestamp; otherwise falls back to the raw line.
func extractMergeKey(line string) string {
	parsed := ParseLine(line)
	if parsed.Timestamp != nil {
		return parsed.Timestamp.UTC().Format("2006-01-02T15:04:05.999999999Z")
	}
	return line
}
