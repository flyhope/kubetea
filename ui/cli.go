package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"os/exec"
)

type CliMsg struct {
	Err error
}

// NewCli 执行一个Cli命令
func NewCli(command string, args ...string) tea.Cmd {
	return NewCliWithCallback(func(err error) tea.Msg {
		return CliMsg{Err: err}
	}, command, args...)
}

// NewCliWithCallback 执行一个Cli命令，并手动设置Callback函数
func NewCliWithCallback(fn func(err error) tea.Msg, command string, args ...string) tea.Cmd {
	cmd := exec.Command(command, args...)
	return NewCmdWithCallback(cmd, fn)
}

// NewCmd 执行一个*exec.Cmd命令
func NewCmd(cmd *exec.Cmd) tea.Cmd {
	return NewCmdWithCallback(cmd, func(err error) tea.Msg {
		return CliMsg{Err: err}
	})
}

// NewCmdWithCallback 执行一个*exec.Cmd命令，并手动设置Callback函数
func NewCmdWithCallback(cmd *exec.Cmd, fn func(err error) tea.Msg) tea.Cmd {
	teaCmd := tea.ExecProcess(cmd, fn)
	return teaCmd
}
