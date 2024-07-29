package comm

import (
	tea "github.com/charmbracelet/bubbletea"
	"time"
)

var Program *tea.Program

// MsgPodCache 缓存数据更新消息
type MsgPodCache time.Time

// MsgUIBack 返回消息
type MsgUIBack bool
