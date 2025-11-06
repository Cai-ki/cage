# 包功能说明

本包提供了一个简洁易用的 LLM（大语言模型）客户端库，基于 OpenAI 兼容的 API 设计。它封装了文本生成、视觉理解和文本嵌入等核心 AI 功能，支持通过环境变量配置 API 密钥、模型参数和基础 URL，能够灵活适配不同的 LLM 服务提供商（如 OpenAI、Ollama 等）。该包采用单例模式设计，提供了开箱即用的全局函数，适合快速集成 AI 能力到各种 Go 应用程序中。

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

Config 结构体用于配置 LLM 客户端的各项参数。APIKey 是访问 LLM API 所需的密钥；BaseURL 支持设置兼容 API 的基础 URL，可用于连接 Ollama 等替代服务；Model 指定默认的文本生成模型；VisionModel 指定默认的视觉模型；EmbedModel 指定默认的嵌入模型；EmbedDim 设置嵌入向量的维度；Temperature 控制生成文本的随机性；TopP 用于核采样，影响词汇选择的多样性。

```go
type LLMClient struct {
	// 包含未导出字段
}
```

LLMClient 是 LLM 客户端的主要结构体，封装了与 LLM API 交互的核心功能。它内部维护配置信息和 OpenAI 客户端实例，但所有字段均为私有，只能通过包级函数间接使用。

## 函数

```go
func LoadConfig() (*Config, error)
```

LoadConfig 函数从环境变量加载配置并返回 Config 结构体实例。它会读取 LLM_API_KEY、LLM_BASE_URL、LLM_MODEL、LLM_VISION_MODEL、LLM_EMBED_MODEL、LLM_EMBED_DIM、LLM_TEMPERATURE 和 LLM_TOPP 等环境变量，为未设置的数值型参数提供默认值。

```go
func Completion(prompt string) (string, error)
```

Completion 函数使用默认配置的 LLM 客户端，根据给定的文本提示生成文本回复。它接收一个字符串参数作为提示，返回生成的文本内容或可能出现的错误。

```go
func CompletionBySystem(prompt string) (string, error)
```

CompletionBySystem 函数与 Completion 类似，但将提示作为系统消息发送给模型。这通常用于需要遵循特定指令或角色的场景，返回生成的文本内容或可能出现的错误。

```go
func Vision(img image.Image) (string, error)
```

Vision 函数使用默认配置的 LLM 客户端分析图像并返回文本描述。它接收一个 image.Image 接口实例，返回对图像的文本描述或可能出现的错误。

```go
func VisionWithPrompt(img image.Image, prompt string) (string, error)
```

VisionWithPrompt 函数在分析图像时允许指定自定义指令。它接收图像和文本提示两个参数，根据提示要求对图像进行分析并返回文本结果或可能出现的错误。

```go
func Embedding(text string) ([]float32, error)
```

Embedding 函数将输入文本转换为向量表示。它接收一个字符串参数，返回对应的浮点数向量切片，向量维度使用配置中的默认值。

```go
func EmbeddingWithDim(text string, dimensions int) ([]float32, error)
```

EmbeddingWithDim 函数允许指定嵌入向量的维度。它接收文本和维度两个参数，返回指定维度的文本嵌入向量，适用于需要控制向量大小的场景。

## 变量与常量

```go
var ErrUnexpectedResponse = errors.New("llm: unexpected API response")
```

ErrUnexpectedResponse 是一个错误变量，当 LLM API 返回意外的响应格式（如空数据）时返回。例如，当嵌入请求返回空数据数组或聊天完成请求返回空选择列表时，会返回此错误。