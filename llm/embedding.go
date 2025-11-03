package llm

import (
	"context"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/packages/param"
)

func (c *LLMClient) embedding(text string) ([]float32, error) {
	return c.embeddingWithDim(text, c.cfg.EmbedDim)
}

func (c *LLMClient) embeddingWithDim(text string, dimensions int) ([]float32, error) {
	resp, err := c.openai.Embeddings.New(
		context.Background(),
		openai.EmbeddingNewParams{
			Input: openai.EmbeddingNewParamsInputUnion{
				OfString: param.Opt[string]{Value: text},
			},
			Model:      c.cfg.EmbedModel,
			Dimensions: param.NewOpt(int64(dimensions)),
		},
	)
	if err != nil {
		return nil, err
	}
	if len(resp.Data) == 0 {
		return nil, ErrUnexpectedResponse
	}
	// Note: The openai-go SDK returns []float64, convert to []float32 if needed
	embedding := resp.Data[0].Embedding
	float32Embedding := make([]float32, len(embedding))
	for i, v := range embedding {
		float32Embedding[i] = float32(v)
	}
	return float32Embedding, nil
}
