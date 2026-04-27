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

var annotateFields []string
var annotateValue string
var annotateOverwrite bool

func init() {
	annotateCmd := &cobra.Command{
		Use:   "annotate",
		Short: "Add or set fields on each JSON log line",
		Long:  `Annotate each JSON log line by injecting one or more key=value fields. Non-JSON lines are passed through unchanged.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			input, err := openInput(args)
			if err != nil {
				return err
			}
			defer input.Close()
			return RunAnnotate(input, os.Stdout, annotateFields, annotateOverwrite)
		},
	}

	annotateCmd.Flags().StringArrayVarP(&annotateFields, "field", "f", nil, "Field to inject in key=value format (repeatable)")
	annotateCmd.Flags().BoolVar(&annotateOverwrite, "overwrite", false, "Overwrite existing fields")
	_ = annotateCmd.MarkFlagRequired("field")

	rootCmd.AddCommand(annotateCmd)
}

func RunAnnotate(r io.Reader, w io.Writer, fields []string, overwrite bool) error {
	pairs, err := parseAnnotatePairs(fields)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		annotated, err := annotateLine(line, pairs, overwrite)
		if err != nil {
			fmt.Fprintln(w, line)
			continue
		}
		fmt.Fprintln(w, annotated)
	}
	return scanner.Err()
}

// parseAnnotatePairs parses a slice of "key=value" strings into a map.
// Returns an error if any entry does not contain a '=' or has an empty key.
func parseAnnotatePairs(fields []string) (map[string]string, error) {
	pairs := make(map[string]string, len(fields))
	for _, f := range fields {
		idx := strings.IndexByte(f, '=')
		if idx < 1 {
			return nil, fmt.Errorf("invalid field annotation %q: expected key=value", f)
		}
		key := f[:idx]
		if strings.TrimSpace(key) == "" {
			return nil, fmt.Errorf("invalid field annotation %q: key must not be blank", f)
		}
		pairs[key] = f[idx+1:]
	}
	return pairs, nil
}

func annotateLine(line string, pairs map[string]string, overwrite bool) (string, error) {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return "", err
	}
	for k, v := range pairs {
		if _, exists := obj[k]; exists && !overwrite {
			continue
		}
		obj[k] = v
	}
	out, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(out), nil
}
