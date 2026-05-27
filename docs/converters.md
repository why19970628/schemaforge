# Converter Rules

SchemaForge focuses on structural conversion. It does not connect to databases, execute SQL, or generate a full application project.

## JSON to Go

Input: a JSON object sample.

Output: Go struct definitions with `json` tags.

Rules:

- Object keys become exported Go field names.
- Common initialisms such as `id`, `json`, `xml`, and `url` are capitalized as `ID`, `JSON`, `XML`, and `URL`.
- Integers become `int64`; decimals become `float64`.
- Arrays use the first item to infer element type.
- Nested objects become nested struct definitions.
- Empty arrays become `[]any`.

## YAML to Go

Input: a YAML object sample.

Output: Go struct definitions with `json` tags.

Current scope:

- Flat YAML object samples are supported.
- Scalars are inferred as `bool`, `int64`, `float64`, or `string`.
- Nested YAML can be added later without changing the CLI mode.

## XML to JSON

Input: an XML document.

Output: formatted JSON.

Rules:

- The XML root element is preserved as the top-level JSON key.
- Attributes use an `@` prefix, for example `id="1"` becomes `"@id": "1"`.
- Text-only elements become strings.
- Repeated children become JSON arrays.
- Mixed text and children store text in `"#text"`.

## SQL Converters

Input: MySQL `CREATE TABLE` DDL.

Unsupported in the first release:

- Query SQL
- Triggers
- Stored procedures
- Views
- Database connections

### SQL Type Mapping

| MySQL type | Go / GORM | Ent | Elasticsearch | MongoDB |
| --- | --- | --- | --- | --- |
| `bigint` | `int64` | `field.Int64` | `long` | `long` |
| other `int` | `int` | `field.Int` | `long` | `long` |
| `tinyint(1)`, `bool` | `bool` | `field.Bool` | `boolean` | `bool` |
| `decimal`, `double`, `float` | `float64` | `field.Float` | `double` | `double` |
| `date`, `time`, `datetime`, `timestamp` | `time.Time` | `field.Time` | `date` | `date` |
| `json` | `string` today | `field.JSON` | `keyword` | `string` |
| text-like types | `string` | `field.String` | `text` for `text`, otherwise `keyword` | `string` |

### SQL to Ent

Output: Ent schema code.

Rules:

- Table names are singularized for schema type names.
- Columns become Ent fields.
- `NOT NULL` columns are required; nullable columns use `Optional().Nillable()`.
- Column comments become `Comment(...)`.
- Unique indexes become `index.Fields(...).Unique()`.

### SQL to GORM

Output: Go model code.

Rules:

- Table names are singularized for struct names.
- Struct fields include `gorm` and `json` tags.
- `PRIMARY KEY`, `AUTO_INCREMENT`, `NOT NULL`, and unique index metadata are included in GORM tags.
- `time.Time` fields automatically add `import "time"`.

### SQL to Elasticsearch

Output: JSON mapping by table name.

Rules:

- Each table becomes an index mapping entry.
- Columns become `properties`.
- String-like fields default to `keyword`; `text` fields become `text`.

### SQL to MongoDB

Output: MongoDB JSON Schema validator by table name.

Rules:

- Each table becomes a validator object.
- `NOT NULL` columns are listed in `required`.
- Column comments become field descriptions.

## UI Behavior

- SQL modes show `Format SQL`; non-SQL modes hide it.
- SQL formatting runs in browser JavaScript and does not call a network service.
- Input blur triggers conversion when the input is non-empty and changed since the last conversion.
- Theme and language preferences are stored in localStorage.
