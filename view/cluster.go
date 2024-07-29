package view

import (
	"fmt"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/flyhope/kubetea/comm"
	"github.com/flyhope/kubetea/k8s"
	"github.com/flyhope/kubetea/ui"
	"github.com/sirupsen/logrus"
	"path/filepath"
	"sort"
	"strconv"
	"time"
)

type clusterModel struct {
	ui.Abstract
	*ui.TableFilter
}

// 更新数据
func (c *clusterModel) updateData(force bool) {
	config := comm.ShowKConfig()
	groupLabels := make(map[string]int)

	pods, err := k8s.PodCache().ShowList(force)
	if err != nil {
		return
	}

	// 获取要展示的label
	labelTotal := make(map[string]int)
	for _, pod := range pods.Items {
		groupLabelValue := pod.Labels[config.ClusterByLabel]
		for idx, filter := range config.ClusterFilters {
			ok, errMatch := filepath.Match(filter, groupLabelValue)
			if errMatch != nil {
				logrus.Warnln(errMatch)
				continue
			}

			if ok {
				groupLabels[groupLabelValue] = idx
				labelTotal[groupLabelValue]++
			}
		}
	}

	// 拼接UI列表数据
	rows := make([]table.Row, 0, len(groupLabels))
	for label := range groupLabels {
		rows = append(rows, table.Row{label, strconv.Itoa(labelTotal[label])})
	}
	sort.Slice(rows, func(i, j int) bool {
		if groupLabels[rows[i][0]] != groupLabels[rows[j][0]] {
			return groupLabels[rows[i][0]] < groupLabels[rows[j][0]]
		}

		return rows[i][0] < rows[j][0]
	})
	c.Table.SetRows(rows)
	c.SubDescs = []string{fmt.Sprintf("数据更新时间：%s", k8s.PodCache().CreatedAt.Format(time.DateTime))}
}

// ShowCluster 获取k8s Pod列表
func ShowCluster() (tea.Model, error) {
	// 渲染UI
	m := &clusterModel{
		TableFilter: ui.NewTableFilter(),
	}
	m.TableFilter.Table = ui.NewTableWithData([]table.Column{
		{Title: "集群", Width: 0},
		{Title: "数量", Width: 10},
	}, nil)
	m.updateData(false)

	m.UpdateEvent = func(msg tea.Msg) (tea.Model, tea.Cmd) {
		switch msgType := msg.(type) {
		// 按键事件
		case tea.KeyMsg:
			switch msgType.String() {
			case "enter":
				row := m.Table.SelectedRow()
				if len(row) == 0 {
					break
				}
				model, err := ShowPod(row[0], m)
				if err != nil {
					logrus.Fatal(err)
				}

				return ui.ViewModel(model)
			case "f5", "ctrl+r":
				m.updateData(true)
			}
		case comm.MsgPodCache, comm.MsgUIBack:
			m.updateData(false)
		}

		return nil, nil
	}
	return m, nil
}
