package parser

import "testing"

func TestParseMySQLDDLColumnsAndIndexes(t *testing.T) {
	input := "CREATE TABLE `users` (`id` bigint NOT NULL AUTO_INCREMENT COMMENT 'primary id', `email` varchar(255) NOT NULL COMMENT 'email', `age` int DEFAULT 0, PRIMARY KEY (`id`), UNIQUE KEY `uk_email` (`email`), KEY `idx_age` (`age`)) COMMENT='users';"

	tables, err := ParseMySQLDDL(input)
	if err != nil {
		t.Fatal(err)
	}
	if len(tables) != 1 {
		t.Fatalf("len(tables) = %d", len(tables))
	}
	table := tables[0]
	if table.Name != "users" || table.Comment != "users" {
		t.Fatalf("unexpected table: %#v", table)
	}
	if len(table.Columns) != 3 {
		t.Fatalf("len(columns) = %d", len(table.Columns))
	}
	if !table.Columns[0].Primary || !table.Columns[0].AutoIncrement || table.Columns[0].Nullable {
		t.Fatalf("id column flags not parsed: %#v", table.Columns[0])
	}
	if !table.Columns[1].Unique || table.Columns[1].Comment != "email" {
		t.Fatalf("email column metadata not parsed: %#v", table.Columns[1])
	}
	if len(table.Indexes) != 2 {
		t.Fatalf("len(indexes) = %d", len(table.Indexes))
	}
	if !table.Indexes[0].Unique || table.Indexes[0].Name != "uk_email" || table.Indexes[0].Columns[0] != "email" {
		t.Fatalf("unique index not parsed: %#v", table.Indexes[0])
	}
}

func TestParseMySQLDDLRejectsNonCreateTable(t *testing.T) {
	if _, err := ParseMySQLDDL("select * from users"); err == nil {
		t.Fatal("expected error")
	}
}
