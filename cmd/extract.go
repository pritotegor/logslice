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

var extractCmd = &cobra.Command{
	Use:   "extract",
	Short: "Extract one or more fields from JSON log lines",
	Long:  `Reads JSON log lines and outputs only the specified fields, optionally as raw values or JSON objects.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunExtract(cmd, os.Stdin)
	},
}

func init() {
	extractCmd.Flags().StringSliceP("fields", "f", nil, "Comma-separated list of fields to extract (required)")
	extractCmd.Flags().BoolP("raw", "r", false, "Output raw field values instead of JSON objects")
	extractCmd.Flags().StringP("separator", "s", "\t", "Separator between fields in raw mode")
	_ = extractCmd.MarkFlagRequired("fields")
	rootCmd.AddCommand(extractCmd)
}

func RunExtract(cmd *cobra.Command, r io.Reader) error {
	fields, _ := cmd.Flags().GetStringSlice("fields")
	rawMode, _ := cmd.Flags().GetBool("raw")
	separator, _ := cmd.Flags().GetString("separator")

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		out, err := extractFields(line, fields, rawMode, separator)
		if err != nil {
			continue
		}
		fmt.Println(out)
	}
	return scanner.Err()
}

func extractFields(line string, fields []string, rawMode bool, separator string) (string, error) {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return "", err
	}

	if rawMode {
		parts := make([]string, 0, len(fields))
		for _, f := range fields {
			if val, ok := obj[f]; ok {
				parts = append(parts, fmt.Sprintf("%v", val))
			} else {
				parts = append(parts, "")
			}
		}
		return strings.Join(parts, separator), nil
	}

	result := make(map[string]interface{}, len(fields))
	for _, f := range fields {
		if val, ok := obj[f]; ok {
			result[f] = val
		}
	}
	b, err := json.Marshal(result)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
