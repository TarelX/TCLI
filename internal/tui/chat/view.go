package chat

import (
	"context"
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/TarelX/TCLI/internal/ai"
)

func (m Model) View() tea.View {
	if m.width == 0 {
		v := tea.NewView("正在初始化...")
		v.AltScreen = true
		return v
	}
	v := tea.NewView(strings.Join([]string{
		m.viewTitleBar(),
		m.viewport.View(),
		m.viewInputBox(),
		m.viewStatusBar(),
	}, "\n"))
	v.AltScreen = true
	return v
}

func (m Model) viewTitleBar() string {
	provider := ""
	if m.client != nil {
		provider = string(m.client.Provider())
		name := m.client.ModelName()
		provider = fmt.Sprintf(" TCli  %s  ● %s ", name, provider)
	} else {
		provider = " TCli "
	}
	return m.theme.TitleBar.Width(m.width).Render(provider)
}

func (m Model) viewInputBox() string {
	return m.theme.InputBox.Width(m.width - 2).Render(m.input.View())
}

func (m Model) viewStatusBar() string {
	tokenInfo := m.tokenStyle().Render(
		fmt.Sprintf("tokens: %s / %s",
			formatNum(m.tokenUsed),
			formatNum(m.tokenMax),
		),
	)
	help := m.theme.Dim.Render("↑↓ 翻历史  /help 命令  Ctrl+C 退出")
	gap := lipgloss.NewStyle().Width(m.width - lipgloss.Width(tokenInfo) - lipgloss.Width(help) - 2).Render("")
	return m.theme.StatusBar.Width(m.width).Render(tokenInfo + gap + help)
}

func (m Model) tokenStyle() lipgloss.Style {
	if m.tokenMax == 0 {
		return m.theme.TokenNormal
	}
	ratio := float64(m.tokenUsed) / float64(m.tokenMax)
	switch {
	case ratio > 0.95:
		return m.theme.TokenDanger
	case ratio > 0.80:
		return m.theme.TokenWarning
	default:
		return m.theme.TokenNormal
	}
}

func (m Model) renderMessages() string {
	var sb strings.Builder

	// 检查是否有用户可见的消息
	hasVisible := false
	for _, msg := range m.messages {
		if msg.Role != ai.RoleSystem {
			hasVisible = true
			break
		}
	}

	// 没有可见消息时显示 Claude Code 风格的欢迎面板
	if !hasVisible && !m.streaming {
		sb.WriteString(m.renderWelcomePanel())
		return sb.String()
	}

	for _, msg := range m.messages {
		switch msg.Role {
		case ai.RoleSystem:
			continue
		case ai.RoleUser:
			sb.WriteString(m.theme.UserMessage.Render("[你] "))
			sb.WriteString(msg.Content)
		case ai.RoleAssistant:
			sb.WriteString(m.theme.AIMessage.Render("[AI] "))
			sb.WriteString(msg.Content)
		}
		sb.WriteString("\n\n")
	}
	// 流式响应中的临时内容
	if m.streaming && m.streamBuf.Len() > 0 {
		sb.WriteString(m.theme.AIMessage.Render("[AI] "))
		sb.WriteString(m.streamBuf.String())
		sb.WriteString(m.spinner.View())
	}
	return sb.String()
}

// logoLines 是 TCLi 的方块风格 ASCII Art，分行存储用于逐行渐变着色
var logoLines = []string{
	` ████████╗  ██████╗ ██╗      ██╗`,
	` ╚══██╔══╝ ██╔════╝ ██║      ██║`,
	`    ██║    ██║      ██║      ██║`,
	`    ██║    ██║      ██║      ██║`,
	`    ██║    ╚██████╗ ███████╗ ██║`,
	`    ╚═╝     ╚═════╝ ╚══════╝ ╚═╝`,
}

