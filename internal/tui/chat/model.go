package chat

import (
	"strings"

	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"github.com/TarelX/TCLI/internal/ai"
	"github.com/TarelX/TCLI/internal/tui/styles"
)

// StreamDeltaMsg 流式响应增量
type StreamDeltaMsg struct{ Text string }

// StreamDoneMsg 流式响应结束
type StreamDoneMsg struct{}

// StreamErrMsg 流式响应错误
type StreamErrMsg struct{ Err error }

// Model 是 chat 界面的 Bubbletea Model
type Model struct {
	messages  []ai.Message
	input     textarea.Model
	viewport  viewport.Model
	spinner   spinner.Model
	streaming bool
	streamBuf strings.Builder
	tokenUsed int
	tokenMax  int
	width     int
	height    int
	client    ai.Client
	theme     styles.Theme
	err       error
}

// New 创建 chat Model
func New(client ai.Client, tokenMax int) Model {
	ta := textarea.New()
	ta.Placeholder = "输入消息... (Enter 发送，Shift+Enter 换行，/help 查看命令)"
	ta.Focus()
	ta.SetHeight(3)
	ta.CharLimit = 0

	sp := spinner.New()
	sp.Spinner = spinner.Dot

	vp := viewport.New()

	return Model{
		input:    ta,
		viewport: vp,
		spinner:  sp,
		tokenMax: tokenMax,
		client:   client,
		theme:    styles.Default(),
	}
}

func (m Model) Init() tea.Cmd {
	return textarea.Blink
}

// Ensure unused imports are kept
var _ = tea.NewView
