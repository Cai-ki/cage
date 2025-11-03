package llm

import (
	"image"

	_ "github.com/Cai-ki/cage/localconfig"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

var defaultClient *LLMClient

type LLMClient struct {
	cfg    *Config
	openai *openai.Client
}

func init() {
	cfg, err := LoadConfig()
	if err != nil {
		// Delay error to first call
		return
	}
	defaultClient, _ = newClient(cfg) // Error handling is done later
}

func newClient(cfg *Config) (*LLMClient, error) {
	// Create client using the new openai-go pattern
	clientOptions := []option.RequestOption{
		option.WithAPIKey(cfg.APIKey),
	}

	if cfg.BaseURL != "" {
		clientOptions = append(clientOptions, option.WithBaseURL(cfg.BaseURL))
	}

	client := openai.NewClient(clientOptions...)

	return &LLMClient{
		cfg:    cfg,
		openai: &client,
	}, nil
}

func initDefaultClient() error {
	if defaultClient == nil {
		cfg, err := LoadConfig()
		if err != nil {
			return err
		}
		var newErr error
		defaultClient, newErr = newClient(cfg)
		if newErr != nil {
			return newErr
		}
	}

	return nil
}

// Completion generates text from a text prompt.
func Completion(prompt string) (string, error) {
	err := initDefaultClient()
	if err != nil {
		return "", err
	}

	return defaultClient.completion(prompt)
}

// Vision analyzes an image and returns a textual description.
// Note: This function's implementation might require updates based on how openai-go handles images.
func Vision(img image.Image) (string, error) {
	err := initDefaultClient()
	if err != nil {
		return "", err
	}

	return defaultClient.vision(img)
}

// VisionWithPrompt analyzes an image with a custom instruction.
func VisionWithPrompt(img image.Image, prompt string) (string, error) {
	err := initDefaultClient()
	if err != nil {
		return "", err
	}

	return defaultClient.visionWithPrompt(img, prompt)
}

// Embedding returns a vector representation of the input text.
func Embedding(text string) ([]float32, error) {
	err := initDefaultClient()
	if err != nil {
		return nil, err
	}

	return defaultClient.embedding(text)
}

// EmbeddingWithDim returns embedding vector with specified dimension.
func EmbeddingWithDim(text string, dimensions int) ([]float32, error) {
	err := initDefaultClient()
	if err != nil {
		return nil, err
	}

	return defaultClient.embeddingWithDim(text, dimensions)
}

// func Transcribe(audio io.Reader) (string, error)
// func TranscribeWithPrompt(audio io.Reader, prompt string) (string, error)
