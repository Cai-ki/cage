# 包功能说明

本包是一个基于大语言模型的智能转换工具集，提供多种数据格式转换和代码生成功能。通过封装复杂的提示词工程和模板处理逻辑，该包简化了常见的数据转换任务，包括 JSON 到 SQL、JSON 到 Go 结构体、SQL 到 JSON Schema、CSV 到 SQL、自然语言到 SQL、JSON 到 Protobuf 等转换功能。此外，还提供了包文档生成、Shell 脚本生成等实用工具。该包采用统一的参数接口设计，支持模板缓存和预处理机制，适用于需要自动化代码生成、数据建模和格式转换的各种应用场景。

## 结构体与接口

```go
type Param interface {
    Prompt() string
    Prepare() error
    Do() (string, error)
}
```

Param 接口定义了所有参数类型需要实现的方法，包括获取提示词模板、预处理参数数据和执行转换操作。该接口为包内的各种转换功能提供了统一的执行框架。

```go
type AnalyzePackageParam struct {
    Dir  string
    Code string
}
```

AnalyzePackageParam 用于包源码分析功能的参数结构体。Dir 字段指定要分析的源码目录路径，Code 字段在预处理后存储读取到的所有 Go 源码文件内容。

```go
type HelloParam struct {
    Input string
    Text  string
}
```

HelloParam 用于示例功能的参数结构体。Input 字段接收用户输入，Text 字段在预处理后存储格式化后的输入内容。

```go
type DescribeToShellScriptParam struct {
    Input string
    Text  string
}
```

DescribeToShellScriptParam 用于生成 Shell 脚本的参数结构体。Input 字段接收自然语言描述，Text 字段在预处理后存储处理后的描述文本。

```go
type DescribeToRunnableShellParam struct {
    Input string
    Text  string
}
```

DescribeToRunnableShellParam 用于生成可运行 Shell 命令的参数结构体。Input 字段接收自然语言描述，Text 字段在预处理后存储处理后的描述文本。

## 函数

```go
func JsonToSql(jstr string) (string, error)
```

JsonToSql 函数将 JSON 字符串转换为 SQL CREATE TABLE 语句。函数接收一个 JSON 字符串参数，自动推断字段类型并生成标准 SQL 语法，表名固定为 "data"，支持嵌套结构的扁平化处理。

```go
func JsonToGoStruct(jstr string) (string, error)
```

JsonToGoStruct 函数将 JSON 字符串转换为 Go 语言结构体定义。函数生成名为 "Data" 的结构体，自动推断字段类型并添加正确的 JSON 标签，字段名采用驼峰命名法。

```go
func SqlToJSONSchema(sql string) (string, error)
```

SqlToJSONSchema 函数将 SQL CREATE TABLE 语句转换为 JSON Schema 格式。函数基于 draft-07 标准，自动映射 SQL 类型到 JSON Schema 类型，处理 NOT NULL 约束为 required 字段。

```go
func CsvToSql(csvSample string) (string, error)
```

CsvToSql 函数根据 CSV 样本数据生成 SQL CREATE TABLE 语句。函数通过分析 CSV 表头和示例行推断字段类型，表名固定为 "csv_data"，所有字段默认允许 NULL。

```go
func DescriptionToSql(desc string) (string, error)
```

DescriptionToSql 函数根据自然语言描述生成 SQL CREATE TABLE 语句。函数解析自然语言描述，智能推断字段名和类型，表名固定为 "data"，特别处理时间戳和金额字段。

```go
func JsonToProto(jstr string) (string, error)
```

JsonToProto 函数将 JSON 字符串转换为 Protobuf message 定义。函数生成名为 "Data" 的 message，自动分配字段编号并正确映射类型，支持嵌套对象和数组类型。

```go
func AnalyzePackage(dir string) (string, error)
```

AnalyzePackage 函数分析指定目录下的 Go 源码包并生成文档。函数递归读取目录中的所有非测试 Go 文件，按文件名排序后拼接内容，通过大语言模型分析包结构和功能。

```go
func GeneratePackageDoc(dir, outputPath string) error
```

GeneratePackageDoc 函数生成包文档并写入指定文件路径。函数先调用 AnalyzePackage 获取分析结果，然后创建必要的输出目录并将文档内容写入文件。

```go
func Hello(input string) (string, error)
```

Hello 函数是一个示例功能，对输入文本进行简单处理并返回结果。函数将输入文本用竖线符号包裹后传递给大语言模型处理。

```go
func DescribeToShellScript(description string) (string, error)
```

DescribeToShellScript 函数根据自然语言描述生成安全的 Shell 脚本。函数接收任务描述，生成完整、安全的 Shell 脚本代码，包含适当的错误处理和安全检查。

```go
func DescribeToRunnableShell(description string) (string, error)
```

DescribeToRunnableShell 函数根据自然语言描述生成可直接运行的 Shell 命令片段。函数生成简洁的命令行指令，适合直接粘贴到终端执行。

```go
func Parse(param Param) (string, error)
```

Parse 函数解析参数并生成大语言模型提示词。函数首先调用参数的 Prepare 方法进行预处理，然后使用模板引擎渲染提示词内容，支持模板缓存优化性能。

```go
func Do(param Param) (string, error)
```

Do 函数执行完整的转换流程。函数内部调用 Parse 生成提示词，然后调用大语言模型的 Completion 方法获取处理结果，是包内各种转换功能的统一执行入口。

## 变量与常量

```go
var analyzePackagePrompt string
```

analyzePackagePrompt 是嵌入的包分析提示词模板内容，用于指导大语言模型如何分析 Go 源码包的结构和功能。

```go
var helloPrompt string
```

helloPrompt 是嵌入的示例功能提示词模板内容，用于 Hello 功能的处理逻辑。

```go
var describeToShellScriptPrompt string
```

describeToShellScriptPrompt 是嵌入的 Shell 脚本生成提示词模板内容，包含生成安全 Shell 脚本的详细指导要求。

```go
var describeToRunnableShellPrompt string
```

describeToRunnableShellPrompt 是嵌入的可运行 Shell 命令生成提示词模板内容，用于生成终端可直接执行的命令片段。

```go
var templateCache map[string]*template.Template
```

templateCache 是模板缓存映射，用于存储已解析的提示词模板，避免重复解析提高性能。键为参数类型名，值为对应的模板对象。