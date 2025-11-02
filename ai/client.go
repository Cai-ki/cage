package ai

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/packages/param"
)

var EmbeddingClient openai.Client
var EmbeddingModel = ""

func init() {
	model := os.Getenv("EMBEDDING_MODEL")
	if model == "" {
		log.Fatal("EMBEDDING_MODEL not found")
	}

	EmbeddingModel = model

	url := os.Getenv("EMBEDDING_URL")
	if url == "" {
		log.Fatal("EMBEDDING_URL not found")
	}

	key := os.Getenv("EMBEDDING_API_KEY")
	if key == "" {
		log.Println("EMBEDDING_API_KEY not found")
	}

	EmbeddingClient = openai.NewClient(
		option.WithBaseURL(url),
		option.WithAPIKey(key),
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
