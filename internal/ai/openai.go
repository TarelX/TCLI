package ai

import (
	"context"
	"errors"
	"io"

	openai "github.com/sashabaranov/go-openai"
)

// OpenAIClient 实现 Client 接口，对接 OpenAI API
type OpenAIClient struct {
	client    *openai.Client
	modelName string
}

// NewOpenAIClient 创建 OpenAI 客户端
// baseURL 为空时使用官方默认地址
func NewOpenAIClient(apiKey, model, baseURL string) *OpenAIClient {
	cfg := openai.DefaultConfig(apiKey)
	if baseURL != "" {
		cfg.BaseURL = baseURL
	}
	return &OpenAIClient{
		client:    openai.NewClientWithConfig(cfg),
		modelName: model,
	}
}

func (c *OpenAIClient) ModelName() string  { return c.modelName }
func (c *OpenAIClient) Provider() Provider { return ProviderOpenAI }

func (c *OpenAIClient) Complete(ctx context.Context, req Request) (string, error) {
	resp, err := c.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       c.modelName,
		Messages:    toOpenAIMessages(req.Messages),
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
	})
	if err != nil {
		return "", err
	}
	if len(resp.Choices) == 0 {
		return "", errors.New("AI 返回了空响应")
	}
	return resp.Choices[0].Message.Content, nil
}

func (c *OpenAIClient) Stream(ctx context.Context, req Request, handler StreamHandler) error {
	stream, err := c.client.CreateChatCompletionStream(ctx, openai.ChatCompletionRequest{
		Model:       c.modelName,
		Messages:    toOpenAIMessages(req.Messages),
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
		Stream:      true,
	})
	if err != nil {
		return err
	}
	defer stream.Close()

	for {
		resp, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			handler("", true, nil)
			return nil
		}
		if err != nil {
			handler("", true, err)
			return err
		}
		if len(resp.Choices) > 0 {
			delta := resp.Choices[0].Delta.Content
			if delta != "" {
				handler(delta, false, nil)
			}
		}
	}
}

func toOpenAIMessages(msgs []Message) []openai.ChatCompletionMessage {
	result := make([]openai.ChatCompletionMessage, len(msgs))
	for i, m := range msgs {
		result[i] = openai.ChatCompletionMessage{
			Role:    string(m.Role),
			Content: m.Content,
		}
	}
	return result
}
