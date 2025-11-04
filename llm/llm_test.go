package llm_test

import (
	"testing"

	"github.com/Cai-ki/cage/llm"
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
