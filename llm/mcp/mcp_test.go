package mcp

import (
	"testing"

	"github.com/openai/openai-go"
)

// 定义测试用的参数结构体
type AddArgs struct {
	A float64 `json:"a"`
	B float64 `json:"b"`
}

type MultiplyArgs struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// 定义测试用的函数
func addFunc(args AddArgs) (interface{}, error) {
	return map[string]interface{}{"result": args.A + args.B}, nil
}

func multiplyFunc(args MultiplyArgs) (interface{}, error) {
	return map[string]interface{}{"result": args.X * args.Y}, nil
}

func failingFunc(args AddArgs) (interface{}, error) {
	return nil, &CustomError{"this function always fails"}
}

type CustomError struct {
	Message string
}

func (e *CustomError) Error() string {
	return e.Message
}

func TestMCPClient_RegisterTool(t *testing.T) {
	client := NewMCPClient()

	client.RegisterTool("add", addFunc, AddArgs{})
	client.RegisterTool("multiply", multiplyFunc, MultiplyArgs{})

	if _, exists := client.tools["add"]; !exists {
		t.Error("Failed to register 'add' tool")
	}
	if _, exists := client.tools["multiply"]; !exists {
		t.Error("Failed to register 'multiply' tool")
	}
}

func TestMCPClient_ExecuteToolCalls_Success(t *testing.T) {
	client := NewMCPClient()
	client.RegisterTool("add", addFunc, AddArgs{})
	client.RegisterTool("multiply", multiplyFunc, MultiplyArgs{})

	// 模拟 OpenAI 返回的 Message，包含 ToolCalls
	message := openai.ChatCompletionMessage{
		ToolCalls: []openai.ChatCompletionMessageToolCall{
			{
				ID:   "call_1",
				Type: "function",
				Function: openai.ChatCompletionMessageToolCallFunction{
					Name:      "add",
					Arguments: `{"a": 2, "b": 3}`,
				},
			},
			{
				ID:   "call_2",
				Type: "function",
				Function: openai.ChatCompletionMessageToolCallFunction{
					Name:      "multiply",
					Arguments: `{"x": 4, "y": 5}`,
				},
			},
		},
	}

	results, err := client.ExecuteToolCalls(message)
	if err != nil {
		t.Fatalf("ExecuteToolCalls failed: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}

	// 检查第一个结果 (add)
	if results[0].OfTool == nil {
		t.Fatal("First result is not a Tool message")
	}
	expectedAddResult := `{"result":5}`
	if results[0].OfTool.Content.OfString != openai.String(expectedAddResult) {
		t.Errorf("Expected add result '%s', got '%s'", expectedAddResult, results[0].OfTool.Content.OfString)
	}

	// 检查第二个结果 (multiply)
	if results[1].OfTool == nil {
		t.Fatal("Second result is not a Tool message")
	}
	expectedMultiplyResult := `{"result":20}`
	if results[1].OfTool.Content.OfString != openai.String(expectedMultiplyResult) {
		t.Errorf("Expected multiply result '%s', got '%s'", expectedMultiplyResult, results[1].OfTool.Content.OfString)
	}
}

func TestMCPClient_ExecuteToolCalls_Failure(t *testing.T) {
	client := NewMCPClient()
	client.RegisterTool("failing", failingFunc, AddArgs{})

	message := openai.ChatCompletionMessage{
		ToolCalls: []openai.ChatCompletionMessageToolCall{
			{
				ID:   "call_fail",
				Type: "function",
				Function: openai.ChatCompletionMessageToolCallFunction{
					Name:      "failing",
					Arguments: `{"a": 1, "b": 1}`,
				},
			},
		},
	}

	results, err := client.ExecuteToolCalls(message)
	if err != nil {
		t.Fatalf("ExecuteToolCalls failed unexpectedly: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	if results[0].OfTool == nil {
		t.Fatal("Result is not a Tool message")
	}

	// 检查错误消息
	expectedErrorMsg := "Error: this function always fails"
	if results[0].OfTool.Content.OfString != openai.String(expectedErrorMsg) {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMsg, results[0].OfTool.Content.OfString)
	}
}

func TestMCPClient_ExecuteToolCalls_UnknownTool(t *testing.T) {
	client := NewMCPClient()
	// 不注册任何工具

	message := openai.ChatCompletionMessage{
		ToolCalls: []openai.ChatCompletionMessageToolCall{
			{
				ID:   "call_unknown",
				Type: "function",
				Function: openai.ChatCompletionMessageToolCallFunction{
					Name:      "unknown_tool",
					Arguments: `{}`,
				},
			},
		},
	}

	_, err := client.ExecuteToolCalls(message)
	if err == nil {
		t.Fatal("Expected error for unknown tool, got nil")
	}

	if err.Error() != "unknown tool: unknown_tool" {
		t.Errorf("Expected error 'unknown tool: unknown_tool', got '%v'", err)
	}
}

func TestGlobalFunctions(t *testing.T) {
	// 测试全局注册和执行函数
	RegisterTool("global_add", addFunc, AddArgs{})

	message := openai.ChatCompletionMessage{
		ToolCalls: []openai.ChatCompletionMessageToolCall{
			{
				ID:   "call_global",
				Type: "function",
				Function: openai.ChatCompletionMessageToolCallFunction{
					Name:      "global_add",
					Arguments: `{"a": 10, "b": 20}`,
				},
			},
		},
	}

	results, err := ExecuteToolCalls(message) // 使用全局函数
	if err != nil {
		t.Fatalf("Global ExecuteToolCalls failed: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 result from global function, got %d", len(results))
	}

	expectedResult := `{"result":30}`
	if results[0].OfTool.Content.OfString != openai.String(expectedResult) {
		t.Errorf("Expected global result '%s', got '%s'", expectedResult, results[0].OfTool.Content.OfString)
	}
}
