package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var whereCmd = &cobra.Command{
	Use:   "where",
	Short: "Filter log lines by numeric or string field comparisons",
	Long: `Filter JSON log lines using field comparisons.
Supports operators: =, !=, <, <=, >, >=, contains, startswith, endswith`,
	RunE: RunWhere,
}

var whereExpressions []string

func init() {
	rootCmd.AddCommand(whereCmd)
	whereCmd.Flags().StringArrayVarP(&whereExpressions, "expr", "e", nil, "Filter expression (e.g. 'status>=400' or 'level=error')")
	_ = whereCmd.MarkFlagRequired("expr")
}

func RunWhere(cmd *cobra.Command, args []string) error {
	reader, err := openInput(args)
	if err != nil {
		return err
	}

	preds, err := parseWhereExpressions(whereExpressions)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		if whereLineMatches(line, preds) {
			fmt.Fprintln(os.Stdout, line)
		}
	}
	return scanner.Err()
}

func whereLineMatches(line string, preds []wherePredicate) bool {
	var fields map[string]interface{}
	if err := json.Unmarshal([]byte(line), &fields); err != nil {
		return false
	}
	for _, p := range preds {
		if !p.evaluate(fields) {
			return false
		}
	}
	return true
}

func runWhereOnReader(r io.Reader, exprs []string, w io.Writer) error {
	preds, err := parseWhereExpressions(exprs)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		if whereLineMatches(line, preds) {
			fmt.Fprintln(w, line)
		}
	}
	return scanner.Err()
}

type wherePredicate struct {
	field    string
	op       string
	operand  string
}

func (p wherePredicate) evaluate(fields map[string]interface{}) bool {
	val, ok := fields[p.field]
	if !ok {
		return false
	}
	strVal := fmt.Sprintf("%v", val)
	switch p.op {
	case "=":
		return strVal == p.operand
	case "!=":
		return strVal != p.operand
	case "contains":
		return strings.Contains(strVal, p.operand)
	case "startswith":
		return strings.HasPrefix(strVal, p.operand)
	case "endswith":
		return strings.HasSuffix(strVal, p.operand)
	case "<", "<=", ">", ">=":
		return compareNumeric(strVal, p.op, p.operand)
	}
	return false
}

func compareNumeric(left, op, right string) bool {
	l, errL := strconv.ParseFloat(left, 64)
	r, errR := strconv.ParseFloat(right, 64)
	if errL != nil || errR != nil {
		return false
	}
	switch op {
	case "<":
		return l < r
	case "<=":
		return l <= r
	case ">":
		return l > r
	case ">=":
		return l >= r
	}
	return false
}

var whereOps = []string{"<=", ">=", "!=", "<", ">", "=", "contains", "startswith", "endswith"}

func parseWhereExpressions(exprs []string) ([]wherePredicate, error) {
	var preds []wherePredicate
	for _, expr := range exprs {
		p, err := parseWhereExpr(expr)
		if err != nil {
			return nil, err
		}
		preds = append(preds, p)
	}
	return preds, nil
}

func parseWhereExpr(expr string) (wherePredicate, error) {
	for _, op := range whereOps {
		if idx := strings.Index(expr, op); idx > 0 {
			field := strings.TrimSpace(expr[:idx])
			operand := strings.TrimSpace(expr[idx+len(op):])
			if field != "" {
				return wherePredicate{field: field, op: op, operand: operand}, nil
			}
		}
	}
	return wherePredicate{}, fmt.Errorf("invalid expression: %q", expr)
}
