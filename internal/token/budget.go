package token

import (
	"github.com/TarelX/TCLI/internal/ai"
)

// Budget 管理对话的 token 预算
type Budget struct {
	Limit      int // 模型的上下文窗口大小
	SystemCost int // 系统 prompt 占用的 token 数
	Reserved   int // 为响应预留的 token 数（默认 2048）
}

// NewBudget 创建预算管理器
func NewBudget(limit, reserved int) *Budget {
	return &Budget{
		Limit:    limit,
		Reserved: reserved,
	}
}

// Available 返回对话历史可用的 token 数
func (b *Budget) Available() int {
	return b.Limit - b.SystemCost - b.Reserved
}

// TrimHistory 在超出预算时截断最早的对话历史
// 保留最后 keepLast 条不被截断（确保最新上下文完整）
// 返回截断后的消息列表和被截断的条数
func (b *Budget) TrimHistory(messages []ai.Message, keepLast int) ([]ai.Message, int) {
	available := b.Available()
	if available <= 0 {
		return messages, 0
	}

	total := EstimateMessages(messages)
	if total <= available {
		return messages, 0
	}

	// 从最早的消息开始删除，但保留最后 keepLast 条
	trimmed := 0
	result := make([]ai.Message, len(messages))
	copy(result, messages)

	// 最多可删除的消息数量（保护最后 keepLast 条）
	maxRemovable := len(result) - keepLast
	if maxRemovable < 0 {
		maxRemovable = 0
	}

	for trimmed < maxRemovable && EstimateMessages(result) > available {
		result = result[1:]
		trimmed++
	}
	return result, trimmed
}

// ModelLimits 常见模型的上下文窗口大小
var ModelLimits = map[string]int{
	"gpt-4o":            128000,
	"gpt-4o-mini":       128000,
	"claude-3-7-sonnet": 200000,
	"claude-3-5-haiku":  200000,
	"deepseek-coder-v3": 128000,
	"qwen2.5-coder:7b":  32768,
}

// LimitForModel 返回模型的 token 上限，未知模型返回默认值 32768
func LimitForModel(model string) int {
	if limit, ok := ModelLimits[model]; ok {
		return limit
	}
	return 32768
}
