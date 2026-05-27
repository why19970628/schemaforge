package parser

import (
	"fmt"
	"strconv"
	"strings"
)

func ParseYAML(input string) (any, error) {
	root := map[string]any{}
	lines := strings.Split(input, "\n")
	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "- ") {
			return nil, fmt.Errorf("top-level arrays are not supported by the built-in YAML parser yet")
		}
		key, value, ok := strings.Cut(line, ":")
		if !ok {
			return nil, fmt.Errorf("invalid YAML line %q", line)
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if key == "" {
			return nil, fmt.Errorf("empty YAML key")
		}
		root[key] = parseScalar(value)
	}
	if len(root) == 0 {
		return nil, fmt.Errorf("no YAML fields found")
	}
	return root, nil
}

func parseScalar(value string) any {
	value = strings.Trim(value, `"'`)
	if value == "" {
		return ""
	}
	if value == "true" {
		return true
	}
	if value == "false" {
		return false
	}
	if i, err := strconv.ParseInt(value, 10, 64); err == nil {
		return i
	}
	if f, err := strconv.ParseFloat(value, 64); err == nil {
		return f
	}
	return value
}
