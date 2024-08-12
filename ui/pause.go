package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// 暂停组件，用于执行命令自动退出时，暂停，看屏幕上的内容
type pause struct {
	lastModel tea.Model
	tips      string
}

func NewPause(lastModel tea.Model, tips string) *pause {
	return &pause{
		lastModel: lastModel,
		tips:      tips,
	}
}

func (m *pause) Init() tea.Cmd {
	return nil
}

func (m *pause) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msgType := msg.(type) {
	case tea.KeyMsg:
		switch msgType.String() {
		case "ctrl+c", "q", "enter":
			return m.lastModel, tea.EnterAltScreen
		}
	}
	return m, nil
}

func (m *pause) View() string {
	size := WindowSize()
	tipsStyle := SubtleStyle.Width(size.Width)

	result := "\n" + tipsStyle.Render(m.tips)
	result += "\npress [enter] to reutrn kubetea"
	return result
}
