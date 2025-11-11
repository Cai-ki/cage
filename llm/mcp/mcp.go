package mcp

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/openai/openai-go"
)

// --- 1. MCPClient 结构体 ---

type MCPClient struct {
	tools map[string]*ToolExecutor
}

// ToolExecutor 用于存储工具的元信息和执行函数
type ToolExecutor struct {
	Name     string
	Func     interface{} // func(args ArgsStruct) (interface{}, error)
	ArgsType interface{} // ArgsStruct 零值
}

// --- 2. 全局默认客户端 ---

var defaultClient *MCPClient

func init() {
	defaultClient = NewMCPClient()
}

// NewMCPClient 创建一个新的 MCPClient 实例
func NewMCPClient() *MCPClient {
	return &MCPClient{
		tools: make(map[string]*ToolExecutor),
	}
}

// --- 3. 注册工具的方法 ---

// RegisterTool 注册一个工具函数
// fn 必须是 func(args ArgsStruct) (interface{}, error) 的形式
// argsType 是 ArgsStruct 的零值
func (c *MCPClient) RegisterTool(name string, fn interface{}, argsType interface{}) {
	c.tools[name] = &ToolExecutor{
		Name:     name,
		Func:     fn,
		ArgsType: argsType,
	}
}

// RegisterTool 全局函数，操作默认客户端
func RegisterTool(name string, fn interface{}, argsType interface{}) {
	defaultClient.RegisterTool(name, fn, argsType)
}

// --- 4. 执行工具调用的方法 ---

// ExecuteToolCalls 接收 OpenAI 返回的 Message，自动解析并执行 ToolCalls
func (c *MCPClient) ExecuteToolCalls(message openai.ChatCompletionMessage) ([]openai.ChatCompletionMessageParamUnion, error) {
	var results []openai.ChatCompletionMessageParamUnion

	if message.ToolCalls == nil {
		return nil, fmt.Errorf("message has no tool_calls")
	}

	for _, tc := range message.ToolCalls {
		executor, ok := c.tools[tc.Function.Name]
		if !ok {
			return nil, fmt.Errorf("unknown tool: %s", tc.Function.Name)
		}

		// 反序列化 arguments
		argsValuePtr := reflect.New(reflect.TypeOf(executor.ArgsType))
		if err := json.Unmarshal([]byte(tc.Function.Arguments), argsValuePtr.Interface()); err != nil {
			return nil, fmt.Errorf("failed to unmarshal arguments for tool %s: %w", tc.Function.Name, err)
		}

		// 调用执行函数
		fn := reflect.ValueOf(executor.Func)
		argsVal := argsValuePtr.Elem()
		fnArgs := []reflect.Value{argsVal}
		fnResults := fn.Call(fnArgs)

		// 检查返回值
		if len(fnResults) != 2 {
			return nil, fmt.Errorf("tool function must return (interface{}, error)")
		}
		result := fnResults[0].Interface()
		errValue := fnResults[1]
		var err error
		if !errValue.IsNil() {
			err = errValue.Interface().(error)
		}

		toolCallID := tc.ID

		if err != nil {
			// 构造错误的 Tool Message
			results = append(results, openai.ChatCompletionMessageParamUnion{
				OfTool: &openai.ChatCompletionToolMessageParam{
					Role: "tool",
					Content: openai.ChatCompletionToolMessageParamContentUnion{
						OfString: openai.String("Error: " + err.Error()),
					},
					ToolCallID: toolCallID,
				},
			})
			continue
		}

		// 序列化结果并构造成功的 Tool Message
		resultBytes, err := json.Marshal(result) // 处理错误
		if err != nil {
			return nil, fmt.Errorf("failed to marshal tool result for %s: %w", tc.Function.Name, err)
		}
		results = append(results, openai.ChatCompletionMessageParamUnion{
			OfTool: &openai.ChatCompletionToolMessageParam{
				Role: "tool",
				Content: openai.ChatCompletionToolMessageParamContentUnion{
					OfString: openai.String(string(resultBytes)),
				},
				ToolCallID: toolCallID,
			},
		})
	}

	return results, nil
}

// ExecuteToolCalls 全局函数，操作默认客户端
func ExecuteToolCalls(message openai.ChatCompletionMessage) ([]openai.ChatCompletionMessageParamUnion, error) {
	return defaultClient.ExecuteToolCalls(message)
}

// --- 5. 获取工具定义的方法 (用于传递给 OpenAI) ---

// GetToolsDefinition 返回注册工具的 JSON Schema 定义，可用于 OpenAI Tools 参数
func (c *MCPClient) GetToolsDefinition() ([]openai.ChatCompletionToolParam, error) {
	// 这里需要一个机制将 Go 结构体转换为 JSON Schema
	// 由于 Go 无法直接从类型推断 JSON Schema，你需要手动定义或使用反射库
	// 为简化，这里返回一个空的实现，你需要根据实际需要实现
	// 例如，可以为每个 ToolExecutor 添加一个 Schema 字段
	// 或者使用如 github.com/invopop/jsonschema 等库来生成
	// 由于这比较复杂，且与你的核心需求（执行）关系不大，这里先不实现
	// 你可以选择在注册工具时同时提供 JSON Schema 定义
	return nil, fmt.Errorf("GetToolsDefinition not implemented, please provide tools JSON manually or implement schema generation")
}

// GetToolsDefinition 全局函数
func GetToolsDefinition() ([]openai.ChatCompletionToolParam, error) {
	return defaultClient.GetToolsDefinition()
}
