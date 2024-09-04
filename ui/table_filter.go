package ui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sirupsen/logrus"
	"strconv"
)

// TableFilter 带有过滤功能的Table
type TableFilter struct {
	Table
	Input       textinput.Model
	UpdateEvent func(msg tea.Msg) (tea.Model, tea.Cmd)
	SubDescs    []string
}

// 上一次的表格组件（必需用上一次的克隆一个，否则新表格设置宽高屏幕会闪）
var lastTableFilter = newTableFilter()

// Reset 重置表格数据
func (m *TableFilter) Reset() {
	m.Input.SetValue("")
	m.Table.filterValue = ""

	if len(m.Table.Rows()) > 0 {
		m.Table.SetCursor(0)
	}
	m.Table.SetRows(nil)
	m.Table.SetColumns(nil)
	m.SubDescs = nil
	m.UpdateEvent = nil
}

// FetchTableFilter 创建一个新的可筛选表格
func FetchTableFilter() *TableFilter {
	tableCloneEntity := *lastTableFilter
	tableClone := &tableCloneEntity
	tableClone.Reset()
	lastTableFilter = tableClone
	return tableClone
}

// Init 初始方法
func (m *TableFilter) Init() tea.Cmd {
	return nil
}

// 创建一个新的TableFilter
func newTableFilter() *TableFilter {
	m := &TableFilter{}

	// 输入框
	m.Input = textinput.New()
	m.Input.Placeholder = "Press / for search"
	m.Input.CharLimit = 156

	m.Table = NewTable()

	return m
}

// View 渲染样式
func (m *TableFilter) View() string {
	var tpl string
	tpl += "%s\n\n"
	tpl += "%s\n\n"

	if len(m.SubDescs) > 0 {
		for idx, desc := range m.SubDescs {
			tpl += SubtleStyle.Render(desc)
			if idx != len(m.SubDescs)-1 {
				tpl += DotStyle
			}
		}
	}

	tableView := BodyStyle.Render(m.Table.View())

	s := fmt.Sprintf(tpl, m.Input.View(), tableView)
	return MainStyle.Render(s)
}

// Update 更新事件
func (m *TableFilter) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// 业务自己实现的更新事件
	if m.UpdateEvent != nil {
		updateModel, updateCmd := m.UpdateEvent(msg)
		if updateModel != nil || updateCmd != nil {
			return updateModel, updateCmd
		}
	}

	switch msgType := msg.(type) {
	// 按键事件
	case tea.KeyMsg:
		switch msgType.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			m.Table.Focus()
			m.Input.Blur()

		default:
			// 输入框指令，过滤表格内容
			if m.Input.Focused() {
				m.Input, _ = m.Input.Update(msg)
				value := m.Input.Value()
				m.Table.FilterRows(value)
				switch msgType.String() {
				case "enter":
					m.Table.Focus()
					m.Input.Blur()
				}
			}

			// Table操作指令
			if m.Table.Focused() {
				switch msgType.String() {
				// 选择控制按键
				case "j", "down", "k", "up", "b", "pgup", "f", "pgdown", "u", "ctrl+u", "d", "ctrl+d", "g", "home", "G", "end":
					m.Table.Model, _ = m.Table.Model.Update(msg)
				// 排序按钮
				case "1", "2", "3", "4", "5", "6", "7", "8", "9", "0":
					sortIndex, errAtoi := strconv.Atoi(msgType.String())
					if errAtoi != nil {
						logrus.Fatal(errAtoi)
					}
					m.Table.SetSortIndex(sortIndex)

				// 搜索按钮
				case "/":
					m.Input.Focus()
					m.Table.Model.Blur()
				}
			}
		}

	// 窗口变化事件
	case tea.WindowSizeMsg:
		m.Table.AutoResize(msgType)
		m.Table.AutoColSize()
	}

	return m, nil
}
