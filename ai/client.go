package ai

import (
	"context"
	"fmt"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/packages/param"
)

var EmbeddingClient openai.Client
var EmbeddingModel = "qwen3-embedding:0.6b"

func init() {
	EmbeddingClient = openai.NewClient(
		option.WithBaseURL("http://localhost:11434/v1"),
	)
}

func EmbeddingText(text string, dim int) ([]float64, error) {
	params := openai.EmbeddingNewParams{
		Input: openai.EmbeddingNewParamsInputUnion{
			OfString: param.Opt[string]{Value: text},
		},
		Model:      EmbeddingModel,
		Dimensions: param.NewOpt(int64(dim)),
	}

	resp, err := EmbeddingClient.Embeddings.New(context.TODO(), params)
	if err != nil {
		return nil, err
	}

	if len(resp.Data) > 0 {
		vec := resp.Data[0].Embedding
		return vec, nil
	}

	return nil, fmt.Errorf("No response")
}
