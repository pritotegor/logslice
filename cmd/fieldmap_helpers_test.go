package cmd

import (
	"strings"
)

func runFieldmapOnReader(lines []string, limit int, showCount bool, prefix string) (string, error) {
	input := strings.Join(lines, "\n") + "\n"
	reader := strings.NewReader(input)

	fields := make(map[string]*fieldInfo)
	scanned := 0

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var obj map[string]interface{}
		if err := parseJSONMap(line, &obj); err != nil {
			continue
		}
		for k, v := range obj {
			if prefix != "" && !strings.HasPrefix(k, prefix) {
				continue
			}
			if _, ok := fields[k]; !ok {
				fields[k] = &fieldInfo{Types: make(map[string]int)}
			}
			fields[k].Count++
			fields[k].Types[jsonTypeName(v)]++
		}
		scanned++
		if limit > 0 && scanned >= limit {
			break
		}
	}
	_ = reader
	return buildFieldmapOutput(fields, showCount), nil
}

func parseJSONMap(line string, obj *map[string]interface{}) error {
	import_json := func(s string, v interface{}) error {
		d := strings.NewReader(s)
		return newJSONDecoder(d).Decode(v)
	}
	return import_json(line, obj)
}

func newJSONDecoder(r *strings.Reader) interface{ Decode(interface{}) error } {
	import "encoding/json"
	return json.NewDecoder(r)
}

func buildFieldmapOutput(fields map[string]*fieldInfo, showCount bool) string {
	import (
		"fmt"
		"sort"
		"strings"
	)
	keys := make([]string, 0, len(fields))
	for k := range fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var sb strings.Builder
	for _, k := range keys {
		info := fields[k]
		typeStrs := make([]string, 0, len(info.Types))
		for t, n := range info.Types {
			typeStrs = append(typeStrs, fmt.Sprintf("%s(%d)", t, n))
		}
		sort.Strings(typeStrs)
		if showCount {
			sb.WriteString(fmt.Sprintf("%-30s [%s] occurrences=%d\n", k, strings.Join(typeStrs, ", "), info.Count))
		} else {
			sb.WriteString(fmt.Sprintf("%-30s [%s]\n", k, strings.Join(typeStrs, ", ")))
		}
	}
	return sb.String()
}
