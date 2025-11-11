package llm

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/openai/openai-go"
)

func (c *LLMClient) completion(prompt string) (string, error) {
	msg, err := c.completionByParams(UserMessage(prompt))
	return msg.Content, err
}

func (c *LLMClient) completionBySystem(prompt string) (string, error) {
	msg, err := c.completionByParams(SystemMessage(prompt))
	return msg.Content, err
}

type MessageFunc func() openai.ChatCompletionMessageParamUnion

func UserMessage(prompt string) MessageFunc {
	return func() openai.ChatCompletionMessageParamUnion {
		return openai.UserMessage(prompt)
	}
}

func SystemMessage(prompt string) MessageFunc {
	return func() openai.ChatCompletionMessageParamUnion {
		return openai.SystemMessage(prompt)
	}
}

func ToolMessage(prompt string, toolCallID string) MessageFunc {
	return func() openai.ChatCompletionMessageParamUnion {
		return openai.ToolMessage(prompt, toolCallID)
	}
}

type ToolFunc func() []openai.ChatCompletionToolParam

type AllowedParam interface {
	Allowed()
}

func (MessageFunc) Allowed() {}
func (ToolFunc) Allowed()    {}

func ToolsByJson(jstr string) ToolFunc {
	return func() []openai.ChatCompletionToolParam {
		tools := []openai.ChatCompletionToolParam{}
		if err := json.Unmarshal([]byte(jstr), &tools); err != nil {
			return nil
		}
		return tools
	}
}

func (c *LLMClient) completionByParams(args ...AllowedParam) (openai.ChatCompletionMessage, error) {
	msgs := []openai.ChatCompletionMessageParamUnion{}
	tools := []openai.ChatCompletionToolParam{}
	for _, arg := range args {
		switch v := arg.(type) {
		case MessageFunc:
			msgs = append(msgs, v())
		case ToolFunc:
			tools = append(tools, v()...)
		default:
			return openai.ChatCompletionMessage{}, fmt.Errorf("unsupported argument type: %T", v)
		}
	}

	params := openai.ChatCompletionNewParams{
		Model:       c.cfg.Model,
		Messages:    msgs,
		Tools:       tools,
		Temperature: openai.Float(c.cfg.Temperature),
		TopP:        openai.Float(c.cfg.TopP),
	}

	resp, err := c.openai.Chat.Completions.New(
		context.Background(),
		params,
	)

	if err != nil {
		return openai.ChatCompletionMessage{}, err
	}
	if len(resp.Choices) == 0 {
		return openai.ChatCompletionMessage{}, ErrUnexpectedResponse
	}
	return resp.Choices[0].Message, nil
}
