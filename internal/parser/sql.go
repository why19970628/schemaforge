package parser

import (
	"fmt"
	"regexp"
	"strings"
)

type Table struct {
	Name    string
	Comment string
	Columns []Column
	Indexes []Index
}

type Column struct {
	Name          string
	Type          string
	Nullable      bool
	Primary       bool
	Unique        bool
	AutoIncrement bool
	Default       string
	Comment       string
}

type Index struct {
	Name    string
	Columns []string
	Unique  bool
}

func ParseMySQLDDL(input string) ([]Table, error) {
	statements := extractCreateTables(input)
	if len(statements) == 0 {
		return nil, fmt.Errorf("no CREATE TABLE statement found")
	}

	tables := make([]Table, 0, len(statements))
	for _, stmt := range statements {
		table := Table{Name: cleanIdent(stmt.name), Comment: tableComment(stmt.options)}
		lines := splitDefinitions(stmt.body)
		for _, line := range lines {
			line = strings.TrimSpace(strings.TrimSuffix(line, ","))
			upper := strings.ToUpper(line)
			switch {
			case strings.HasPrefix(line, "`"):
				col, err := parseColumn(line)
				if err != nil {
					return nil, err
				}
				table.Columns = append(table.Columns, col)
			case strings.HasPrefix(upper, "PRIMARY KEY"):
				cols := parseIndexColumns(line)
				for i := range table.Columns {
					for _, name := range cols {
						if table.Columns[i].Name == name {
							table.Columns[i].Primary = true
							table.Columns[i].Nullable = false
						}
					}
				}
			case strings.HasPrefix(upper, "UNIQUE KEY"), strings.HasPrefix(upper, "UNIQUE INDEX"):
				idx := parseIndex(line, true)
				table.Indexes = append(table.Indexes, idx)
				for i := range table.Columns {
					for _, name := range idx.Columns {
						if table.Columns[i].Name == name {
							table.Columns[i].Unique = true
						}
					}
				}
			case strings.HasPrefix(upper, "KEY"), strings.HasPrefix(upper, "INDEX"):
				table.Indexes = append(table.Indexes, parseIndex(line, false))
			}
		}
		tables = append(tables, table)
	}
	return tables, nil
}

type createTableStatement struct {
	name    string
	body    string
	options string
}

func extractCreateTables(input string) []createTableStatement {
	re := regexp.MustCompile(`(?is)CREATE\s+TABLE\s+(?:IF\s+NOT\s+EXISTS\s+)?([` + "`" + `\w.]+)\s*\(`)
	matches := re.FindAllStringSubmatchIndex(input, -1)
	var statements []createTableStatement
	for _, m := range matches {
		name := input[m[2]:m[3]]
		open := m[1] - 1
		close := matchingParen(input, open)
		if close < 0 {
			continue
		}
		end := close + 1
		for end < len(input) && input[end] != ';' {
			end++
		}
		statements = append(statements, createTableStatement{
			name:    name,
			body:    input[open+1 : close],
			options: input[close+1 : end],
		})
	}
	return statements
}

func matchingParen(input string, open int) int {
	depth := 0
	inQuote := byte(0)
	for i := open; i < len(input); i++ {
		ch := input[i]
		if inQuote != 0 {
			if ch == inQuote && (i == 0 || input[i-1] != '\\') {
				inQuote = 0
			}
			continue
		}
		switch ch {
		case '\'', '"', '`':
			inQuote = ch
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}

func splitDefinitions(body string) []string {
	var parts []string
	var current strings.Builder
	depth := 0
	inQuote := byte(0)
	for i := 0; i < len(body); i++ {
		ch := body[i]
		if inQuote != 0 {
			current.WriteByte(ch)
			if ch == inQuote && (i == 0 || body[i-1] != '\\') {
				inQuote = 0
			}
			continue
		}
		switch ch {
		case '\'', '"', '`':
			inQuote = ch
			current.WriteByte(ch)
		case '(':
			depth++
			current.WriteByte(ch)
		case ')':
			depth--
			current.WriteByte(ch)
		case ',':
			if depth == 0 {
				parts = append(parts, current.String())
				current.Reset()
			} else {
				current.WriteByte(ch)
			}
		default:
			current.WriteByte(ch)
		}
	}
	if strings.TrimSpace(current.String()) != "" {
		parts = append(parts, current.String())
	}
	return parts
}

func parseColumn(line string) (Column, error) {
	end := strings.Index(line[1:], "`")
	if end < 0 {
		return Column{}, fmt.Errorf("invalid column definition %q", line)
	}
	name := line[1 : end+1]
	rest := strings.TrimSpace(line[end+2:])
	fields := strings.Fields(rest)
	if len(fields) == 0 {
		return Column{}, fmt.Errorf("missing type for column %q", name)
	}
	colType := fields[0]
	upper := strings.ToUpper(rest)
	return Column{
		Name:          name,
		Type:          colType,
		Nullable:      !strings.Contains(upper, "NOT NULL"),
		AutoIncrement: strings.Contains(upper, "AUTO_INCREMENT"),
		Default:       captureQuotedOrToken(rest, "DEFAULT"),
		Comment:       captureQuotedOrToken(rest, "COMMENT"),
	}, nil
}

func parseIndex(line string, unique bool) Index {
	name := "idx"
	nameRe := regexp.MustCompile("(?i)(?:UNIQUE\\s+)?(?:KEY|INDEX)\\s+`?([\\w_]+)`?")
	if m := nameRe.FindStringSubmatch(line); len(m) > 1 {
		name = m[1]
	}
	return Index{Name: name, Columns: parseIndexColumns(line), Unique: unique}
}

func parseIndexColumns(line string) []string {
	start := strings.Index(line, "(")
	end := strings.LastIndex(line, ")")
	if start < 0 || end < start {
		return nil
	}
	raw := strings.Split(line[start+1:end], ",")
	cols := make([]string, 0, len(raw))
	for _, item := range raw {
		item = strings.TrimSpace(item)
		item = strings.Trim(item, "` ")
		if n := strings.Index(item, "("); n >= 0 {
			item = item[:n]
		}
		if item != "" {
			cols = append(cols, item)
		}
	}
	return cols
}

func captureQuotedOrToken(input, keyword string) string {
	re := regexp.MustCompile(`(?i)` + keyword + `\s+(?:'([^']*)'|([^\s,]+))`)
	if m := re.FindStringSubmatch(input); len(m) > 0 {
		if m[1] != "" {
			return m[1]
		}
		return m[2]
	}
	return ""
}

func tableComment(options string) string {
	re := regexp.MustCompile(`(?i)COMMENT\s*=\s*'([^']*)'`)
	if m := re.FindStringSubmatch(options); len(m) > 1 {
		return m[1]
	}
	return ""
}

func cleanIdent(s string) string {
	s = strings.TrimSpace(s)
	s = strings.Trim(s, "`")
	if dot := strings.LastIndex(s, "."); dot >= 0 {
		s = strings.Trim(s[dot+1:], "`")
	}
	return s
}
