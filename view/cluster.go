package view

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/flyhope/kubetea/comm"
	"github.com/flyhope/kubetea/k8s"
	"github.com/flyhope/kubetea/lang"
	"github.com/flyhope/kubetea/ui"
	"github.com/nicksnyder/go-i18n/v2/i18n"
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
	//c.SubDescs = []string{fmt.Sprintf("数据更新时间：%s", k8s.PodCache().CreatedAt.Format(time.DateTime))}
	c.SubDescs = []string{lang.Txt(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "cluster.data_update_time",
			Other: "Data update time: {{.UpdateTime}}",
		},
		TemplateData: map[string]interface{}{
			"UpdateTime": k8s.PodCache().CreatedAt.Format(time.DateTime),
		},
	})}
}

// ShowCluster 获取k8s Pod列表
func ShowCluster() (tea.Model, error) {
	// 渲染UI
	m := &clusterModel{
		TableFilter: ui.FetchTableFilter(),
	}
	m.TableFilter.SetColumns(comm.ShowKubeteaConfig().ShowTemplateColumn(comm.ConfigTemplateCluster))
	m.Table.SetSortIndex(comm.ShowKubeteaConfig().Sort.Cluster)
	m.updateData(false)

	m.UpdateEvent = func(msg tea.Msg) (tea.Model, tea.Cmd) {
		switch msgType := msg.(type) {
		// 按键事件
		case tea.KeyMsg:
			switch msgType.String() {
			case "f5", "ctrl+r":
				m.updateData(true)
			}

			// 仅在未输入状态下，响应按键事件
			if !m.TableFilter.Input.Focused() {
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
					return model, nil
				}
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
