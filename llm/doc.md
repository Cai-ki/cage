# 包功能说明

本包提供了一个基于 OpenAI API 的 LLM（大语言模型）客户端封装，支持文本生成、视觉分析、文本嵌入等多种 AI 功能。设计目标是简化 AI 服务的集成过程，通过环境变量配置和默认客户端模式降低使用门槛。该包支持兼容 OpenAI API 的第三方服务（如 Ollama），并提供了灵活的参数配置和工具调用机制，适用于需要 AI 能力的各种应用场景，如智能对话、图像分析、语义搜索等。

## 结构体与接口

```go
type Config struct {
    APIKey      string
    BaseURL     string
    Model       string
    VisionModel string
    EmbedModel  string
    EmbedDim    int
    Temperature float64
    TopP        float64
}
```

Config 结构体用于配置 LLM 客户端参数。APIKey 是 API 访问密钥；BaseURL 支持兼容 API 的服务地址；Model 指定默认文本模型；VisionModel 指定默认视觉模型；EmbedModel 指定默认嵌入模型；EmbedDim 设置嵌入向量的维度；Temperature 控制生成文本的随机性；TopP 用于核采样，控制生成文本的多样性。

```go
type LLMClient struct {
    // 未导出字段
}
```

LLMClient 是 LLM 客户端的主要结构体，封装了与 OpenAI API 的交互逻辑。它通过内部配置和 OpenAI 客户端实例提供各种 AI 功能。

```go
type MCPClient struct {
    // 未导出字段
}
```

MCPClient 是模型控制协议客户端，用于管理和执行工具调用。它维护了一个工具注册表，能够解析和执行来自 AI 模型的工具调用请求。

```go
type ToolExecutor struct {
    Name     string
    Func     interface{}
    ArgsType interface{}
}
```

ToolExecutor 用于存储工具的元信息和执行函数。Name 是工具名称；Func 是工具执行函数，必须符合 func(args ArgsStruct) (interface{}, error) 的签名；ArgsType 是参数结构的零值，用于 JSON 反序列化。

## 函数

```go
func LoadConfig() (*Config, error)
```

LoadConfig 从环境变量加载配置并返回 Config 实例。它会读取 LLM_API_KEY、LLM_BASE_URL、LLM_MODEL 等环境变量，为未设置的数值型参数提供默认值。

```go
func Completion(prompt string) (string, error)
```

Completion 使用默认客户端根据文本提示生成文本回复。它接收用户输入的提示文本，返回 AI 生成的文本内容。

```go
func CompletionBySystem(prompt string) (string, error)
```

CompletionBySystem 使用系统角色消息生成文本回复。适合需要系统指令的场景，如设定 AI 行为模式。

```go
func CompletionByParams(args ...AllowedParam) (openai.ChatCompletionMessage, error)
```

CompletionByParams 支持灵活的参数组合生成文本回复。可以接收消息函数和工具函数等多种参数类型，返回完整的聊天完成消息。

```go
func Vision(img image.Image) (string, error)
```

Vision 分析图像并返回文本描述。使用默认的视觉提示词，适合基础的图像理解任务。

```go
func VisionWithPrompt(img image.Image, prompt string) (string, error)
```

VisionWithPrompt 使用自定义提示词分析图像。允许指定具体的分析要求，如图像中特定内容的识别或描述。

```go
func Embedding(text string) ([]float32, error)
```

Embedding 返回输入文本的向量表示。使用配置中指定的嵌入模型和维度，返回 float32 类型的向量数组。

```go
func EmbeddingWithDim(text string, dimensions int) ([]float32, error)
```

EmbeddingWithDim 返回指定维度的文本嵌入向量。允许覆盖配置中的默认维度设置，适合需要特定向量大小的场景。

```go
func UserMessage(prompt string) MessageFunc
```

UserMessage 创建用户角色消息的函数。返回的 MessageFunc 可用于 CompletionByParams 参数。

```go
func SystemMessage(prompt string) MessageFunc
```

SystemMessage 创建系统角色消息的函数。用于设定 AI 的系统指令或行为约束。

```go
func ToolMessage(prompt string, toolCallID string) MessageFunc
```

ToolMessage 创建工具角色消息的函数。用于向 AI 返回工具调用的执行结果。

```go
func ToolsByJson(jstr string) ToolFunc
```

ToolsByJson 从 JSON 字符串创建工具定义函数。将 JSON 格式的工具定义转换为 OpenAI 工具参数。

```go
func NewMCPClient() *MCPClient
```

NewMCPClient 创建新的 MCP 客户端实例。返回一个空的工具注册表，用于管理工具调用。

```go
func RegisterTool(name string, fn interface{}, argsType interface{})
```

RegisterTool 向默认 MCP 客户端注册工具函数。name 是工具名称，fn 是工具执行函数，argsType 是参数结构的零值。

```go
func ExecuteToolCalls(message openai.ChatCompletionMessage) ([]openai.ChatCompletionMessageParamUnion, error)
```

ExecuteToolCalls 执行 OpenAI 返回消息中的工具调用。自动解析工具参数，调用注册的工具函数，并返回工具执行结果的消息数组。

```go
func GetToolsDefinition() ([]openai.ChatCompletionToolParam, error)
```

GetToolsDefinition 返回注册工具的 JSON Schema 定义。目前未实现完整功能，需要手动提供工具定义或实现 schema 生成逻辑。

## 变量与常量

```go
var ErrUnexpectedResponse = errors.New("llm: unexpected API response")
```

ErrUnexpectedResponse 在 API 返回意外响应时返回，如空数据数组或缺失必要字段。