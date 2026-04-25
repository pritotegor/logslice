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

var (
	compactFields []string
	compactDropEmpty bool
)

func init() {
	compactCmd := &cobra.Command{
		Use:   "compact [file]",
		Short: "Remove or trim fields from JSON log lines",
		Long:  `Remove specified fields from JSON log lines. Optionally drop lines with empty or null values.`,
		RunE:  RunCompact,
	}
	compactCmd.Flags().StringSliceVarP(&compactFields, "fields", "f", nil, "Fields to remove (comma-separated)")
	compactCmd.Flags().BoolVar(&compactDropEmpty, "drop-empty", false, "Drop fields with empty string or null values")
	rootCmd.AddCommand(compactCmd)
}

func RunCompact(cmd *cobra.Command, args []string) error {
	var reader io.Reader
	var err error
	if len(args) > 0 {
		reader, err = openInput(args[0])
		if err != nil {
			return err
		}
	} else {
		reader = os.Stdin
	}

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		result, err := compactLine(line, compactFields, compactDropEmpty)
		if err != nil {
			fmt.Fprintln(os.Stderr, "warn: "+err.Error())
			fmt.Println(line)
			continue
		}
		fmt.Println(result)
	}
	return scanner.Err()
}

func compactLine(line string, fields []string, dropEmpty bool) (string, error) {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line, fmt.Errorf("skipping non-JSON line")
	}

	removeSet := make(map[string]struct{}, len(fields))
	for _, f := range fields {
		removeSet[f] = struct{}{}
	}

	for key := range obj {
		if _, ok := removeSet[key]; ok {
			delete(obj, key)
			continue
		}
		if dropEmpty {
			val := obj[key]
			if val == nil {
				delete(obj, key)
				continue
			}
			if s, ok := val.(string); ok && s == "" {
				delete(obj, key)
			}
		}
	}

	out, err := json.Marshal(obj)
	if err != nil {
		return line, err
	}
	return string(out), nil
}
