package token

import "github.com/TarelX/TCLI/internal/ai"

// EstimateMessages 估算消息列表的 token 数（快速估算版本）
// 精确计算需要 tiktoken，此处用字符数 / 4 作为近似值（英文约 4 字符/token，中文约 1.5 字符/token）
func EstimateMessages(msgs []ai.Message) int {
	total := 0
	for _, m := range msgs {
		total += EstimateText(m.Content)
		total += 4 // 每条消息的结构开销
	}
	return total
}

// EstimateText 估算单段文本的 token 数
func EstimateText(text string) int {
	if text == "" {
		return 0
	}
	runes := []rune(text)
	// 简单区分中文（CJK）和其他字符
	cjk := 0
	for _, r := range runes {
		if r >= 0x4E00 && r <= 0x9FFF {
			cjk++
		}
	}
	others := len(runes) - cjk
	// 中文约 1 rune ≈ 1.2 token，其他约 1 rune ≈ 0.25 token
	return int(float64(cjk)*1.2 + float64(others)*0.25)
}
