package generator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/format"
	"sort"
	"strings"
)

func GoStructFromJSON(input, rootName string) (string, error) {
	dec := json.NewDecoder(strings.NewReader(input))
	dec.UseNumber()
	var value any
	if err := dec.Decode(&value); err != nil {
		return "", err
	}
	return GoStructFromValue(value, rootName)
}

func GoStructFromValue(value any, rootName string) (string, error) {
	var buf bytes.Buffer
	seen := map[string]bool{}
	if err := writeGoStruct(&buf, exportName(rootName), value, seen); err != nil {
		return "", err
	}
	src, err := format.Source(buf.Bytes())
	if err != nil {
		return buf.String(), nil
	}
	return string(src), nil
}

func writeGoStruct(buf *bytes.Buffer, name string, value any, seen map[string]bool) error {
	obj, ok := value.(map[string]any)
	if !ok {
		return fmt.Errorf("root value must be an object")
	}
	if seen[name] {
		return nil
	}
	seen[name] = true

	keys := make([]string, 0, len(obj))
	for key := range obj {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	fmt.Fprintf(buf, "type %s struct {\n", name)
	nested := make([]struct {
		name  string
		value any
	}, 0)
	for _, key := range keys {
		fieldName := exportName(key)
		fieldType := goType(fieldName, obj[key])
		fmt.Fprintf(buf, "\t%s %s `json:\"%s\"`\n", fieldName, fieldType, key)
		if m, ok := obj[key].(map[string]any); ok {
			nested = append(nested, struct {
				name  string
				value any
			}{fieldName, m})
		}
	}
	buf.WriteString("}\n\n")
	for _, item := range nested {
		if err := writeGoStruct(buf, item.name, item.value, seen); err != nil {
			return err
		}
	}
	return nil
}

func goType(name string, value any) string {
	switch v := value.(type) {
	case bool:
		return "bool"
	case json.Number:
		if strings.Contains(v.String(), ".") {
			return "float64"
		}
		return "int64"
	case float64:
		return "float64"
	case int64, int:
		return "int64"
	case string:
		return "string"
	case []any:
		if len(v) == 0 {
			return "[]any"
		}
		return "[]" + goType(name, v[0])
	case map[string]any:
		return exportName(name)
	default:
		return "any"
	}
}
