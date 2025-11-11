package llm_test

import (
	"testing"

	"github.com/Cai-ki/cage/llm"
	"github.com/Cai-ki/cage/llm/mcp"
	"github.com/Cai-ki/cage/media"
	"github.com/Cai-ki/cage/sugar"
)

func TestCompletion(t *testing.T) {
	txt, err := llm.Completion("请你给我讲个冷笑话。")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(txt)
}

func TestCompletionBySystem(t *testing.T) {
	txt, err := llm.CompletionBySystem("请你给我讲个冷笑话。")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(txt)
}

func TestCompletionByParams(t *testing.T) {
	msg, err := llm.CompletionByParams(llm.UserMessage("调用function : add解决问题，1 + 1 = ？"),
		llm.ToolsByJson(`
[
  {
    "type": "function",
    "function": {
      "name": "add",
      "description": "将两个数字相加并返回结果。",
      "parameters": {
        "type": "object",
        "properties": {
          "a": {
            "type": "number",
            "description": "第一个加数"
          },
          "b": {
            "type": "number",
            "description": "第二个加数"
          }
        },
        "required": ["a", "b"]
      }
    }
  }
]		
`),
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(msg.Content, "\n", msg.ToolCalls[0].RawJSON())

	// 定义测试用的参数结构体
	type AddArgs struct {
		A float64 `json:"a"`
		B float64 `json:"b"`
	}

	// 定义测试用的函数
	addFunc := func(args AddArgs) (interface{}, error) {
		return map[string]interface{}{"result": args.A + args.B}, nil
	}

	mcp.RegisterTool("add", addFunc, AddArgs{})
	res, err := mcp.ExecuteToolCalls(msg)
	t.Log(res[0].OfTool.Content.OfString)
}

func TestVision(t *testing.T) {
	txt, err := llm.Vision(sugar.Must(media.Screenshot()))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(txt)
}

func TestVisionWithPrompt(t *testing.T) {
	txt, err := llm.VisionWithPrompt(sugar.Must(media.Screenshot()), "忽略图片文本内容，请你告诉我图片整体是什么颜色？")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(txt)
}

func TestEmbedding(t *testing.T) {
	arr, err := llm.Embedding("hello, world!")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(arr)
}

func TestEmbeddingWithDim(t *testing.T) {
	arr, err := llm.EmbeddingWithDim("hello, world!", 8)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(arr)
}
