package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/flyhope/kubetea/comm"
	"github.com/nsf/termbox-go"
	"github.com/sirupsen/logrus"
)

const (
	DotChar = " • "
)

var (
	// MainStyle 框架样式
	MainStyle = lipgloss.NewStyle().MarginTop(1).MarginLeft(2)

	// BodyStyle 主体样式
	BodyStyle = lipgloss.NewStyle().MarginTop(0).MarginLeft(0)

	// SubtleStyle 灰色样式
	SubtleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))

	// DotStyle 圆点样式
	DotStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("236")).Render(DotChar)
)

// MainView 渲染主体页面
func MainView(s string) string {
	return MainStyle.Render(s + "\n\n")
}

// RunProgram 启动一个程序
func RunProgram(model tea.Model) (tea.Model, error) {
	comm.Program = tea.NewProgram(model, tea.WithAltScreen())
	return comm.Program.Run()
}

// ViewModel 开始展示一个Model
func ViewModel(model tea.Model) (tea.Model, tea.Cmd) {
	return model, tea.Sequence(tea.ExitAltScreen, tea.EnterAltScreen)
}

// WindowSize 获取窗口大小
func WindowSize() tea.WindowSizeMsg {
	err := termbox.Init()
	if err != nil {
		logrus.Fatal(err)
	}
	defer termbox.Close()

	width, height := termbox.Size()
	return tea.WindowSizeMsg{
		Width:  width,
		Height: height,
	}
}

// Abstract 含有上级来源的Model
type Abstract struct {
	LastModel tea.Model
}

func (m *Abstract) GoBack() (tea.Model, tea.Cmd) {
	model, cmd := ViewModel(m.LastModel)
	model.Update(comm.MsgUIBack(true))
	return model, cmd
}
