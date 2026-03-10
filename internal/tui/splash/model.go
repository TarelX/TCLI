package splash

import (
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/TarelX/TCLI/internal/tui/styles"
)

// phase 动画阶段
type phase int

const (
	phaseTyping phase = iota // 逐字打印 "TCli"
	phaseLogo                // 展开完整 ASCII Art Logo
	phaseInfo                // 显示版本 + 加载点
	phaseDone                // 动画结束
)

// tickMsg 内部定时器消息
type tickMsg time.Time

// DoneMsg 动画结束后发送给父 Model，通知切换到主界面
type DoneMsg struct{}

const logoFull = `
████████╗ ██████╗██╗     ██╗
   ██╔══╝██╔════╝██║     ██║
   ██║   ██║     ██║     ██║
   ██║   ██║     ██║     ██║
   ██║   ╚██████╗███████╗██║
   ╚═╝    ╚═════╝╚══════╝╚═╝`

const logoCompact = `
 _____ _____ _  _
|_   _/  __ \ |(_)
  | | | /  \/ | _
  | | | |   | || |
 _| |_| \__/\ ||_|
 \___/ \____/_||_|`

const subtitle = "Terminal Code Intelligence CLI"
const appName = "TCli"

// Model 是启动动画的 Bubbletea Model
type Model struct {
	phase      phase
	typedChars int // 已打印的字符数（阶段一）
	dotCount   int // 加载点数量（阶段三）
	width      int
	version    string
	modelName  string
	theme      styles.Theme
}

func New(width int, version, modelName string) Model {
	return Model{
		width:     width,
		version:   version,
		modelName: modelName,
		theme:     styles.Default(),
	}
}

func (m Model) Init() tea.Cmd {
	return tick(80 * time.Millisecond)
}

func tick(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		return m, nil
	case tickMsg:
		switch m.phase {
		case phaseTyping:
			if m.typedChars < len([]rune(appName)) {
				m.typedChars++
				return m, tick(100 * time.Millisecond)
			}
			m.phase = phaseLogo
			return m, tick(300 * time.Millisecond)

		case phaseLogo:
			m.phase = phaseInfo
			return m, tick(500 * time.Millisecond)

		case phaseInfo:
			if m.dotCount < 3 {
				m.dotCount++
				return m, tick(180 * time.Millisecond)
			}
			m.phase = phaseDone
			return m, func() tea.Msg { return DoneMsg{} }
		}
	}
	return m, nil
}

func (m Model) View() tea.View {
	center := lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center)

	switch m.phase {
	case phaseTyping:
		runes := []rune(appName)
		typed := string(runes[:m.typedChars])
		styled := lipgloss.NewStyle().
			Foreground(m.theme.GradientColors[0]).
			Bold(true).
			Render(typed)
		v := tea.NewView("\n\n\n" + center.Render(styled) + "\n")
		v.AltScreen = true
		return v

	case phaseLogo, phaseInfo:
		logo := logoFull
		if m.width < 50 {
			logo = logoCompact
		}
		lines := strings.Split(strings.TrimPrefix(logo, "\n"), "\n")
		var rendered []string
		grad := m.theme.GradientColors
		for i, line := range lines {
			idx := (i * (len(grad) - 1)) / max(len(lines)-1, 1)
			colored := lipgloss.NewStyle().
				Foreground(grad[idx]).
				Bold(true).
				Render(line)
			rendered = append(rendered, center.Render(colored))
		}

		result := "\n" + strings.Join(rendered, "\n")
		result += "\n\n" + center.Render(m.theme.Dim.Render(subtitle))

		if m.phase == phaseInfo {
			dots := strings.Repeat("●", m.dotCount) + strings.Repeat("○", 3-m.dotCount)
			dotsStyled := lipgloss.NewStyle().
				Foreground(m.theme.GradientColors[2]).
				Render(dots)
			info := m.theme.Bold.Render(m.version) +
				m.theme.Dim.Render("  ·  ") +
				m.theme.Dim.Render(m.modelName) +
				m.theme.Dim.Render("  ·  初始化中 ") +
				dotsStyled
			result += "\n" + center.Render(info)
		}
		v := tea.NewView(result)
		v.AltScreen = true
		return v
	}
	v := tea.NewView("")
	v.AltScreen = true
	return v
}
