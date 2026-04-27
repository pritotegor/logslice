package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

var fieldmapCmd = &cobra.Command{
	Use:   "fieldmap [file]",
	Short: "Discover all field names and their value types across log lines",
	RunE:  RunFieldmap,
}

func init() {
	fieldmapCmd.Flags().IntP("limit", "n", 0, "Stop after scanning N lines (0 = all)")
	fieldmapCmd.Flags().BoolP("count", "c", false, "Show occurrence count per field")
	fieldmapCmd.Flags().StringP("prefix", "p", "", "Only show fields matching prefix")
	rootCmd.AddCommand(fieldmapCmd)
}

type fieldInfo struct {
	Types  map[string]int
	Count  int
}

func RunFieldmap(cmd *cobra.Command, args []string) error {
	limit, _ := cmd.Flags().GetInt("limit")
	showCount, _ := cmd.Flags().GetBool("count")
	prefix, _ := cmd.Flags().GetString("prefix")

	reader, err := openInput(args)
	if err != nil {
		return err
	}

	fields := make(map[string]*fieldInfo)
	scanner := bufio.NewScanner(reader)
	scanned := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var obj map[string]interface{}
		if err := json.Unmarshal([]byte(line), &obj); err != nil {
			continue
		}
		for k, v := range obj {
			if prefix != "" && !strings.HasPrefix(k, prefix) {
				continue
			}
			if _, ok := fields[k]; !ok {
				fields[k] = &fieldInfo{Types: make(map[string]int)}
			}
			fields[k].Count++
			fields[k].Types[jsonTypeName(v)]++
		}
		scanned++
		if limit > 0 && scanned >= limit {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	keys := make([]string, 0, len(fields))
	for k := range fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	w := io.Writer(os.Stdout)
	for _, k := range keys {
		info := fields[k]
		typeStrs := make([]string, 0, len(info.Types))
		for t, n := range info.Types {
			typeStrs = append(typeStrs, fmt.Sprintf("%s(%d)", t, n))
		}
		sort.Strings(typeStrs)
		if showCount {
			fmt.Fprintf(w, "%-30s [%s] occurrences=%d\n", k, strings.Join(typeStrs, ", "), info.Count)
		} else {
			fmt.Fprintf(w, "%-30s [%s]\n", k, strings.Join(typeStrs, ", "))
		}
	}
	return nil
}

func jsonTypeName(v interface{}) string {
	switch v.(type) {
	case string:
		return "string"
	case float64:
		return "number"
	case bool:
		return "bool"
	case nil:
		return "null"
	case map[string]interface{}:
		return "object"
	case []interface{}:
		return "array"
	default:
		return "unknown"
	}
}
