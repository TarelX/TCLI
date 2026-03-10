package chat

import (
	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"github.com/TarelX/TCLI/internal/ai"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.input.SetWidth(msg.Width - 2)
		m.viewport.SetWidth(msg.Width)
		m.viewport.SetHeight(msg.Height - 7) // 减去标题栏、输入框、状态栏高度
		m.refreshViewport()

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			if !m.streaming {
				return m.handleSubmit()
			}
		}

	case StreamDeltaMsg:
		m.streamBuf.WriteString(msg.Text)
		m.refreshViewport()
		return m, nil

	case StreamDoneMsg:
		m.streaming = false
		// 将完整响应写入消息历史
		if m.streamBuf.Len() > 0 {
			// TODO: 计算 token 并更新 m.tokenUsed
		}
		m.streamBuf.Reset()

	case StreamErrMsg:
		m.streaming = false
		m.err = msg.Err
		m.streamBuf.Reset()
	}

	// 更新子组件
	var inputCmd, vpCmd, spCmd tea.Cmd

	m.input, inputCmd = m.input.Update(msg)
	cmds = append(cmds, inputCmd)

	m.viewport, vpCmd = m.viewport.Update(msg)
	cmds = append(cmds, vpCmd)

	if m.streaming {
		m.spinner, spCmd = m.spinner.Update(msg)
		cmds = append(cmds, spCmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) handleSubmit() (Model, tea.Cmd) {
	text := m.input.Value()
	if text == "" {
		return m, nil
	}

	// 处理 / 命令
	if len(text) > 0 && text[0] == '/' {
		return m.handleSlashCommand(text)
	}

	// 普通消息
	m.messages = append(m.messages, ai.UserMessage(text))
	m.input.Reset()
	m.streaming = true
	m.refreshViewport()

	return m, m.startStream()
}

func (m Model) handleSlashCommand(cmd string) (Model, tea.Cmd) {
	switch cmd {
	case "/clear":
		m.messages = nil
		m.tokenUsed = 0
		m.input.Reset()
		m.refreshViewport()
	case "/copy":
		// TODO: 复制最后一条回复到剪贴板
	case "/help":
		// TODO: 显示帮助信息
	}
	return m, nil
}

func (m *Model) refreshViewport() {
	m.viewport.SetContent(m.renderMessages())
	m.viewport.GotoBottom()
}

// 以下为占位符，确保 import 不被移除
var _ = textarea.New
var _ = viewport.New
var _ ai.Message