// renderWelcomePanel 渲染带渐变Logo的持久欢迎面板
func (m Model) renderWelcomePanel() string {
	w := m.width - 4
	if w < 40 {
		w = 40
	}

	// 渲染渐变 Logo（紫→蓝→青）
	grad := m.theme.GradientColors
	var logoRendered strings.Builder
	for i, line := range logoLines {
		idx := (i * (len(grad) - 1)) / max(len(logoLines)-1, 1)
		colored := lipgloss.NewStyle().
			Foreground(grad[idx]).
			Bold(true).
			Render(line)
		logoRendered.WriteString(colored + "\n")
	}

	// 模型名称
	modelName := "未配置"
	providerName := ""
	if m.client != nil {
		modelName = m.client.ModelName()
		providerName = string(m.client.Provider())
	}

	// 构建左侧信息栏
	var left strings.Builder
	left.WriteString(m.theme.Bold.Render("Welcome!"))
	left.WriteString("\n\n")
	left.WriteString(fmt.Sprintf("  🤖 模型: %s", m.theme.Bold.Render(modelName)))
	if providerName != "" {
		left.WriteString(m.theme.Dim.Render(" · " + providerName))
	}
	left.WriteString("\n")
	if m.projectType != "" {
		left.WriteString(fmt.Sprintf("  📁 项目: %s", m.theme.Bold.Render(m.projectType)))
		left.WriteString("\n")
	}
	if m.gitBranch != "" {
		left.WriteString(fmt.Sprintf("  🌿 分支: %s", m.theme.Bold.Render(m.gitBranch)))
		left.WriteString("\n")
	}
	if m.workDir != "" {
		left.WriteString(fmt.Sprintf("  📂 目录: %s", m.theme.Dim.Render(m.workDir)))
		left.WriteString("\n")
	}

	// 构建右侧提示栏
	var right strings.Builder
	right.WriteString(m.theme.Success.Render("Tips"))
	right.WriteString("\n")
	right.WriteString("  直接输入问题开始对话\n")
	right.WriteString("  支持多轮连续对话\n")
	right.WriteString("  AI 会记住上下文\n")
	right.WriteString("\n")
	right.WriteString(m.theme.Success.Render("可用命令"))
	right.WriteString("\n")
	right.WriteString("  /clear  清空对话历史\n")
	right.WriteString("  /copy   复制最后回复\n")
	right.WriteString("  /help   查看帮助\n")

	// 用边框包裹
	leftContent := left.String()
	rightContent := right.String()

	halfW := w/2 - 2
	if halfW < 20 {
		halfW = 20
	}

	leftBox := lipgloss.NewStyle().Width(halfW).Render(leftContent)
	rightBox := lipgloss.NewStyle().
		Width(halfW).
		BorderLeft(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#4A4A8A")).
		PaddingLeft(2).
		Render(rightContent)

	body := lipgloss.JoinHorizontal(lipgloss.Top, leftBox, rightBox)

	// 居中Logo
	logoCenter := lipgloss.NewStyle().Width(w).Align(lipgloss.Center).Render(logoRendered.String())

	// 版本号（Logo下方居中）
	versionLine := lipgloss.NewStyle().Width(w).Align(lipgloss.Center).Render(
		m.theme.Dim.Render(fmt.Sprintf("TCli %s", m.version)),
	)

	// 外层边框
	panel := lipgloss.NewStyle().
		Width(w).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#BF5AF2")).
		Padding(1, 2).
		Render(
			logoCenter + "\n" + versionLine + "\n\n" + body,
		)

	return "\n" + panel + "\n"
}

// startStream 启动流式 AI 请求，通过 tea.Cmd 将增量注入 bubbletea 消息循环
func (m Model) startStream() tea.Cmd {
	client := m.client
	messages := make([]ai.Message, len(m.messages))
	copy(messages, m.messages)

	return func() tea.Msg {
		// 用 channel 桥接流式回调和 bubbletea 消息循环
		ch := make(chan tea.Msg, 64)

		go func() {
			defer close(ch)
			req := ai.Request{
				Messages:    messages,
				MaxTokens:   8192,
				Temperature: 0.7,
				Stream:      true,
			}
			err := client.Stream(context.Background(), req, func(delta string, done bool, streamErr error) {
				if streamErr != nil {
					ch <- StreamErrMsg{Err: streamErr}
					return
				}
				if done {
					ch <- StreamDoneMsg{}
					return
				}
				if delta != "" {
					ch <- StreamDeltaMsg{Text: delta}
				}
			})
			if err != nil {
				ch <- StreamErrMsg{Err: err}
			}
		}()

		// 返回第一条消息，后续通过 waitForStream 持续读取
		msg, ok := <-ch
		if !ok {
			return StreamDoneMsg{}
		}
		// 把 channel 存下来给后续 waitForStream 使用
		streamCh = ch
		return msg
	}
}

// streamCh 用于在多次 tea.Cmd 调用之间传递 channel
// 注意：这是包级变量，单实例安全
var streamCh chan tea.Msg

// waitForStream 持续从 channel 读取流式消息
func waitForStream() tea.Msg {
	if streamCh == nil {
		return StreamDoneMsg{}
	}
	msg, ok := <-streamCh
	if !ok {
		streamCh = nil
		return StreamDoneMsg{}
	}
	return msg
}

func formatNum(n int) string {
	if n >= 1000 {
		return fmt.Sprintf("%d,%03d", n/1000, n%1000)
	}
	return fmt.Sprintf("%d", n)
}
