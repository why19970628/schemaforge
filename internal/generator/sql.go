package generator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/format"
	"strings"

	"github.com/why19970628/schemaforge/internal/parser"
)

func EntSchemas(tables []parser.Table) (string, error) {
	var out strings.Builder
	for _, table := range tables {
		name := exportName(singularTableName(table.Name))
		fmt.Fprintf(&out, "package schema\n\n")
		fmt.Fprintf(&out, "import (\n\t\"entgo.io/ent\"\n\t\"entgo.io/ent/schema/field\"\n\t\"entgo.io/ent/schema/index\"\n)\n\n")
		fmt.Fprintf(&out, "type %s struct {\n\tent.Schema\n}\n\n", name)
		fmt.Fprintf(&out, "func (%s) Fields() []ent.Field {\n\treturn []ent.Field{\n", name)
		for _, col := range table.Columns {
			fmt.Fprintf(&out, "\t\t%s,\n", entField(col))
		}
		fmt.Fprintf(&out, "\t}\n}\n\n")
		fmt.Fprintf(&out, "func (%s) Indexes() []ent.Index {\n\treturn []ent.Index{\n", name)
		for _, idx := range table.Indexes {
			if len(idx.Columns) == 0 {
				continue
			}
			fmt.Fprintf(&out, "\t\tindex.Fields(%s)", quoteList(idx.Columns))
			if idx.Unique {
				out.WriteString(".Unique()")
			}
			out.WriteString(",\n")
		}
		fmt.Fprintf(&out, "\t}\n}\n")
	}
	return out.String(), nil
}

func GORMModels(tables []parser.Table) (string, error) {
	var buf bytes.Buffer
	buf.WriteString("package models\n\n")
	if gormNeedsTime(tables) {
		buf.WriteString("import \"time\"\n\n")
	}
	for _, table := range tables {
		name := exportName(singularTableName(table.Name))
		fmt.Fprintf(&buf, "type %s struct {\n", name)
		for _, col := range table.Columns {
			field := exportName(col.Name)
			tag := gormTag(col)
			fmt.Fprintf(&buf, "\t%s %s `gorm:\"%s\" json:\"%s\"`\n", field, sqlGoType(col.Type), tag, col.Name)
		}
		fmt.Fprintf(&buf, "}\n\n")
		fmt.Fprintf(&buf, "func (%s) TableName() string {\n\treturn %q\n}\n\n", name, table.Name)
	}
	src, err := format.Source(buf.Bytes())
	if err != nil {
		return buf.String(), nil
	}
	return string(src), nil
}

func gormNeedsTime(tables []parser.Table) bool {
	for _, table := range tables {
		for _, col := range table.Columns {
			if sqlGoType(col.Type) == "time.Time" {
				return true
			}
		}
	}
	return false
}

func ESMappings(tables []parser.Table) (string, error) {
	root := map[string]any{}
	for _, table := range tables {
		props := map[string]any{}
		for _, col := range table.Columns {
			props[col.Name] = map[string]any{"type": esType(col.Type)}
		}
		root[table.Name] = map[string]any{"mappings": map[string]any{"properties": props}}
	}
	data, err := json.MarshalIndent(root, "", "  ")
	return string(data), err
}

func MongoSchemas(tables []parser.Table) (string, error) {
	root := map[string]any{}
	for _, table := range tables {
		properties := map[string]any{}
		required := []string{}
		for _, col := range table.Columns {
			properties[col.Name] = map[string]any{
				"bsonType":    mongoType(col.Type),
				"description": strings.TrimSpace(col.Comment),
			}
			if !col.Nullable {
				required = append(required, col.Name)
			}
		}
		root[table.Name] = map[string]any{
			"validator": map[string]any{
				"$jsonSchema": map[string]any{
					"bsonType":   "object",
					"required":   required,
					"properties": properties,
				},
			},
		}
	}
	data, err := json.MarshalIndent(root, "", "  ")
	return string(data), err
}

func entField(col parser.Column) string {
	var b strings.Builder
	fmt.Fprintf(&b, "field.%s(%q)", entType(col.Type), col.Name)
	if col.Primary || col.AutoIncrement {
		b.WriteString(".Unique()")
	}
	if col.Nullable {
		b.WriteString(".Optional().Nillable()")
	}
	if col.Comment != "" {
		fmt.Fprintf(&b, ".Comment(%q)", col.Comment)
	}
	return b.String()
}

func entType(t string) string {
	t = strings.ToLower(t)
	switch {
	case strings.Contains(t, "bool"), strings.Contains(t, "tinyint(1)"):
		return "Bool"
	case strings.Contains(t, "bigint"):
		return "Int64"
	case strings.Contains(t, "int"):
		return "Int"
	case strings.Contains(t, "decimal"), strings.Contains(t, "double"), strings.Contains(t, "float"):
		return "Float"
	case strings.Contains(t, "time"), strings.Contains(t, "date"):
		return "Time"
	case strings.Contains(t, "json"):
		return "JSON"
	default:
		return "String"
	}
}

func sqlGoType(t string) string {
	t = strings.ToLower(t)
	switch {
	case strings.Contains(t, "bool"), strings.Contains(t, "tinyint(1)"):
		return "bool"
	case strings.Contains(t, "bigint"):
		return "int64"
	case strings.Contains(t, "int"):
		return "int"
	case strings.Contains(t, "decimal"), strings.Contains(t, "double"), strings.Contains(t, "float"):
		return "float64"
	case strings.Contains(t, "time"), strings.Contains(t, "date"):
		return "time.Time"
	default:
		return "string"
	}
}

func gormTag(col parser.Column) string {
	parts := []string{"column:" + col.Name, "type:" + col.Type}
	if col.Primary {
		parts = append(parts, "primaryKey")
	}
	if col.AutoIncrement {
		parts = append(parts, "autoIncrement")
	}
	if !col.Nullable {
		parts = append(parts, "not null")
	}
	if col.Unique {
		parts = append(parts, "uniqueIndex")
	}
	return strings.Join(parts, ";")
}

func esType(t string) string {
	t = strings.ToLower(t)
	switch {
	case strings.Contains(t, "bool"), strings.Contains(t, "tinyint(1)"):
		return "boolean"
	case strings.Contains(t, "int"):
		return "long"
	case strings.Contains(t, "decimal"), strings.Contains(t, "double"), strings.Contains(t, "float"):
		return "double"
	case strings.Contains(t, "time"), strings.Contains(t, "date"):
		return "date"
	case strings.Contains(t, "text"):
		return "text"
	default:
		return "keyword"
	}
}

func mongoType(t string) string {
	t = strings.ToLower(t)
	switch {
	case strings.Contains(t, "bool"), strings.Contains(t, "tinyint(1)"):
		return "bool"
	case strings.Contains(t, "int"):
		return "long"
	case strings.Contains(t, "decimal"), strings.Contains(t, "double"), strings.Contains(t, "float"):
		return "double"
	case strings.Contains(t, "time"), strings.Contains(t, "date"):
		return "date"
	default:
		return "string"
	}
}

func quoteList(items []string) string {
	quoted := make([]string, 0, len(items))
	for _, item := range items {
		quoted = append(quoted, fmt.Sprintf("%q", item))
	}
	return strings.Join(quoted, ", ")
}
