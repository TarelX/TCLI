package ai

import (
	"context"
	"errors"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// AnthropicClient 实现 Client 接口，对接 Anthropic Claude API
type AnthropicClient struct {
	client    *anthropic.Client
	modelName string
}

// NewAnthropicClient 创建 Anthropic 客户端
func NewAnthropicClient(apiKey, model, baseURL string) *AnthropicClient {
	opts := []option.RequestOption{
		option.WithAPIKey(apiKey),
	}
	if baseURL != "" {
		opts = append(opts, option.WithBaseURL(baseURL))
	}
	client := anthropic.NewClient(opts...)
	return &AnthropicClient{
		client:    &client,
		modelName: model,
	}
}

func (c *AnthropicClient) ModelName() string  { return c.modelName }
func (c *AnthropicClient) Provider() Provider { return ProviderAnthropic }

func (c *AnthropicClient) Complete(ctx context.Context, req Request) (string, error) {
	messages, system := splitMessages(req.Messages)

	params := anthropic.MessageNewParams{
		Model:     anthropic.Model(c.modelName),
		MaxTokens: int64(req.MaxTokens),
		Messages:  messages,
	}
	if system != "" {
		params.System = []anthropic.TextBlockParam{
			{Text: system},
		}
	}

	msg, err := c.client.Messages.New(ctx, params)
	if err != nil {
		return "", err
	}
	if len(msg.Content) == 0 {
		return "", errors.New("AI 返回了空响应")
	}
	// 取第一个 text block
	for _, block := range msg.Content {
		if block.Type == "text" {
			return block.AsText().Text, nil
		}
	}
	return "", errors.New("响应中没有文本内容")
}

func (c *AnthropicClient) Stream(ctx context.Context, req Request, handler StreamHandler) error {
	messages, system := splitMessages(req.Messages)

	params := anthropic.MessageNewParams{
		Model:     anthropic.Model(c.modelName),
		MaxTokens: int64(req.MaxTokens),
		Messages:  messages,
	}
	if system != "" {
		params.System = []anthropic.TextBlockParam{
			{Text: system},
		}
	}

	stream := c.client.Messages.NewStreaming(ctx, params)
	for stream.Next() {
		event := stream.Current()
		// 只处理文本增量事件
		if event.Type == "content_block_delta" {
			if event.Delta.Type == "text_delta" {
				handler(event.Delta.Text, false, nil)
			}
		}
	}
	if err := stream.Err(); err != nil {
		handler("", true, err)
		return err
	}
	handler("", true, nil)
	return nil
}

// splitMessages 将消息列表中的 system 消息分离出来
func splitMessages(msgs []Message) ([]anthropic.MessageParam, string) {
	var system string
	var params []anthropic.MessageParam
	for _, m := range msgs {
		switch m.Role {
		case RoleSystem:
			system = m.Content
		case RoleUser:
			params = append(params, anthropic.NewUserMessage(anthropic.NewTextBlock(m.Content)))
		case RoleAssistant:
			params = append(params, anthropic.NewAssistantMessage(anthropic.NewTextBlock(m.Content)))
		}
	}
	return params, system
}
