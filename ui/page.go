package ui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
	"github.com/hokaccha/go-prettyjson"
	"github.com/sirupsen/logrus"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

type PageView struct {
	Title        string            // 标题
	Content      string            // 内容
	Referer      tea.Model         // 来源页面
	ready        bool              // 是否就绪
	viewport     viewport.Model    // 内容显示容器
	Input        textinput.Model   // 搜索条
	filterIndexs []*pageViewFilter // 过滤匹配的内容索引（倒序存储）
	filterOffset int               // 当前过滤的偏移量
}

// PageViewJson 以页面形式展示一个json
func PageViewJson(title string, content any, referer tea.Model) *PageView {
	jsonData, err := prettyjson.Marshal(content)
	if err != nil {
		logrus.Fatal(err)
	}
	return PageViewContent(title, string(jsonData), referer)
}

// PageViewContent 以页面形展形示一条内容
func PageViewContent(title, content string, referer tea.Model) *PageView {
	m := &PageView{
		Title:   title,
		Content: content,
		Referer: referer,
		Input:   textinput.New(),
	}

	size := WindowSize()
	m.Resize(size)
	return m
}

func (m *PageView) Init() tea.Cmd {
	return nil
}

// 过滤器位置记录
type pageViewFilter struct {
	start     int
	end       int
	line      int
	selfColor string
}

// 前往当前指定筛选的位置
func (m *PageView) gotoFilterOffset() {
	var filter *pageViewFilter
	if m.filterOffset < len(m.filterIndexs) {
		filter = m.filterIndexs[m.filterOffset]
	}
	if filter != nil {
		m.viewport.SetYOffset(filter.line)
	}
	m.updateViewContent()
}

// 更新筛选条件
func (m *PageView) updateInputFilter() {
	m.filterIndexs = nil
	filterContent := m.Input.Value()
	if strings.TrimSpace(m.Input.Prompt) == "/" && filterContent != "" {
		var filterOffset int
		regex, err := regexp.Compile(filterContent)
		if err != nil {
			logrus.WithFields(logrus.Fields{"filter": filterContent}).Warnln(err)
			m.viewport.SetContent(m.Content)
			return
		}
		indexs := regex.FindAllStringIndex(m.Content, -1)

		if len(indexs) > 0 {
			// 找到所有颜色控制字符
			cmdColorIndexs := rexexCmdColor.FindAllStringIndex(m.Content, -1)
			slices.Reverse(cmdColorIndexs)
			slices.Reverse(indexs)

		ForMatchIndex:
			for _, index := range indexs {
				var selfColor string
				for _, colorIndex := range cmdColorIndexs {
					// 找到离自己最近的颜色
					if index[0] >= colorIndex[0] {
						// 不处理颜色字符
						if index[1] <= colorIndex[1] {
							continue ForMatchIndex
						}
						selfColor = m.Content[colorIndex[0]:colorIndex[1]]
						break
					}
				}

				// 记录查找位置
				filterIndex := &pageViewFilter{
					start:     index[0],
					end:       index[1],
					line:      strings.Count(m.Content[:index[0]], "\n"),
					selfColor: selfColor,
				}
				if filterIndex.line >= m.viewport.YOffset {
					filterOffset = len(m.filterIndexs)
				}
				m.filterIndexs = append(m.filterIndexs, filterIndex)
			}
		}
		m.filterOffset = filterOffset
	}

	m.gotoFilterOffset()
}

// 颜色控制字符正则
var rexexCmdColor = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// 重置输入框
func (m *PageView) inputReset() {
	m.Input.Reset()
	m.Input.Blur()
	m.filterIndexs = nil
}

