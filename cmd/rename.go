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

var renameFields []string

func init() {
	renameCmd := &cobra.Command{
		Use:   "rename",
		Short: "Rename fields in JSON log lines",
		Long:  `Rename one or more fields in each JSON log line using old=new pairs.`,
		RunE:  RunRename,
	}

	renameCmd.Flags().StringArrayVarP(&renameFields, "field", "f", nil, "Field rename mapping as old=new (repeatable)")
	_ = renameCmd.MarkFlagRequired("field")
	rootCmd.AddCommand(renameCmd)
}

func RunRename(cmd *cobra.Command, args []string) error {
	mappings, err := parseRenameMappings(renameFields)
	if err != nil {
		return err
	}

	var reader io.Reader
	if len(args) > 0 {
		f, err := os.Open(args[0])
		if err != nil {
			return fmt.Errorf("open file: %w", err)
		}
		defer f.Close()
		reader = f
	} else {
		reader = os.Stdin
	}

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		fmt.Fprintln(cmd.OutOrStdout(), renameLineFields(line, mappings))
	}
	return scanner.Err()
}

func parseRenameMappings(pairs []string) (map[string]string, error) {
	mappings := make(map[string]string, len(pairs))
	for _, p := range pairs {
		parts := strings.SplitN(p, "=", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return nil, fmt.Errorf("invalid rename mapping %q: expected old=new", p)
		}
		mappings[parts[0]] = parts[1]
	}
	return mappings, nil
}

func renameLineFields(line string, mappings map[string]string) string {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}
	for oldKey, newKey := range mappings {
		if val, ok := obj[oldKey]; ok {
			delete(obj, oldKey)
			obj[newKey] = val
		}
	}
	out, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return string(out)
}
