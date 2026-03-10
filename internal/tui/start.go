package tui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/TarelX/TCLI/internal/ai"
)

// StartChat 启动交互式 Chat TUI
// showSplash: 是否显示启动动画
func StartChat(client ai.Client, version string, showSplash bool) error {
	model := NewRootModel(client, version, showSplash)
	p := tea.NewProgram(model)
	_, err := p.Run()
	return err
}
