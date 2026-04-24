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

var sortCmd = &cobra.Command{
	Use:   "sort [file]",
	Short: "Sort log lines by a field or timestamp",
	Long:  `Read log lines and emit them sorted by a JSON field value (lexicographic) or by timestamp.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		field, _ := cmd.Flags().GetString("field")
		reverse, _ := cmd.Flags().GetBool("reverse")
		r, err := openInput(args)
		if err != nil {
			return err
		}
		if closer, ok := r.(io.Closer); ok {
			defer closer.Close()
		}
		return RunSort(r, os.Stdout, field, reverse)
	},
}

func init() {
	sortCmd.Flags().StringP("field", "f", "", "JSON field to sort by (default: timestamp)")
	sortCmd.Flags().BoolP("reverse", "r", false, "Sort in descending order")
	rootCmd.AddCommand(sortCmd)
}

type sortEntry struct {
	raw string
	key string
}

func RunSort(r io.Reader, w io.Writer, field string, reverse bool) error {
	scanner := bufio.NewScanner(r)
	var entries []sortEntry

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		key := extractSortKey(line, field)
		entries = append(entries, sortEntry{raw: line, key: key})
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("reading input: %w", err)
	}

	sort.SliceStable(entries, func(i, j int) bool {
		if reverse {
			return entries[i].key > entries[j].key
		}
		return entries[i].key < entries[j].key
	})

	bw := bufio.NewWriter(w)
	for _, e := range entries {
		fmt.Fprintln(bw, e.raw)
	}
	return bw.Flush()
}

func extractSortKey(line, field string) string {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}
	if field == "" {
		// default to timestamp
		for _, k := range []string{"time", "timestamp", "ts", "@timestamp"} {
			if v, ok := obj[k]; ok {
				return fmt.Sprintf("%v", v)
			}
		}
		return line
	}
	if v, ok := obj[field]; ok {
		return fmt.Sprintf("%v", v)
	}
	return ""
}
