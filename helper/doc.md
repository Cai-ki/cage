# 包功能说明

本包是一个基于大语言模型的智能代码生成与文档生成工具集，主要提供多种格式转换和自动化文档生成功能。该包通过封装大语言模型的自然语言处理能力，实现了从 JSON、SQL、CSV 等不同格式之间的智能转换，以及 Go 代码包的自动化文档生成。典型使用场景包括数据库建模辅助、协议缓冲区定义生成、Shell 脚本自动生成以及 Go 项目文档自动化生成等，能够显著提升开发效率并减少手动编写重复代码的工作量。

## 结构体与接口

本包未定义公开的结构体与接口。

## 函数

```go
func JsonToSql(jstr string) (string, error)
```

将 JSON 字符串转换为标准的 SQL CREATE TABLE 语句。函数接收一个 JSON 字符串作为输入，通过大语言模型智能推断字段类型和表结构，生成适用于通用 SQL 语法的建表语句。表名固定为 "data"，会自动处理嵌套结构的扁平化映射，并根据字段存在情况设置 NULL 约束。

```go
func JsonToGoStruct(jstr string) (string, error)
```

将 JSON 字符串转换为 Go 语言结构体定义。函数根据输入的 JSON 字符串自动生成对应的 Go struct，结构体名为 "Data"，字段名采用驼峰命名法并包含正确的 json 标签。支持递归推断嵌套对象和数组类型，确保生成的代码符合 Go 语言规范。

```go
func SqlToJSONSchema(sql string) (string, error)
```

将 SQL CREATE TABLE 语句转换为 JSON Schema 格式。函数解析输入的 SQL 建表语句，根据字段类型映射到 JSON Schema 的相应类型，支持 draft-07 标准。会自动处理 NOT NULL 约束并将其转换为 required 字段，表名会作为 schema 的标题或标识符。

```go
func CsvToSql(csvSample string) (string, error)
```

根据 CSV 样本数据生成 SQL CREATE TABLE 语句。函数接收包含表头和示例行的 CSV 数据，通过分析示例值推断各列的字段类型。表名固定为 "csv_data"，所有字段默认允许 NULL，当类型推断不确定时默认使用 TEXT 类型。

```go
func DescriptionToSql(desc string) (string, error)
```

根据自然语言描述生成 SQL CREATE TABLE 语句。函数解析用户对表结构的自然语言描述，智能推断字段名和数据类型。时间字段使用 DATETIME 或 TIMESTAMP 类型，金额字段使用 DECIMAL(18,8) 精度，生成的表名为 "data"。

```go
func JsonToProto(jstr string) (string, error)
```

将 JSON 字符串转换为 Protocol Buffers 的 message 定义。函数根据输入的 JSON 示例生成对应的 protobuf message，message 名为 "Data"。字段编号从 1 开始连续分配，支持嵌套消息定义和重复字段，输出仅包含 message 块而不包含语法声明等额外内容。

```go
func AnalyzePackageBySourceCode(dir string) (string, error)
```

递归分析指定目录下的 Go 源代码并生成包功能文档。函数会遍历目录中的所有非测试 Go 文件，提取源代码内容后通过大语言模型分析包的可导出功能。返回包含包功能说明的文档字符串，适用于自动化文档生成流程。

```go
func GeneratePackageDoc(dir, outputPath string) error
```

生成完整的包文档并写入指定文件。函数首先调用 AnalyzePackageBySourceCode 分析源代码，然后将生成的文档写入指定的输出路径。会自动创建输出目录（如果不存在），提供一站式的包文档生成解决方案。

```go
func DescribeToShellScript(description string) (string, error)
```

根据自然语言描述生成安全的 Shell 脚本。函数接收对脚本功能的描述，生成符合 DevOps 最佳实践的 Bash 脚本。脚本包含严格模式设置、变量引用安全处理、错误检查和中文注释，确保生成的操作安全可靠。

```go
func DescribeToRunnableShell(description string) (string, error)
```

根据自然语言描述生成可直接在终端运行的 Shell 命令片段。与 DescribeToShellScript 不同，此函数生成的是适合直接粘贴到终端执行的命令序列，不包含脚本文件头。命令会进行安全优化，避免系统污染和危险操作，适合快速执行临时任务。

## 变量与常量

本包未定义公开的变量与常量。