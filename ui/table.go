package ui

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/flyhope/kubetea/comm"
	"slices"
	"sort"
	"strings"
	"unicode/utf8"
)

// Table 加强版
// 1. 自适应宽高
// 2. 支持筛选
type Table struct {
	table.Model
	filterValue   string         // 过滤字符表达式
	originRows    []table.Row    // 存储原始未筛选的数据
	originColumns []table.Column // 存储原始未筛选的表头
}

// SetRows 设置完整数据
func (t *Table) SetRows(rows []table.Row) {
	t.originRows = rows
	t.AutoColSize()
	t.FilterRows(t.filterValue)
}

// SetColumns 设置表格字段
func (t *Table) SetColumns(c []table.Column) {
	t.originColumns = c
	if len(t.originRows) > 0 {
		t.AutoColSize()
	}
}

// FilterRows 根据输入的筛选条件，筛选数据
func (t *Table) FilterRows(value string) {
	t.filterValue = value

	if t.filterValue == "" {
		t.Model.SetRows(t.originRows)
	} else {
		rows := make([]table.Row, 0, len(t.Model.Rows()))
		for _, row := range t.originRows {
			for _, item := range row {
				if strings.Contains(strings.ToLower(item), strings.ToLower(t.filterValue)) {
					rows = append(rows, row)
					break
				}
			}
		}
		t.Model.SetRows(rows)
	}
}

// AutoResize 自动设置Table大小
func (t *Table) AutoResize(msg tea.WindowSizeMsg) {
	width := msg.Width - 4
	t.Model.SetWidth(width)
	t.Model.SetHeight(msg.Height - 8)
}

// AutoColSize 自动计算列宽
func (t *Table) AutoColSize() {
	cols := slices.Clone(t.originColumns)
	surplusWidth := t.Model.Width()                        // 剩余宽度
	autoColsIdx := make(map[int]int, len(t.originColumns)) // 自动宽度的索列数 => 最大字符数
	for idx, col := range cols {
		if col.Width == 0 {
			var colMaxWidth int
			for _, row := range t.originRows {
				colMaxWidth = max(colMaxWidth, utf8.RuneCountInString(row[idx]))
			}
			autoColsIdx[idx] = colMaxWidth
		} else {
			surplusWidth -= col.Width
		}
		surplusWidth -= 2
	}
	// 获取最大字符累加总数
	var totalStringCount int
	for _, val := range autoColsIdx {
		totalStringCount += val
	}
	// 计算每列数值
	for idx, val := range autoColsIdx {
		ratio := float64(val) / float64(totalStringCount)
		cols[idx].Width = max(int(float64(surplusWidth)*ratio), 3)
	}
	t.Model.SetColumns(cols)
}

// NewTable 带着数据和默认样式创建Table
func NewTable() Table {
	model := table.New()

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true).
		Align(lipgloss.Center)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)

	model.SetStyles(s)
	model.Focus()

	t := &Table{
		Model: model,
	}

	t.AutoResize(WindowSize())

	return *t
}

// TableRowsSort 对Rows进行指定名称排序
func TableRowsSort(rows []table.Row, sortByMap comm.SortMap) {
	sort.Slice(rows, func(i, j int) bool {
		sortValueI := sortByMap.SortVal(rows[i][0])
		sortValueJ := sortByMap.SortVal(rows[j][0])

		if sortValueI != sortValueJ {
			return sortValueI < sortValueJ
		}

		return rows[i][0] < rows[j][0]
	})
}
