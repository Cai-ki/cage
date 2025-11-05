package helper

import "github.com/Cai-ki/cage/llm"

func JsonToSql(jstr string) (string, error) {
	prompt := `
你是一个专业的数据库建模助手。请根据用户提供的 JSON 字符串，生成一条标准的 SQL CREATE TABLE 语句。

要求如下：
1. 表名使用 "data"。
2. 自动推断每个字段的数据类型：
   - 如果值是字符串 → 使用 TEXT 或 VARCHAR(255)
   - 如果值是整数 → 使用 INTEGER
   - 如果值是浮点数 → 使用 REAL 或 DECIMAL
   - 如果值是布尔值 → 使用 BOOLEAN
   - 如果值是 null 或缺失 → 默认使用 TEXT，并在注释中标注“可能为 NULL”
3. 对于嵌套对象或数组：
   - 不要尝试创建子表；
   - 将嵌套结构扁平化，使用下划线连接路径作为字段名（例如：user.name → user_name）；
   - 如果嵌套过深或结构复杂，请在 SQL 注释中提醒用户“此处为嵌套结构，建议人工复核”。
4. 所有字段设为 NOT NULL 仅当 JSON 中该字段在所有示例中都存在且非 null；否则允许 NULL。
5. 使用通用 SQL 语法，避免使用 MySQL、PostgreSQL 等特定方言。
6. 在生成的 SQL 前后不要添加任何解释、标记或多余内容，只需输出纯 SQL 语句。
7. 如果 JSON 无效或为空，请返回错误信息：-- ERROR: 无效或空的 JSON 输入

现在，请根据以下 JSON 内容生成 SQL：

` + "`" + jstr + "`" + `
`
	return llm.Completion(prompt)
}

func JsonToGoStruct(jstr string) (string, error) {
	prompt := `
你是一个 Go 语言专家。请根据提供的 JSON 字符串，生成一个符合 Go 语言规范的结构体（struct）。

要求：
1. 结构体名为 "Data"。
2. 字段名使用大写首字母（公开字段），并根据 JSON key 合理驼峰化（如 "user_id" → UserID）。
3. 自动推断字段类型：
   - 字符串 → string
   - 整数 → int64（避免溢出）
   - 浮点数 → float64
   - 布尔 → bool
   - 数组 → []T（递归推断 T）
   - 嵌套对象 → 嵌套 struct（命名如 InnerData）
4. 为每个字段添加 json 标签，保持原始 key（如 "json:"user_id""）。
5. 不要包含包声明、import 或其他无关代码。
6. 如果 JSON 无效或为空，返回：-- ERROR: 无效或空的 JSON 输入

输入 JSON：
` + "`" + jstr + "`" + `
`
	return llm.Completion(prompt)
}

func SqlToJSONSchema(sql string) (string, error) {
	prompt := `
你是一个数据建模工程师。请将以下 SQL 的 CREATE TABLE 语句转换为标准的 JSON Schema（draft-07）。

要求：
1. 输出必须是合法的 JSON。
2. 每个字段根据 SQL 类型映射为 JSON Schema 类型（如 INTEGER → "type": "integer"）。
3. 如果字段含 NOT NULL，则 required 数组包含该字段名。
4. 表名作为 "$id" 或 "title"。
5. 忽略索引、外键、默认值等高级特性，仅关注字段名和基本类型。
6. 支持常见类型：INTEGER, BIGINT, REAL, DECIMAL, TEXT, VARCHAR, BOOLEAN, DATETIME, TIMESTAMP。
7. 若输入不是 CREATE TABLE 语句，返回 {"error": "无效的 SQL 输入"}

SQL 语句：
` + "`" + sql + "`" + `
`
	return llm.Completion(prompt)
}

func CsvToSql(csvSample string) (string, error) {
	prompt := `
你是一个数据库工程师。请根据提供的 CSV 前几行（包含表头和示例行）生成 CREATE TABLE SQL。

要求：
1. 表名为 "csv_data"。
2. 通过示例行推断字段类型（如 "123" → INTEGER，"3.14" → REAL，"2025-01-01" → TEXT 或 DATETIME）。
3. 所有字段默认允许 NULL。
4. 使用通用 SQL 语法。
5. 如果无法推断类型，使用 TEXT。
6. 输出仅为 SQL，无额外说明。

CSV 数据（第一行为列名）：
` + "`" + csvSample + "`" + `
`
	return llm.Completion(prompt)
}

func DescriptionToSql(desc string) (string, error) {
	prompt := `
你是一个智能数据库助手。请根据以下自然语言描述，生成一条 CREATE TABLE 语句。

示例描述："一张表记录用户交易，包含用户ID（字符串）、交易金额（小数）、是否成功（布尔值）、交易时间（时间戳）"

要求：
1. 表名为 "data"。
2. 合理推断字段名和类型。
3. 时间相关字段用 DATETIME 或 TIMESTAMP。
4. 金额用 DECIMAL(18,8)。
5. 输出仅为 SQL。

描述：
` + "`" + desc + "`" + `
`
	return llm.Completion(prompt)
}

func JsonToProto(jstr string) (string, error) {
	prompt := `
你是一个 protobuf 专家。请根据提供的 JSON 示例，生成一个 .proto 文件中的 message 定义。

要求：
1. message 名为 "Data"。
2. 字段编号从 1 开始连续分配。
3. 类型映射：
   - string → string
   - 整数 → int64
   - 浮点 → double
   - 布尔 → bool
   - 数组 → repeated T
   - 嵌套对象 → nested message（命名如 Data_Sub）
4. 不要包含 syntax、package、import 等，仅输出 message 块。
5. 如果 JSON 无效，返回：// ERROR: 无效输入

JSON：
` + "`" + jstr + "`" + `
`
	return llm.Completion(prompt)
}
