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

var flattenCmd = &cobra.Command{
	Use:   "flatten [file]",
	Short: "Flatten nested JSON log fields into dot-notation keys",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunFlatten(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(flattenCmd)
	flattenCmd.Flags().StringP("prefix", "p", "", "Optional key prefix to apply to all output fields")
	flattenCmd.Flags().StringP("separator", "s", ".", "Separator to use between nested key segments")
}

func RunFlatten(cmd *cobra.Command, args []string) error {
	sep, _ := cmd.Flags().GetString("separator")
	prefix, _ := cmd.Flags().GetString("prefix")
	if sep == "" {
		sep = "."
	}

	var r io.Reader
	if len(args) == 1 {
		f, err := os.Open(args[0])
		if err != nil {
			return fmt.Errorf("open file: %w", err)
		}
		defer f.Close()
		r = f
	} else {
		r = os.Stdin
	}

	scanner := bufio.NewScanner(r)
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var obj map[string]interface{}
		if err := json.Unmarshal([]byte(line), &obj); err != nil {
			fmt.Fprintln(w, line)
			continue
		}
		flat := make(map[string]interface{})
		flattenMap(obj, prefix, sep, flat)
		out, err := json.Marshal(flat)
		if err != nil {
			fmt.Fprintln(w, line)
			continue
		}
		fmt.Fprintln(w, string(out))
	}
	return scanner.Err()
}

func flattenMap(obj map[string]interface{}, prefix, sep string, out map[string]interface{}) {
	for k, v := range obj {
		key := k
		if prefix != "" {
			key = prefix + sep + k
		}
		switch child := v.(type) {
		case map[string]interface{}:
			flattenMap(child, key, sep, out)
		default:
			out[key] = v
		}
	}
}
