package tui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/TarelX/TCLI/internal/ai"
	"github.com/TarelX/TCLI/internal/tui/chat"
	"github.com/TarelX/TCLI/internal/tui/splash"
)

// RootModel 顶层 Model，负责管理 splash → 主界面的切换
type RootModel struct {
	splash     splash.Model
	chat       chat.Model
	showSplash bool
	width      int
	height     int
}

// NewRootModel 创建根 Model
func NewRootModel(client ai.Client, version string, showSplash bool) RootModel {
	modelName := ""
	tokenMax := 128000
	if client != nil {
		modelName = client.ModelName()
	}

	return RootModel{
		splash:     splash.New(80, version, modelName),
		chat:       chat.New(client, tokenMax),
		showSplash: showSplash,
	}
}

func (m RootModel) Init() tea.Cmd {
	if m.showSplash {
		return m.splash.Init()
	}
	return m.chat.Init()
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	if m.showSplash {
		newSplash, cmd := m.splash.Update(msg)
		m.splash = newSplash.(splash.Model)

		// 动画结束，切换到 chat 界面
		if _, ok := msg.(splash.DoneMsg); ok {
			m.showSplash = false
			return m, m.chat.Init()
		}
		return m, cmd
	}

	newChat, cmd := m.chat.Update(msg)
	m.chat = newChat.(chat.Model)
	return m, cmd
}

func (m RootModel) View() tea.View {
	if m.showSplash {
		return m.splash.View()
	}
	return m.chat.View()
}
