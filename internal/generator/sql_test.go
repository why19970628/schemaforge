package generator

import (
	"strings"
	"testing"

	"github.com/why19970628/schemaforge/internal/parser"
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

func TestGORMModelsWithOptions(t *testing.T) {
	tables := []parser.Table{{
		Name: "users",
		Columns: []parser.Column{
			{Name: "id", Type: "bigint unsigned", Nullable: false, Primary: true, AutoIncrement: true, Comment: "primary id"},
			{Name: "nickname", Type: "varchar(64)", Nullable: true, Comment: "display name"},
			{Name: "deleted_at", Type: "datetime", Nullable: true, Comment: "deleted time"},
			{Name: "payload", Type: "json", Nullable: true, Comment: "raw payload"},
		},
	}}
	out, err := GORMModelsWithOptions(tables, Options{
		PackageName:      "model",
		NullableStrategy: NullablePointer,
		JSONTags:         false,
		Comments:         true,
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"package model",
		"// primary id",
		"ID uint64",
		"Nickname *string",
		"DeletedAt *time.Time",
		"Payload *json.RawMessage",
		`gorm:"column:nickname;type:varchar(64)"`,
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("missing %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, `json:"`) {
		t.Fatalf("did not expect json tags:\n%s", out)
	}
}

func TestGORMModelsWithSQLNullStrategy(t *testing.T) {
	tables := []parser.Table{{Name: "users", Columns: []parser.Column{
		{Name: "name", Type: "varchar(64)", Nullable: true},
		{Name: "age", Type: "int", Nullable: true},
		{Name: "enabled", Type: "tinyint(1)", Nullable: true},
		{Name: "created_at", Type: "datetime", Nullable: true},
	}}}
	out, err := GORMModelsWithOptions(tables, Options{NullableStrategy: NullableSQLNull, JSONTags: true})
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"database/sql", "Name", "sql.NullString", "Age", "sql.NullInt64", "Enabled", "sql.NullBool", "CreatedAt", "sql.NullTime"} {
		if !strings.Contains(out, want) {
			t.Fatalf("missing %q:\n%s", want, out)
		}
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
