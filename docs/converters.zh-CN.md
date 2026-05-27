# 转换规则

SchemaForge 专注结构转换，不连接数据库、不执行 SQL，也不生成完整应用工程。

## JSON 转 Go

输入：JSON object 样例。

输出：带 `json` tag 的 Go struct。

规则：

- object key 会转换成导出的 Go 字段名。
- 常见缩写如 `id`、`json`、`xml`、`url` 会转换成 `ID`、`JSON`、`XML`、`URL`。
- 整数推断为 `int64`，小数推断为 `float64`。
- 数组使用第一个元素推断元素类型。
- 嵌套 object 会生成嵌套 struct。
- 空数组推断为 `[]any`。

## YAML 转 Go

输入：YAML object 样例。

输出：带 `json` tag 的 Go struct。

当前范围：

- 支持常见扁平 YAML object。
- 标量会推断为 `bool`、`int64`、`float64` 或 `string`。
- 嵌套 YAML 后续可以在不改变 CLI 模式的前提下增强。

## XML 转 JSON

输入：XML 文档。

输出：格式化 JSON。

规则：

- XML 根节点会保留为 JSON 顶层 key。
- 属性使用 `@` 前缀，例如 `id="1"` 会变成 `"@id": "1"`。
- 纯文本节点会变成字符串。
- 重复子节点会变成 JSON 数组。
- 混合文本和子节点时，文本存放在 `"#text"`。

## SQL 转换器

输入：MySQL `CREATE TABLE` DDL。

首版暂不支持：

- 查询 SQL
- 触发器
- 存储过程
- 视图
- 数据库连接

### SQL 类型映射

| MySQL 类型 | Go / GORM | Ent | Elasticsearch | MongoDB |
| --- | --- | --- | --- | --- |
| `bigint` | `int64` | `field.Int64` | `long` | `long` |
| 其他 `int` | `int` | `field.Int` | `long` | `long` |
| `tinyint(1)`, `bool` | `bool` | `field.Bool` | `boolean` | `bool` |
| `decimal`, `double`, `float` | `float64` | `field.Float` | `double` | `double` |
| `date`, `time`, `datetime`, `timestamp` | `time.Time` | `field.Time` | `date` | `date` |
| `json` | 当前为 `string` | `field.JSON` | `keyword` | `string` |
| 文本类型 | `string` | `field.String` | `text` 类型为 `text`，其他为 `keyword` | `string` |

### SQL 转 Ent

输出：Ent schema 代码。

规则：

- 表名会单数化为 schema 类型名。
- 字段会生成 Ent field。
- `NOT NULL` 字段为必填；nullable 字段使用 `Optional().Nillable()`。
- 字段注释会生成 `Comment(...)`。
- 唯一索引会生成 `index.Fields(...).Unique()`。

### SQL 转 GORM

输出：Go model 代码。

规则：

- 表名会单数化为 struct 名。
- struct 字段包含 `gorm` 和 `json` tag。
- `PRIMARY KEY`、`AUTO_INCREMENT`、`NOT NULL` 和唯一索引信息会写入 GORM tag。
- 出现 `time.Time` 字段时会自动添加 `import "time"`。

### SQL 转 Elasticsearch

输出：按表名组织的 JSON mapping。

规则：

- 每张表生成一个 index mapping。
- 字段生成到 `properties`。
- 字符串字段默认 `keyword`；`text` 类型为 `text`。

### SQL 转 MongoDB

输出：按表名组织的 MongoDB JSON Schema validator。

规则：

- 每张表生成一个 validator。
- `NOT NULL` 字段会加入 `required`。
- 字段注释会作为 description。

## UI 行为

- SQL 模式显示 `Format SQL`，非 SQL 模式隐藏。
- SQL 格式化在浏览器 JavaScript 中完成，不调用网络服务。
- 输入框失焦时，如果输入非空且相对上次转换有变化，会自动转换。
- 主题和语言偏好保存在 localStorage。
