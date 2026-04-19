package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

func runFilter(filePath, fromStr, toStr, fieldFilter string) error {
	var from, to time.Time
	var err error

	if fromStr != "" {
		if from, err = time.Parse(time.RFC3339, fromStr); err != nil {
			return fmt.Errorf("invalid --from time: %w", err)
		}
	}
	if toStr != "" {
		if to, err = time.Parse(time.RFC3339, toStr); err != nil {
			return fmt.Errorf("invalid --to time: %w", err)
		}
	}

	var filterKey, filterVal string
	if fieldFilter != "" {
		parts := strings.SplitN(fieldFilter, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("--field must be in key=value format")
		}
		filterKey, filterVal = parts[0], parts[1]
	}

	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("cannot open file: %w", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		var entry map[string]interface{}
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			continue
		}
		if !matchesTime(entry, from, to) {
			continue
		}
		if filterKey != "" && !matchesField(entry, filterKey, filterVal) {
			continue
		}
		fmt.Println(line)
	}
	return scanner.Err()
}

func matchesTime(entry map[string]interface{}, from, to time.Time) bool {
	for _, key := range []string{"time", "timestamp", "ts", "@timestamp"} {
		val, ok := entry[key]
		if !ok {
			continue
		}
		ts, err := time.Parse(time.RFC3339, fmt.Sprintf("%v", val))
		if err != nil {
			continue
		}
		if !from.IsZero() && ts.Before(from) {
			return false
		}
		if !to.IsZero() && ts.After(to) {
			return false
		}
		return true
	}
	return true
}

func matchesField(entry map[string]interface{}, key, val string) bool {
	v, ok := entry[key]
	if !ok {
		return false
	}
	return fmt.Sprintf("%v", v) == val
}
