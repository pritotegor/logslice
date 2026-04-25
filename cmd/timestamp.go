package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var timestampCmd = &cobra.Command{
	Use:   "timestamp",
	Short: "Add or reformat a timestamp field in JSON log lines",
	RunE:  RunTimestamp,
}

func init() {
	timestampCmd.Flags().String("field", "timestamp", "Field name to read/write the timestamp")
	timestampCmd.Flags().String("from-format", "", "Input time format (alias or Go layout); auto-detected if empty")
	timestampCmd.Flags().String("to-format", "rfc3339", "Output time format (alias or Go layout)")
	timestampCmd.Flags().Bool("add", false, "Add current UTC timestamp to lines that lack the field")
	timestampCmd.Flags().String("file", "", "Input file (default: stdin)")
	rootCmd.AddCommand(timestampCmd)
}

func RunTimestamp(cmd *cobra.Command, args []string) error {
	field, _ := cmd.Flags().GetString("field")
	fromFmt, _ := cmd.Flags().GetString("from-format")
	toFmt, _ := cmd.Flags().GetString("to-format")
	add, _ := cmd.Flags().GetBool("add")
	filePath, _ := cmd.Flags().GetString("file")

	outLayout := ResolveTimeFormat(toFmt)

	var inLayout string
	if fromFmt != "" {
		inLayout = ResolveTimeFormat(fromFmt)
	}

	reader, err := openInput(filePath)
	if err != nil {
		return err
	}
	if closer, ok := reader.(io.Closer); ok {
		defer closer.Close()
	}

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		out, err := reformatTimestamp(line, field, inLayout, outLayout, add)
		if err != nil {
			fmt.Fprintln(os.Stderr, "warn:", err)
			fmt.Println(line)
			continue
		}
		fmt.Println(out)
	}
	return scanner.Err()
}

func reformatTimestamp(line, field, inLayout, outLayout string, add bool) (string, error) {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line, nil
	}

	val, exists := obj[field]
	if !exists {
		if !add {
			return line, nil
		}
		obj[field] = time.Now().UTC().Format(outLayout)
		b, err := json.Marshal(obj)
		if err != nil {
			return line, err
		}
		return string(b), nil
	}

	raw, ok := val.(string)
	if !ok {
		return line, nil
	}

	var t time.Time
	var err error
	if inLayout != "" {
		t, err = time.Parse(inLayout, raw)
	} else {
		t, err = parseTimestamp(raw)
	}
	if err != nil {
		return line, fmt.Errorf("could not parse timestamp %q: %w", raw, err)
	}

	obj[field] = t.UTC().Format(outLayout)
	b, err := json.Marshal(obj)
	if err != nil {
		return line, err
	}
	return string(b), nil
}
