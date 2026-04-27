package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var splitField string
var splitOutDir string
var splitPrefix string

func init() {
	splitCmd := &cobra.Command{
		Use:   "split",
		Short: "Split log lines into separate files by field value",
		Long:  `Reads log lines and writes each line to a separate file based on the value of a JSON field.`,
		RunE:  RunSplit,
	}
	splitCmd.Flags().StringVarP(&splitField, "field", "f", "", "JSON field to split on (required)")
	splitCmd.Flags().StringVarP(&splitOutDir, "out", "o", ".", "Output directory for split files")
	splitCmd.Flags().StringVar(&splitPrefix, "prefix", "split_", "Filename prefix for output files")
	_ = splitCmd.MarkFlagRequired("field")
	rootCmd.AddCommand(splitCmd)
}

func RunSplit(cmd *cobra.Command, args []string) error {
	reader, err := openInput(args)
	if err != nil {
		return err
	}
	if rc, ok := reader.(io.ReadCloser); ok {
		defer rc.Close()
	}

	if err := os.MkdirAll(splitOutDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	writers := map[string]*os.File{}
	defer func() {
		for _, f := range writers {
			f.Close()
		}
	}()

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		key := splitKeyValue(line, splitField)
		if key == "" {
			key = "_unknown"
		}
		f, ok := writers[key]
		if !ok {
			filename := filepath.Join(splitOutDir, splitPrefix+sanitizeFilename(key)+".log")
			f, err = os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			if err != nil {
				return fmt.Errorf("failed to open output file %s: %w", filename, err)
			}
			writers[key] = f
		}
		fmt.Fprintln(f, line)
	}
	return scanner.Err()
}

func splitKeyValue(line, field string) string {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return ""
	}
	v, ok := obj[field]
	if !ok {
		return ""
	}
	return fmt.Sprintf("%v", v)
}

func sanitizeFilename(s string) string {
	replacer := strings.NewReplacer(
		"/", "_",
		"\\", "_",
		":", "_",
		"*", "_",
		"?", "_",
		"\"", "_",
		"<", "_",
		">", "_",
		"|", "_",
		" ", "_",
	)
	return replacer.Replace(s)
}
