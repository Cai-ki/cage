# 包功能说明

本包是一个基于大语言模型的智能转换工具集，提供多种数据格式转换和代码生成功能。核心设计目标是通过自然语言处理技术简化开发过程中的常见任务，包括 JSON 到 SQL、Go 结构体、Protobuf 的转换，SQL 到 JSON Schema 的转换，CSV 到 SQL 的转换，以及根据自然语言描述生成 SQL 表和 Shell 脚本等。该包适用于快速原型开发、数据建模、API 设计和自动化脚本生成等场景，能够显著提升开发效率。

## 结构体与接口

```go
type Param interface {
    Prompt() string
    Prepare() error
    Do() (string, error)
}
```

Param 接口定义了所有参数类型的统一行为，包含获取提示词模板、准备数据和执行转换三个方法。所有具体的参数结构体都需要实现此接口，确保与核心处理逻辑的兼容性。

```go
type AnalyzePackageParam struct {
    Dir  string
    Code string
}
```

AnalyzePackageParam 用于包源码分析任务的参数配置。Dir 字段指定要分析的 Go 包目录路径，Code 字段在 Prepare 方法执行后会被填充为整理后的源码内容。

```go
type HelloParam struct {
    Input string
    Text  string
}
```

HelloParam 用于简单的文本处理任务。Input 字段接收原始输入字符串，Text 字段在 Prepare 方法中会被处理为带管道符号的格式化文本。

```go
type DescribeToShellScriptParam struct {
    Input string
    Text  string
}
```

DescribeToShellScriptParam 用于将自然语言描述转换为 Shell 脚本。Input 字段接收自然语言描述，Text 字段在 Prepare 方法中直接使用输入内容。

```go
type DescribeToRunnableShellParam struct {
    Input string
    Text  string
}
```

DescribeToRunnableShellParam 用于将自然语言描述转换为可运行的 Shell 命令片段。Input 字段接收自然语言描述，Text 字段在 Prepare 方法中直接使用输入内容。

## 函数

```go
func JsonToSql(jstr string) (string, error)
```

JsonToSql 将 JSON 字符串转换为标准的 SQL CREATE TABLE 语句。函数会自动推断字段数据类型，处理嵌套结构的扁平化，并生成通用的 SQL 语法。适用于快速从 JSON 样本创建数据库表结构。

```go
func JsonToGoStruct(jstr string) (string, error)
```

JsonToGoStruct 将 JSON 字符串转换为 Go 语言结构体定义。函数会自动推断字段类型，生成符合 Go 规范的公开字段和 json 标签，支持嵌套对象和数组的递归处理。适用于 API 开发和数据绑定场景。

```go
func SqlToJSONSchema(sql string) (string, error)
```

SqlToJSONSchema 将 SQL CREATE TABLE 语句转换为标准的 JSON Schema 定义。函数会根据 SQL 字段类型映射到对应的 JSON Schema 类型，处理 NOT NULL 约束，并生成符合 draft-07 规范的 JSON 输出。适用于 API 文档生成和数据验证。

```go
func CsvToSql(csvSample string) (string, error)
```

CsvToSql 根据 CSV 样本数据（包含表头和示例行）生成 CREATE TABLE SQL 语句。函数会通过示例行推断字段类型，默认所有字段允许 NULL，使用通用 SQL 语法。适用于从 CSV 文件快速创建数据库表。

```go
func DescriptionToSql(desc string) (string, error)
```

DescriptionToSql 根据自然语言描述生成 CREATE TABLE 语句。函数会合理推断字段名和类型，对时间字段使用 DATETIME 或 TIMESTAMP，对金额字段使用 DECIMAL(18,8)。适用于快速原型设计和概念验证。

```go
func JsonToProto(jstr string) (string, error)
```

JsonToProto 将 JSON 示例转换为 Protobuf message 定义。函数会自动映射数据类型，分配连续的字段编号，支持嵌套消息和重复字段。适用于 gRPC 服务开发和协议定义。

```go
func AnalyzePackage(dir string) (string, error)
```

AnalyzePackage 分析指定目录下的 Go 包源码，生成包的功能说明文档。函数会递归读取所有非测试的 .go 文件，按文件名排序后提交给大语言模型进行分析。适用于自动化文档生成和代码理解。

```go
func GeneratePackageDoc(dir, outputPath string) error
```

GeneratePackageDoc 分析包源码并生成文档文件。函数会创建必要的输出目录，将分析结果写入指定路径。适用于集成到构建流程或文档工具链中。

```go
func Hello(input string) (string, error)
```

Hello 是一个简单的示例函数，对输入文本进行格式化处理并返回结果。主要用于演示包的基本使用方式和参数处理流程。

```go
func DescribeToShellScript(description string) (string, error)
```

DescribeToShellScript 根据自然语言描述生成安全的 Shell 脚本。函数会生成完整的脚本文件，包含适当的注释和安全检查。适用于自动化部署和系统管理任务。

```go
func DescribeToRunnableShell(description string) (string, error)
```

DescribeToRunnableShell 根据自然语言描述生成可直接运行的 Shell 命令片段。与 DescribeToShellScript 不同，此函数生成的是适合直接粘贴到终端执行的命令序列。适用于快速执行一次性任务。

```go
func Parse(param Param) (string, error)
```

Parse 方法根据参数类型解析对应的提示词模板，并执行数据准备和模板渲染。函数内部维护模板缓存以提高性能，确保相同类型的参数重复使用时不会重复解析模板。

```go
func Do(param Param) (string, error)
```

Do 方法是执行转换任务的核心入口，封装了模板解析和大语言模型调用的完整流程。函数首先调用 Parse 生成提示词，然后通过系统提示方式调用大语言模型完成实际转换任务。

## 变量与常量

本包未定义公开的变量或常量。