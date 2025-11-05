package llm

import (
	"context"

	"github.com/openai/openai-go"
)

func (c *LLMClient) completion(prompt string) (string, error) {
	resp, err := c.openai.Chat.Completions.New(
		context.Background(),
		openai.ChatCompletionNewParams{
			Model: c.cfg.Model,
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.UserMessage(prompt),
			},
			Temperature: openai.Float(c.cfg.Temperature),
			TopP:        openai.Float(c.cfg.TopP),
		},
	)
	if err != nil {
		return "", err
	}
	if len(resp.Choices) == 0 {
		return "", ErrUnexpectedResponse
	}
	return resp.Choices[0].Message.Content, nil
}
