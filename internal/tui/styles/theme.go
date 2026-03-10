package styles

import (
	"image/color"

	"charm.land/lipgloss/v2"
)

// Theme 定义整个 TUI 的色彩和样式
type Theme struct {
	// 消息气泡
	UserMessage lipgloss.Style
	AIMessage   lipgloss.Style

	// 布局容器
	TitleBar  lipgloss.Style
	StatusBar lipgloss.Style
	InputBox  lipgloss.Style

	// 状态指示
	TokenNormal  lipgloss.Style
	TokenWarning lipgloss.Style // 用量 > 80%
	TokenDanger  lipgloss.Style // 用量 > 95%

	// 文字辅助
	Dim     lipgloss.Style
	Bold    lipgloss.Style
	Success lipgloss.Style
	Error   lipgloss.Style

	// 启动动画渐变色（从 Logo 顶部到底部）
	GradientColors []color.Color
}

// Default 返回深色主题（自适应终端背景）
func Default() Theme {
	return Theme{
		UserMessage: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7DF9FF")).
			Bold(true),

		AIMessage: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E8E8E8")),

		TitleBar: lipgloss.NewStyle().
			Background(lipgloss.Color("#1A1A2E")).
			Foreground(lipgloss.Color("#A9A9C8")).
			Padding(0, 1),

		StatusBar: lipgloss.NewStyle().
			Background(lipgloss.Color("#16213E")).
			Foreground(lipgloss.Color("#777799")).
			Padding(0, 1),

		InputBox: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#4A4A8A")).
			Padding(0, 1),

		TokenNormal: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#55EFC4")),

		TokenWarning: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FDCB6E")),

		TokenDanger: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF7675")).
			Bold(true),

		Dim: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#555577")),

		Bold: lipgloss.NewStyle().
			Bold(true),

		Success: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00B894")),

		Error: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#D63031")).
			Bold(true),

		// 紫→蓝→青 渐变，用于启动动画 Logo
		GradientColors: []color.Color{
			lipgloss.Color("#BF5AF2"),
			lipgloss.Color("#9B59B6"),
			lipgloss.Color("#7D56F4"),
			lipgloss.Color("#5DADE2"),
			lipgloss.Color("#48C9B0"),
		},
	}
}
