package ai_test

import (
	"testing"

	"github.com/Cai-ki/cage/ai"
	_ "github.com/Cai-ki/cage/localconfig"
)

func TestEmbeddingText(t *testing.T) {
	a, err := ai.EmbeddingText("text", 4)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(a)
}