// Update 事件更新
func (m *PageView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msgType := msg.(type) {
	case tea.KeyMsg:
		if k := msgType.String(); k == "ctrl+c" {
			return ViewModel(m.Referer)
		}

		switch msgType.String() {
		default:
			if m.Input.Focused() {
				prompt := strings.TrimSpace(m.Input.Prompt)
				m.Input, cmd = m.Input.Update(msg)

				// 输入模式
				switch msgType.String() {
				case "esc":
					m.Input.SetValue("")
					m.Input.Blur()
					m.updateInputFilter()
					return m, tea.Sequence(tea.ExitAltScreen, tea.EnterAltScreen)
				default:
					switch prompt {
					// 输入命令模式
					case ":":
						switch msgType.String() {
						case "q":
							return ViewModel(m.Referer)
						case "enter":
							line, err := strconv.Atoi(m.Input.Value())
							if err != nil {
								logrus.Warnln(err)
							}
							m.viewport.SetYOffset(line)
							m.inputReset()
							cmds = append(cmds, tea.ExitAltScreen, tea.EnterAltScreen)
						}
					// 查找字符串模式
					case "/":
						switch msgType.String() {
						case "enter":
							m.Input.Blur()
						}
					}
					m.updateInputFilter()
				}

				return m, cmd
			} else {
				// 控制模式
				switch msgType.String() {
				case "g":
					m.viewport.GotoTop()
					return m, nil
				case "G":
					m.viewport.GotoBottom()
					return m, nil
				case ":": // 命令模式
					m.Input.Reset()
					m.Input.Prompt = ": "
					return m, m.Input.Focus()

				case "/": // 查找匹配模式
					m.Input.Reset()
					m.Input.Prompt = "/ "
					return m, m.Input.Focus()
				}

				// 查找中模式
				if m.Input.Value() != "" {
					// 查找字符串命令模式
					switch msgType.String() {
					// 下一个
					case "n":
						m.filterOffset--
						if m.filterOffset < 0 {
							m.filterOffset = len(m.filterIndexs) - 1
						}
						m.gotoFilterOffset()
					// 上一个
					case "N":
						m.filterOffset++
						if m.filterOffset >= len(m.filterIndexs) {
							m.filterOffset = 0
						}
						m.gotoFilterOffset()
					// 退出查找
					case "esc":
						m.inputReset()
						cmds = append(cmds, tea.ExitAltScreen, tea.EnterAltScreen)
					}
				}

				// Handle keyboard and mouse events in the viewport
				m.viewport, cmd = m.viewport.Update(msg)
				if key.Matches(msgType, m.viewport.KeyMap.Up) ||
					key.Matches(msgType, m.viewport.KeyMap.Down) ||
					key.Matches(msgType, m.viewport.KeyMap.HalfPageUp) ||
					key.Matches(msgType, m.viewport.KeyMap.HalfPageDown) ||
					key.Matches(msgType, m.viewport.KeyMap.PageUp) ||
					key.Matches(msgType, m.viewport.KeyMap.PageDown) {
					m.updateViewContent()
				}
				cmds = append(cmds, cmd)
			}
		}

	case tea.WindowSizeMsg:
		if cmdResize := m.Resize(msgType); cmdResize != nil {
			cmds = append(cmds, cmdResize)
		}
	}

	return m, tea.Sequence(cmds...)
}

// Resize 窗口大小重置时渲染
func (m *PageView) Resize(msgType tea.WindowSizeMsg) tea.Cmd {
	headerHeight := lipgloss.Height(m.headerView())
	verticalMarginHeight := headerHeight

	if !m.ready {
		// Since this program is using the full size of the viewport we
		// need to wait until we've received the window dimensions before
		// we can initialize the viewport. The initial dimensions come in
		// quickly, though asynchronously, which is why we wait for them
		// here.
		m.viewport = viewport.New(msgType.Width, msgType.Height-verticalMarginHeight)
		m.viewport.YPosition = headerHeight
		m.viewport.HighPerformanceRendering = false
		m.viewport.SetContent(m.Content)
		m.ready = true

		// This is only necessary for high performance rendering, which in
		// most cases you won't need.
		//
		// Render the viewport one line below the header.
		m.viewport.YPosition = headerHeight + 1
	} else {
		m.viewport.Width = msgType.Width
		m.viewport.Height = msgType.Height - verticalMarginHeight
	}

	return tea.Sequence(tea.ExitAltScreen, tea.EnterAltScreen)
}

// 展示内容
func (m *PageView) updateViewContent() {
	if len(m.filterIndexs) == 0 {
		m.viewport.SetContent(m.Content)
		return
	}

	lineStart := m.viewport.YOffset
	lineEnd := m.viewport.YOffset + m.viewport.Height
	content := m.Content

	// 高亮文本
	hilight := color.New(color.BgBlue)
	for _, filter := range m.filterIndexs {
		if filter.line < lineStart {
			break
		}
		if filter.line > lineEnd {
			continue
		}

		// 高亮前部分 + 高亮部分 + 原颜色（高亮会重置颜色） + 高亮后部分
		content = content[:filter.start] + hilight.Sprintf("%s", content[filter.start:filter.end]) + filter.selfColor + content[filter.end:]
	}

	m.viewport.SetContent(content)
	return
}

// View 展示数据
func (m *PageView) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}
	result := fmt.Sprintf("%s\n%s", m.headerView(), m.viewport.View())

	return result
}

var (
	headerStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	pageTitleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return headerStyle.BorderStyle(b).Padding(0, 1)
	}()

	pageInfoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return headerStyle.BorderStyle(b)
	}()
)

// 展示头部
func (m *PageView) headerView() string {
	// 输入模式，顶部显示输入框
	if m.Input.Focused() || m.Input.Value() != "" {
		var info string
		if strings.TrimSpace(m.Input.Prompt) == "/" {
			foundTotal := len(m.filterIndexs)
			if foundTotal == 0 {
				info = pageInfoStyle.Render("查找")
			} else {
				info = pageInfoStyle.Render(fmt.Sprintf("%d/%d", foundTotal-m.filterOffset, foundTotal))
			}
		} else {
			info = pageInfoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
		}
		title := pageTitleStyle.Width(m.viewport.Width - lipgloss.Width(info) - 2).Render(m.Input.View())
		return lipgloss.JoinHorizontal(lipgloss.Center, title, info)
	}

	// 页面展示模式，展示标题
	title := pageTitleStyle.Render(m.Title)
	info := pageInfoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line, info)
}
