# 包功能说明

本包提供了一个简洁易用的 LLM（大语言模型）客户端库，基于 OpenAI 兼容的 API 设计。它封装了文本生成、视觉理解、文本嵌入等核心 AI 功能，支持通过环境变量配置 API 密钥、基础 URL 和模型参数。该包采用单例模式设计，自动初始化默认客户端，简化了调用流程，适用于需要快速集成 AI 能力的各种应用场景，如智能对话系统、图像分析工具和语义搜索应用。

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

Config 结构体用于配置 LLM 客户端的各项参数。APIKey 是访问 API 所需的密钥；BaseURL 支持兼容 OpenAI API 的自部署服务；Model 指定默认的文本生成模型；VisionModel 指定视觉分析模型；EmbedModel 指定文本嵌入模型；EmbedDim 设置嵌入向量的维度；Temperature 控制生成文本的随机性；TopP 用于核采样，控制生成文本的多样性。

```go
type LLMClient struct {
    // 包含未导出字段
}
```

LLMClient 是 LLM 客户端的主要结构体，封装了与 OpenAI 兼容 API 的交互逻辑。它内部维护配置信息和 OpenAI 客户端实例，提供文本生成、视觉分析和文本嵌入等功能的实现。

## 函数

```go
func LoadConfig() (*Config, error)
```

LoadConfig 函数从环境变量加载配置信息并返回 Config 结构体实例。它会读取 LLM_API_KEY、LLM_BASE_URL、LLM_MODEL 等环境变量，并将字符串类型的值转换为相应的数据类型。如果环境变量未设置，会使用默认值。

```go
func Completion(prompt string) (string, error)
```

Completion 函数使用默认客户端对给定的文本提示生成补全内容。它接收一个文本提示作为参数，返回生成的文本内容。内部会确保默认客户端已初始化，然后调用客户端的 completion 方法。

```go
func CompletionBySystem(prompt string) (string, error)
```

CompletionBySystem 函数使用系统消息角色生成文本内容。它接收一个系统提示作为参数，返回生成的文本内容。这种方法适用于需要系统指令引导的对话场景。

```go
func CompletionByParams(args ...AllowedParam) (openai.ChatCompletionMessage, error)
```

CompletionByParams 函数支持通过参数灵活配置聊天补全请求。它接收可变数量的 AllowedParam 参数，返回完整的聊天补全消息对象。允许组合用户消息、系统消息和工具调用等不同类型的参数。

```go
func Vision(img image.Image) (string, error)
```

Vision 函数对输入的图像进行分析并返回文本描述。它接收一个 image.Image 对象作为参数，使用默认的视觉提示"Describe this image."，返回对图像的文本描述。

```go
func VisionWithPrompt(img image.Image, prompt string) (string, error)
```

VisionWithPrompt 函数使用自定义提示对图像进行分析。它接收图像和自定义提示文本作为参数，返回根据提示分析得到的文本结果。适用于需要特定图像分析任务的场景。

```go
func Embedding(text string) ([]float32, error)
```

Embedding 函数将输入文本转换为向量表示。它接收文本字符串作为参数，返回对应的浮点数向量切片。向量维度使用配置中的默认值，适用于语义搜索和文本相似度计算。

```go
func EmbeddingWithDim(text string, dimensions int) ([]float32, error)
```

EmbeddingWithDim 函数将文本转换为指定维度的向量表示。它接收文本和期望的向量维度作为参数，返回对应维度的浮点数向量。允许覆盖配置中的默认维度设置。

```go
func UserMessage(prompt string) MessageFunc
```

UserMessage 函数创建一个用户消息的参数生成函数。它接收提示文本作为参数，返回一个 MessageFunc 函数，该函数在调用时会生成对应的用户消息参数。

```go
func SystemMessage(prompt string) MessageFunc
```

SystemMessage 函数创建一个系统消息的参数生成函数。它接收系统提示文本作为参数，返回一个 MessageFunc 函数，用于生成系统角色消息参数。

```go
func ToolMessage(prompt string, toolCallID string) MessageFunc
```

ToolMessage 函数创建一个工具消息的参数生成函数。它接收消息内容和工具调用 ID 作为参数，返回一个 MessageFunc 函数，用于生成工具调用结果消息。

```go
func ToolsByJson(jstr string) ToolFunc
```

ToolsByJson 函数从 JSON 字符串创建工具参数。它接收 JSON 格式的字符串作为参数，解析后返回一个 ToolFunc 函数，该函数提供工具调用定义供模型使用。

## 变量与常量

```go
var ErrUnexpectedResponse = errors.New("llm: unexpected API response")
```

ErrUnexpectedResponse 变量表示 API 返回了意外的响应格式。当 API 响应缺少预期的数据字段或返回空结果时返回此错误。