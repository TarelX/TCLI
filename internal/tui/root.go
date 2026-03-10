package tui

import (
	"os"
	"path/filepath"

	tea "charm.land/bubbletea/v2"
	"github.com/TarelX/TCLI/internal/ai"
	cctx "github.com/TarelX/TCLI/internal/context"
	"github.com/TarelX/TCLI/internal/token"
	"github.com/TarelX/TCLI/internal/tui/chat"
)

// RootModel 顶层 Model，直接使用 chat（欢迎面板内置在 chat 中）
type RootModel struct {
	chat   chat.Model
	width  int
	height int
}

// NewRootModel 创建根 Model，收集项目信息传给 chat 用于欢迎面板
func NewRootModel(client ai.Client, version string) RootModel {
	tokenMax := 128000
	if client != nil {
		tokenMax = token.LimitForModel(client.ModelName())
	}

	// 收集项目信息
	gitInfo := cctx.CollectGitInfo()
	projectType := cctx.Collect(0).ProjectType

	workDir := ""
	if wd, err := os.Getwd(); err == nil {
		workDir = filepath.Base(wd)
	}

	return RootModel{
		chat: chat.New(client, tokenMax, version, projectType, gitInfo.Branch, workDir),
	}
}

func (m RootModel) Init() tea.Cmd {
	return m.chat.Init()
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	newChat, cmd := m.chat.Update(msg)
	m.chat = newChat.(chat.Model)
	return m, cmd
}

func (m RootModel) View() tea.View {
	return m.chat.View()
}
