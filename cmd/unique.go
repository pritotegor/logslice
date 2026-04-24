package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

var uniqueField string
var uniqueCount bool

func init() {
	uniqueCmd := &cobra.Command{
		Use:   "unique",
		Short: "Print unique values for a given JSON field across log lines",
		RunE:  RunUnique,
	}
	uniqueCmd.Flags().StringVarP(&uniqueField, "field", "f", "", "JSON field to extract unique values from (required)")
	uniqueCmd.Flags().BoolVarP(&uniqueCount, "count", "c", false, "Print occurrence count alongside each unique value")
	_ = uniqueCmd.MarkFlagRequired("field")
	rootCmd.AddCommand(uniqueCmd)
}

func RunUnique(cmd *cobra.Command, args []string) error {
	reader, err := openInput(args)
	if err != nil {
		return err
	}
	if rc, ok := reader.(io.Closer); ok {
		defer rc.Close()
	}

	counts, order, err := collectUniqueValues(reader, uniqueField)
	if err != nil {
		return err
	}

	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	for _, val := range order {
		if uniqueCount {
			fmt.Fprintf(w, "%s\t%d\n", val, counts[val])
		} else {
			fmt.Fprintln(w, val)
		}
	}
	return nil
}

func collectUniqueValues(r io.Reader, field string) (map[string]int, []string, error) {
	counts := make(map[string]int)
	var order []string

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		val := extractUniqueFieldValue(line, field)
		if val == "" {
			continue
		}
		if _, seen := counts[val]; !seen {
			order = append(order, val)
		}
		counts[val]++
	}
	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}
	return counts, order, nil
}

func extractUniqueFieldValue(line, field string) string {
	var fields map[string]interface{}
	if err := json.Unmarshal([]byte(line), &fields); err != nil {
		return ""
	}
	v, ok := fields[field]
	if !ok {
		return ""
	}
	switch val := v.(type) {
	case string:
		return val
	case float64:
		return fmt.Sprintf("%g", val)
	case bool:
		return fmt.Sprintf("%t", val)
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return ""
		}
		return string(b)
	}
}
