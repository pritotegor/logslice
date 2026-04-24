package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/spf13/cobra"
)

var countCmd = &cobra.Command{
	Use:   "count",
	Short: "Count log lines, optionally grouped by a field value",
	RunE:  RunCount,
}

func init() {
	countCmd.Flags().StringP("field", "f", "", "Group counts by this JSON field")
	countCmd.Flags().StringP("input", "i", "", "Input file (default: stdin)")
	countCmd.Flags().BoolP("sort", "s", false, "Sort output by count descending")
	rootCmd.AddCommand(countCmd)
}

func RunCount(cmd *cobra.Command, args []string) error {
	field, _ := cmd.Flags().GetString("field")
	inputFile, _ := cmd.Flags().GetString("input")
	sortByCount, _ := cmd.Flags().GetBool("sort")

	reader, err := openInput(inputFile)
	if err != nil {
		return fmt.Errorf("opening input: %w", err)
	}
	if closer, ok := reader.(io.Closer); ok && inputFile != "" {
		defer closer.Close()
	}

	counts, total, err := countLines(reader, field)
	if err != nil {
		return err
	}

	if field == "" {
		fmt.Fprintf(os.Stdout, "%d\n", total)
		return nil
	}

	type kv struct {
		Key   string
		Count int
	}
	pairs := make([]kv, 0, len(counts))
	for k, v := range counts {
		pairs = append(pairs, kv{k, v})
	}
	if sortByCount {
		sort.Slice(pairs, func(i, j int) bool {
			return pairs[i].Count > pairs[j].Count
		})
	} else {
		sort.Slice(pairs, func(i, j int) bool {
			return pairs[i].Key < pairs[j].Key
		})
	}
	for _, p := range pairs {
		fmt.Fprintf(os.Stdout, "%s\t%d\n", p.Key, p.Count)
	}
	return nil
}

func countLines(r io.Reader, field string) (map[string]int, int, error) {
	counts := make(map[string]int)
	total := 0
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		total++
		if field == "" {
			continue
		}
		var fields map[string]interface{}
		if err := json.Unmarshal([]byte(line), &fields); err != nil {
			counts["<unparsed>"]++
			continue
		}
		val, ok := fields[field]
		if !ok {
			counts["<missing>"]++
			continue
		}
		counts[fmt.Sprintf("%v", val)]++
	}
	return counts, total, scanner.Err()
}
