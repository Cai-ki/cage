package llm

import (
	"bytes"
	"context"
	"encoding/base64"
	"image"
	"image/png"

	"github.com/openai/openai-go"
)

func (c *LLMClient) vision(img image.Image) (string, error) {
	return c.visionWithPrompt(img, "Describe this image.")
}

func (c *LLMClient) visionWithPrompt(img image.Image, prompt string) (string, error) {
	// Encode image to PNG in memory
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return "", err
	}

	imageDataStr := "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes()) // Placeholder - needs base64 import

	resp, err := c.openai.Chat.Completions.New(
		context.Background(),
		openai.ChatCompletionNewParams{
			Model: c.cfg.VisionModel,
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.UserMessage([]openai.ChatCompletionContentPartUnionParam{
					{
						OfImageURL: &openai.ChatCompletionContentPartImageParam{
							ImageURL: openai.ChatCompletionContentPartImageImageURLParam{
								URL:    imageDataStr,
								Detail: "auto",
							},
						},
					},
					{
						OfText: &openai.ChatCompletionContentPartTextParam{
							Text: prompt,
						},
					},
				}),
			},
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
