package chat

import (
	"strings"
	"time"

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
	messages    []ai.Message
	input       textarea.Model
	viewport    viewport.Model
	spinner     spinner.Model
	streaming   bool
	streamStart time.Time // 流式请求开始时间，用于显示思考耗时
	streamBuf   *strings.Builder
	tokenUsed   int
	tokenMax    int
	width       int
	height      int
	client      ai.Client
	theme       styles.Theme
	err         error
	version     string // 版本号，用于欢迎面板
	projectType string // 项目类型，用于欢迎面板
	gitBranch   string // Git分支，用于欢迎面板
	workDir     string // 工作目录，用于欢迎面板

	// 自动补全状态
	completions    []string // 当前候选列表
	completionIdx  int      // 当前选中的候选索引
	completionMode string   // 补全模式："" | "file" | "cmd"
	lastInput      string   // 上一次输入内容，用于检测变化
}

// New 创建 chat Model
func New(client ai.Client, tokenMax int, version, projectType, gitBranch, workDir string) Model {
	ta := textarea.New()
	ta.Placeholder = "输入消息... (Enter 发送, /help 查看命令)"
	ta.Focus()
	ta.SetHeight(3)
	ta.CharLimit = 0
	ta.ShowLineNumbers = false
	ta.Prompt = "❯ "

	sp := spinner.New()
	sp.Spinner = spinner.Dot

	vp := viewport.New()

	return Model{
		input:       ta,
		viewport:    vp,
		spinner:     sp,
		streamBuf:   &strings.Builder{},
		tokenMax:    tokenMax,
		client:      client,
		theme:       styles.Default(),
		version:     version,
		projectType: projectType,
		gitBranch:   gitBranch,
		workDir:     workDir,
	}
}

func (m Model) Init() tea.Cmd {
	return textarea.Blink
}

// Ensure unused imports are kept
var _ = tea.NewView
