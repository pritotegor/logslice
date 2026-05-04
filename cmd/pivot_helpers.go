package cmd

import (
	"bufio"
	"encoding/json"
	"io"
	"strings"
)

// pivotCell holds the aggregated numeric value for a (row, col) pair.
type pivotCell = float64

// collectPivot reads lines from r and builds the pivot table.
// If valField is empty, each occurrence increments the cell by 1 (count).
// If valField is set, the numeric value of that field is summed.
func collectPivot(
	r io.Reader,
	rowField, colField, valField string,
) (table map[string]map[string]pivotCell, rowOrder []string, colSet map[string]struct{}, err error) {
	table = make(map[string]map[string]pivotCell)
	colSet = make(map[string]struct{})
	seen := make(map[string]bool)

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var fields map[string]interface{}
		if jsonErr := json.Unmarshal([]byte(line), &fields); jsonErr != nil {
			continue
		}

		rowVal, ok := stringField(fields, rowField)
		if !ok {
			continue
		}
		colVal, ok := stringField(fields, colField)
		if !ok {
			continue
		}

		var increment float64 = 1
		if valField != "" {
			increment = numericField(fields, valField)
		}

		if !seen[rowVal] {
			seen[rowVal] = true
			rowOrder = append(rowOrder, rowVal)
		}
		colSet[colVal] = struct{}{}

		if table[rowVal] == nil {
			table[rowVal] = make(map[string]pivotCell)
		}
		table[rowVal][colVal] += increment
	}
	err = scanner.Err()
	return
}

// stringField extracts a field value as a string from a decoded JSON map.
func stringField(fields map[string]interface{}, key string) (string, bool) {
	v, ok := fields[key]
	if !ok {
		return "", false
	}
	switch t := v.(type) {
	case string:
		return t, true
	default:
		return "", false
	}
}

// numericField extracts a numeric field value; returns 0 if missing or non-numeric.
func numericField(fields map[string]interface{}, key string) float64 {
	v, ok := fields[key]
	if !ok {
		return 0
	}
	switch t := v.(type) {
	case float64:
		return t
	case int:
		return float64(t)
	}
	return 0
}
