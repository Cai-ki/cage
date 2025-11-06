package llm

import (
	"os"

	"github.com/Cai-ki/cage/sugar"
)

type Config struct {
	APIKey      string
	BaseURL     string // 支持兼容 API（如 Ollama）
	Model       string // 默认文本模型
	VisionModel string // 默认视觉模型
	EmbedModel  string // 默认 embedding 模型
	EmbedDim    int
	Temperature float64
	TopP        float64
}

func LoadConfig() (*Config, error) {
	return &Config{
		APIKey:      os.Getenv("LLM_API_KEY"),
		BaseURL:     os.Getenv("LLM_BASE_URL"),
		Model:       os.Getenv("LLM_MODEL"),
		VisionModel: os.Getenv("LLM_VISION_MODEL"),
		EmbedModel:  os.Getenv("LLM_EMBED_MODEL"),
		EmbedDim:    sugar.StrToTWithDefault(os.Getenv("LLM_EMBED_DIM"), 0),
		Temperature: sugar.StrToTWithDefault(os.Getenv("LLM_TEMPERATURE"), 0.0),
		TopP:        sugar.StrToTWithDefault(os.Getenv("LLM_TOPP"), 0.8),
	}, nil
}
