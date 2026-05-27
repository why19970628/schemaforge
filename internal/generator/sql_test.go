package generator

import (
	"strings"
	"testing"

	"github.com/schemaforge/schemaforge/internal/parser"
)

func TestGORMModelsAddsTimeImportOnlyWhenNeeded(t *testing.T) {
	withTime := []parser.Table{{Name: "events", Columns: []parser.Column{{Name: "created_at", Type: "datetime", Nullable: false}}}}
	out, err := GORMModels(withTime)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, `import "time"`) || !strings.Contains(out, "CreatedAt time.Time") {
		t.Fatalf("time import or field missing:\n%s", out)
	}

	withoutTime := []parser.Table{{Name: "users", Columns: []parser.Column{{Name: "email", Type: "varchar(255)", Nullable: false}}}}
	out, err = GORMModels(withoutTime)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out, `import "time"`) {
		t.Fatalf("unexpected time import:\n%s", out)
	}
}

func TestSQLTypeMapping(t *testing.T) {
	cases := []struct {
		mysql string
		goTyp string
		ent   string
		es    string
		mongo string
	}{
		{"bigint", "int64", "Int64", "long", "long"},
		{"tinyint(1)", "bool", "Bool", "boolean", "bool"},
		{"decimal(10,2)", "float64", "Float", "double", "double"},
		{"datetime", "time.Time", "Time", "date", "date"},
		{"varchar(255)", "string", "String", "keyword", "string"},
	}
	for _, tc := range cases {
		if got := sqlGoType(tc.mysql); got != tc.goTyp {
			t.Fatalf("sqlGoType(%q) = %q", tc.mysql, got)
		}
		if got := entType(tc.mysql); got != tc.ent {
			t.Fatalf("entType(%q) = %q", tc.mysql, got)
		}
		if got := esType(tc.mysql); got != tc.es {
			t.Fatalf("esType(%q) = %q", tc.mysql, got)
		}
		if got := mongoType(tc.mysql); got != tc.mongo {
			t.Fatalf("mongoType(%q) = %q", tc.mysql, got)
		}
	}
}
