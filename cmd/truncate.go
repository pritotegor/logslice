package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var truncateCmd = &cobra.Command{
	Use:   "truncate",
	Short: "Truncate string fields in JSON log lines to a maximum length",
	RunE:  RunTruncate,
}

func init() {
	truncateCmd.Flags().StringSliceP("fields", "f", nil, "Fields to truncate (truncates all string fields if omitted)")
	truncateCmd.Flags().IntP("max-len", "n", 100, "Maximum length of string field values")
	truncateCmd.Flags().StringP("suffix", "s", "...", "Suffix appended to truncated values")
	truncateCmd.Flags().StringP("input", "i", "", "Input file (default: stdin)")
	rootCmd.AddCommand(truncateCmd)
}

func RunTruncate(cmd *cobra.Command, args []string) error {
	fields, _ := cmd.Flags().GetStringSlice("fields")
	maxLen, _ := cmd.Flags().GetInt("max-len")
	suffix, _ := cmd.Flags().GetString("suffix")
	inputFile, _ := cmd.Flags().GetString("input")

	reader, err := openInput(inputFile)
	if err != nil {
		return err
	}
	if closer, ok := reader.(io.Closer); ok {
		defer closer.Close()
	}

	fieldSet := make(map[string]bool, len(fields))
	for _, f := range fields {
		fieldSet[f] = true
	}

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		out, err := truncateLine(line, fieldSet, maxLen, suffix)
		if err != nil {
			fmt.Fprintln(os.Stderr, "skipping invalid line:", err)
			continue
		}
		fmt.Println(out)
	}
	return scanner.Err()
}

func truncateLine(line string, fields map[string]bool, maxLen int, suffix string) (string, error) {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return "", err
	}
	for k, v := range obj {
		if len(fields) > 0 && !fields[k] {
			continue
		}
		if s, ok := v.(string); ok && len(s) > maxLen {
			obj[k] = s[:maxLen] + suffix
		}
	}
	b, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
