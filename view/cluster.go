package view

import (
	"fmt"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/flyhope/kubetea/comm"
	"github.com/flyhope/kubetea/k8s"
	"github.com/flyhope/kubetea/ui"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"path/filepath"
	"sort"
	"time"
)

type clusterModel struct {
	ui.Abstract
	*ui.TableFilter
}

// 用于显示的集群每行数据
type clusterRow struct {
	Name    string
	Pods    []*v1.Pod
	SortNum int
}

// 更新数据
func (c *clusterModel) updateData(force bool) {
	config := comm.ShowKubeteaConfig()

	pods, err := k8s.PodCache().ShowList(force)
	if err != nil {
		return
	}

	// 获取要展示的label
	clusterRows := make(map[string]*clusterRow, len(config.ClusterFilters))
	for _, pod := range pods.Items {
		groupLabelValue := pod.Labels[config.ClusterByLabel]
		for idx, filter := range config.ClusterFilters {
			ok, errMatch := filepath.Match(filter, groupLabelValue)
			if errMatch != nil {
				logrus.Warnln(errMatch)
				continue
			}

			if ok {
				row := clusterRows[groupLabelValue]
				if row == nil {
					row = &clusterRow{
						Name:    groupLabelValue,
						SortNum: idx,
					}
					clusterRows[groupLabelValue] = row
				}
				row.Pods = append(row.Pods, &pod)
				break
			}
		}
	}

	// 拼接UI列表数据
	rows := make([]table.Row, 0, len(clusterRows))
	for _, row := range clusterRows {
		rows = append(rows, TemplateRender(comm.ConfigTemplateCluster, row))
	}

	// 排序
	sort.Slice(rows, func(i, j int) bool {
		if clusterRows[rows[i][0]].SortNum != clusterRows[rows[j][0]].SortNum {
			return clusterRows[rows[i][0]].SortNum < clusterRows[rows[j][0]].SortNum
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
	m.TableFilter.Table = ui.NewTableWithData(comm.ShowKubeteaConfig().ShowTemplateColumn(comm.ConfigTemplateCluster), nil)
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

// countByPodLabelKeyValue 根据Pod的Label标签，统计符合条件的POD的数量
func countByPodLabelKeyValue(pods []*v1.Pod, key, value string) (result int) {
	for _, pod := range pods {
		if pod.Labels[key] == value {
			result++
		}
	}
	return result
}
