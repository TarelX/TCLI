package ai

import "context"

// Provider 标识 AI 服务提供商
type Provider string

const (
	ProviderOpenAI     Provider = "openai"
	ProviderAnthropic  Provider = "anthropic"
	ProviderCompatible Provider = "compatible"
)

// Role 消息角色
type Role string

const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleSystem    Role = "system"
)

// Message 单条对话消息
type Message struct {
	Role    Role
	Content string
}

// Request 统一请求结构
type Request struct {
	Messages    []Message
	MaxTokens   int
	Temperature float32
	Stream      bool
}

// StreamHandler 流式响应回调
// delta: 本次收到的文本片段
// done: 是否已结束
// err: 错误（非 nil 时 done 也为 true）
type StreamHandler func(delta string, done bool, err error)

// Client 统一 AI 客户端接口
type Client interface {
	// Complete 非流式调用，返回完整响应
	Complete(ctx context.Context, req Request) (string, error)
	// Stream 流式调用，通过 handler 回调逐步返回内容
	Stream(ctx context.Context, req Request, handler StreamHandler) error
	// ModelName 返回当前使用的模型名
	ModelName() string
	// Provider 返回 provider 类型
	Provider() Provider
}

// NewMessage 创建一条消息的快捷方法
func NewMessage(role Role, content string) Message {
	return Message{Role: role, Content: content}
}

func UserMessage(content string) Message      { return NewMessage(RoleUser, content) }
func AssistantMessage(content string) Message { return NewMessage(RoleAssistant, content) }
func SystemMessage(content string) Message    { return NewMessage(RoleSystem, content) }
