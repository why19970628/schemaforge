# Examples

Run examples from the repository root:

```bash
go run ./cmd/schemaforge convert json-go -i examples/json/user.json
go run ./cmd/schemaforge convert yaml-go -i examples/yaml/user.yaml
go run ./cmd/schemaforge convert xml-json -i examples/xml/user.xml
go run ./cmd/schemaforge convert sql-ent -i examples/sql/users.sql
go run ./cmd/schemaforge convert sql-gorm -i examples/sql/users.sql
go run ./cmd/schemaforge convert sql-es -i examples/sql/users.sql
go run ./cmd/schemaforge convert sql-mongo -i examples/sql/users.sql
```
