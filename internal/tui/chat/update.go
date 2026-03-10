package chat

import (
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/TarelX/TCLI/internal/ai"
	cctx "github.com/TarelX/TCLI/internal/context"
	"github.com/TarelX/TCLI/internal/prompt"
	"github.com/TarelX/TCLI/internal/token"
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
		// 继续读取下一条流式消息
		return m, waitForStream

	case StreamDoneMsg:
		m.streaming = false
		// 将完整响应写入消息历史
		if m.streamBuf.Len() > 0 {
			resp := m.streamBuf.String()
			m.messages = append(m.messages, ai.AssistantMessage(resp))
			m.tokenUsed = token.EstimateMessages(m.messages)
		}
		m.streamBuf.Reset()
		m.refreshViewport()

	case StreamErrMsg:
		m.streaming = false
		m.err = msg.Err
		// 显示错误信息在消息区域
		if msg.Err != nil {
			m.messages = append(m.messages, ai.AssistantMessage(
				fmt.Sprintf("⚠️ 请求出错：%v", msg.Err),
			))
		}
		m.streamBuf.Reset()
		m.refreshViewport()
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
	text := strings.TrimSpace(m.input.Value())
	if text == "" {
		return m, nil
	}

	// 处理 / 命令
	if len(text) > 0 && text[0] == '/' {
		return m.handleSlashCommand(text)
	}

	// 注入 system prompt（首次发送时）
	if len(m.messages) == 0 {
		m.messages = append(m.messages, ai.SystemMessage(prompt.SystemBase))
	}

	// 解析 @file 引用，将文件内容注入消息
	text = m.expandFileRefs(text)

	// 普通消息
	m.messages = append(m.messages, ai.UserMessage(text))
	m.input.Reset()
	m.streaming = true
	m.tokenUsed = token.EstimateMessages(m.messages)
	m.refreshViewport()

	return m, tea.Batch(m.startStream(), m.spinner.Tick)
}

func (m Model) handleSlashCommand(input string) (Model, tea.Cmd) {
	m.input.Reset()

	switch strings.TrimSpace(input) {
	case "/clear":
		m.messages = nil
		m.tokenUsed = 0
		m.refreshViewport()

	case "/copy":
		// 复制最后一条 AI 回复到剪贴板
		for i := len(m.messages) - 1; i >= 0; i-- {
			if m.messages[i].Role == ai.RoleAssistant {
				if err := clipboardCopy(m.messages[i].Content); err != nil {
					m.messages = append(m.messages, ai.AssistantMessage("⚠️ 复制失败："+err.Error()))
				} else {
					m.messages = append(m.messages, ai.AssistantMessage("✓ 已复制到剪贴板"))
				}
				break
			}
		}
		m.refreshViewport()

	case "/help":
		helpText := `可用命令：
  /clear   清空对话历史
  /copy    复制最后一条 AI 回复到剪贴板
  /help    显示此帮助信息

快捷键：
  Enter       发送消息
  Ctrl+C      退出`
		m.messages = append(m.messages, ai.AssistantMessage(helpText))
		m.refreshViewport()

	default:
		m.messages = append(m.messages, ai.AssistantMessage("未知命令："+input+"，输入 /help 查看可用命令"))
		m.refreshViewport()
	}

	return m, nil
}

func (m *Model) refreshViewport() {
	m.viewport.SetContent(m.renderMessages())
	m.viewport.GotoBottom()
}

// expandFileRefs 解析消息中的 @文件路径 引用，将文件内容注入消息
// 支持格式：@main.go  @./src/utils.go  @internal/ai/client.go
func (m Model) expandFileRefs(text string) string {
	// 用正则匹配 @filepath 模式
	re := regexp.MustCompile(`@([\w./\\-]+\.\w+)`)
	matches := re.FindAllStringSubmatch(text, -1)
	if len(matches) == 0 {
		return text
	}

	var fileContents strings.Builder
	loadedFiles := 0
	for _, match := range matches {
		filePath := match[1]
		content, err := cctx.ReadFile(filePath)
		if err != nil {
			continue
		}
		fileContents.WriteString(fmt.Sprintf("\n\n## 文件：%s\n\n%s", filePath, content))
		loadedFiles++
	}

	if loadedFiles > 0 {
		// 移除原文中的 @filepath，追加文件内容
		cleaned := re.ReplaceAllString(text, "`$1`")
		return cleaned + fileContents.String()
	}
	return text
}

// clipboardCopy 跨平台复制到剪贴板
func clipboardCopy(text string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("clip")
	case "darwin":
		cmd = exec.Command("pbcopy")
	default:
		cmd = exec.Command("xclip", "-selection", "clipboard")
	}
	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}
