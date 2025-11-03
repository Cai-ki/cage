package llm_test

import (
	"testing"

	"github.com/Cai-ki/cage/llm"
	"github.com/Cai-ki/cage/media"
	"github.com/Cai-ki/cage/sugar"
)

func TestCompletion(t *testing.T) {
	txt, err := llm.Completion("1 + 1 = ?")
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
	txt, err := llm.VisionWithPrompt(sugar.Must(media.Screenshot()), "图片整体基调是什么颜色？")
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
