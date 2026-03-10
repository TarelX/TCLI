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
	for _, msg := range m.messages {
		switch msg.Role {
		case ai.RoleUser:
			sb.WriteString(m.theme.UserMessage.Render("[你] "))
			sb.WriteString(msg.Content)
		case ai.RoleAssistant:
			sb.WriteString(m.theme.AIMessage.Render("[AI] "))
			// TODO Phase 2: 用 glamour 渲染 Markdown
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
