# 包功能说明

本包是一个基于大语言模型的智能代码生成与文档生成工具集，主要提供多种格式转换和自动化文档生成功能。包的核心设计目标是通过自然语言处理技术简化开发工作流，包括 JSON 到 SQL/Go 结构体/Protobuf 的转换、CSV 到 SQL 的转换、自然语言到 SQL 的转换，以及 Shell 脚本生成等。该包适用于快速原型开发、数据库设计、API 协议定义和自动化脚本编写等场景，能够显著提升开发效率。

## 结构体与接口

本包未定义可被外部包访问的结构体与接口。

## 函数

```go
func JsonToSql(jstr string) (string, error)
```

将 JSON 字符串转换为标准的 SQL CREATE TABLE 语句。函数会自动推断字段数据类型，处理嵌套结构的扁平化，并生成通用的 SQL 语法。参数 jstr 为输入的 JSON 字符串，返回生成的 SQL 语句或错误信息。

```go
func JsonToGoStruct(jstr string) (string, error)
```

将 JSON 字符串转换为符合 Go 语言规范的结构体定义。函数会自动推断字段类型，生成驼峰命名的公开字段，并为每个字段添加 json 标签。参数 jstr 为输入的 JSON 字符串，返回生成的 Go 结构体代码或错误信息。

```go
func SqlToJSONSchema(sql string) (string, error)
```

将 SQL 的 CREATE TABLE 语句转换为标准的 JSON Schema（draft-07 版本）。函数会根据 SQL 字段类型映射到对应的 JSON Schema 类型，处理 NOT NULL 约束，并生成合法的 JSON 输出。参数 sql 为输入的 SQL 语句，返回生成的 JSON Schema 或错误信息。

```go
func CsvToSql(csvSample string) (string, error)
```

根据 CSV 样本数据（包含表头和示例行）生成 CREATE TABLE SQL 语句。函数会通过示例行推断字段类型，默认所有字段允许 NULL，并使用通用 SQL 语法。参数 csvSample 为 CSV 格式的字符串样本，返回生成的 SQL 语句。

```go
func DescriptionToSql(desc string) (string, error)
```

根据自然语言描述生成 CREATE TABLE 语句。函数会合理推断字段名和类型，对金额字段使用 DECIMAL(18,8)，时间相关字段使用 DATETIME 或 TIMESTAMP。参数 desc 为自然语言描述文本，返回生成的 SQL 语句。

```go
func JsonToProto(jstr string) (string, error)
```

将 JSON 示例转换为 Protobuf 的 message 定义。函数会自动映射类型，分配字段编号，处理嵌套对象和数组类型。参数 jstr 为输入的 JSON 字符串，返回生成的 Protobuf message 定义或错误信息。

```go
func AnalyzePackageBySourceCode(dir string) (string, error)
```

分析指定目录下的 Go 源代码包，生成格式化的 Markdown 文档。函数会递归读取目录中的所有非测试 Go 文件，通过大模型分析包的公开接口和功能。参数 dir 为源码目录路径，返回生成的 Markdown 文档内容。

```go
func GeneratePackageDoc(dir, outputPath string) error
```

生成 Go 包的文档并写入指定文件。函数会调用 AnalyzePackageBySourceCode 分析源码，然后将结果写入指定的输出路径。参数 dir 为源码目录路径，outputPath 为输出文件路径，返回可能的错误信息。

```go
func DescribeToShellScript(description string) (string, error)
```

根据自然语言描述生成安全、可直接执行的 Shell 脚本。生成的脚本包含严格模式设置、变量引用保护、错误处理和中文注释。参数 description 为自然语言描述，返回生成的完整 Shell 脚本内容。

```go
func DescribeToRunnableShell(description string) (string, error)
```

根据自然语言描述生成可直接粘贴到终端运行的 Shell 命令片段。输出为纯 Shell 命令，包含错误处理和中文注释，适用于快速命令行操作。参数 description 为自然语言描述，返回可执行的 Shell 命令字符串。

## 变量与常量

本包未定义可被外部包访问的变量与常量。