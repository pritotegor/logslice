package cmd

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert log files between formats (json, csv, raw)",
	RunE:  RunConvert,
}

func init() {
	convertCmd.Flags().StringP("input", "i", "", "Input file (default: stdin)")
	convertCmd.Flags().StringP("from", "f", "json", "Source format: json, csv, raw")
	convertCmd.Flags().StringP("to", "t", "raw", "Target format: json, csv, raw")
	convertCmd.Flags().StringSliceP("fields", "k", nil, "Fields to include in output (default: all)")
	rootCmd.AddCommand(convertCmd)
}

func RunConvert(cmd *cobra.Command, args []string) error {
	inputFile, _ := cmd.Flags().GetString("input")
	fromFmt, _ := cmd.Flags().GetString("from")
	toFmt, _ := cmd.Flags().GetString("to")
	fields, _ := cmd.Flags().GetStringSlice("fields")

	if err := ValidateFormat(fromFmt); err != nil {
		return fmt.Errorf("invalid source format: %w", err)
	}
	if err := ValidateFormat(toFmt); err != nil {
		return fmt.Errorf("invalid target format: %w", err)
	}

	reader, err := openInput(inputFile)
	if err != nil {
		return err
	}
	if closer, ok := reader.(*os.File); ok && closer != os.Stdin {
		defer closer.Close()
	}

	writer := NewWriter(os.Stdout, toFmt, fields)
	scanner := bufio.NewScanner(reader)

	var lineNum int
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		fields_map, err := convertParseLine(line, fromFmt)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warning: skipping line %d: %v\n", lineNum, err)
			continue
		}
		if err := writer.WriteLine(line, fields_map); err != nil {
			return fmt.Errorf("line %d: %w", lineNum, err)
		}
	}
	return scanner.Err()
}

func convertParseLine(line, format string) (map[string]interface{}, error) {
	switch format {
	case "json":
		var m map[string]interface{}
		if err := json.Unmarshal([]byte(line), &m); err != nil {
			return nil, fmt.Errorf("invalid JSON: %w", err)
		}
		return m, nil
	case "csv":
		return parseCSVLine(line)
	default:
		return map[string]interface{}{"message": line}, nil
	}
}

func parseCSVLine(line string) (map[string]interface{}, error) {
	r := csv.NewReader(strings.NewReader(line))
	records, err := r.Read()
	if err != nil {
		return nil, fmt.Errorf("invalid CSV: %w", err)
	}
	m := make(map[string]interface{}, len(records))
	for i, v := range records {
		m[fmt.Sprintf("field%d", i+1)] = v
	}
	return m, nil
}
