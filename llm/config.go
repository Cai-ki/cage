package llm

import (
	"os"
	"strconv"
)

type Config struct {
	APIKey      string
	BaseURL     string // 支持兼容 API（如 Ollama）
	Model       string // 默认文本模型
	VisionModel string // 默认视觉模型
	EmbedModel  string // 默认 embedding 模型
	EmbedDim    int
}

func LoadConfig() (*Config, error) {
	embedDim, err := strconv.Atoi(os.Getenv("LLM_EMBED_DIM"))
	if err != nil {
		embedDim = 0
	}

	return &Config{
		APIKey:      os.Getenv("LLM_API_KEY"),
		BaseURL:     os.Getenv("LLM_BASE_URL"),
		Model:       os.Getenv("LLM_MODEL"),
		VisionModel: os.Getenv("LLM_VISION_MODEL"),
		EmbedModel:  os.Getenv("LLM_EMBED_MODEL"),
		EmbedDim:    embedDim,
	}, nil
}
