package view

import (
	"bytes"
	"fmt"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/flyhope/kubetea/comm"
	"github.com/flyhope/kubetea/k8s"
	"github.com/flyhope/kubetea/ui"
	"github.com/sirupsen/logrus"
	"io"
	v1 "k8s.io/api/core/v1"
	"sort"
	"time"
)

// 容器列表页面
type containerModel struct {
	ui.Abstract
	*ui.TableFilter
	PodName string
}

func (m *containerModel) updateData(force bool) {
	pod, lastUpdate, err := k8s.PodCache().Show(m.PodName, force)
	if err != nil {
		logrus.Fatal(err)
	}

	// 获取 container
	rows := make([]table.Row, 0, len(pod.Status.ContainerStatuses))
	for _, container := range pod.Status.ContainerStatuses {
		rows = append(rows, table.Row{
			container.Name,
			container.Image,
			getContainerState(container),
			getBoolString(container.Ready),
		})
	}
	sort.Slice(rows, func(i, j int) bool {
		return rows[i][0] < rows[j][0]
	})

	// 展示 init container
	initRows := make([]table.Row, 0, len(pod.Status.InitContainerStatuses))
	for _, container := range pod.Status.InitContainerStatuses {
		initRows = append(initRows, table.Row{
			container.Name,
			container.Image,
			getContainerState(container),
			getBoolString(container.Ready),
		})
	}
	sort.Slice(initRows, func(i, j int) bool {
		return initRows[i][0] < initRows[j][0]
	})
	rows = append(rows, initRows...)

	m.Table.SetRows(rows)
	m.SubDescs = []string{
		fmt.Sprintf("合计：%d", len(rows)),
		fmt.Sprintf("数据更新时间：%s", lastUpdate.Format(time.DateTime)),
	}
}

func ShowContainer(podName string, lastModel tea.Model) (tea.Model, error) {
	// 渲染UI
	m := &containerModel{
		Abstract:    ui.Abstract{LastModel: lastModel},
		TableFilter: ui.NewTableFilter(),
		PodName:     podName,
	}
	m.Table = ui.NewTableWithData([]table.Column{
		{Title: "容器名称", Width: 0},
		{Title: "镜像地址", Width: 0},
		{Title: "状态", Width: 4},
		{Title: "就绪", Width: 4},
	}, nil)
	m.updateData(false)

	m.UpdateEvent = func(msg tea.Msg) (tea.Model, tea.Cmd) {
		row := m.Table.SelectedRow()
		switch msgType := msg.(type) {
		// 按键事件
		case tea.KeyMsg:
			switch msgType.String() {
			case "alt+left", "ctrl+left":
				return m.GoBack()
			}
		// 数据更新事件
		case comm.MsgPodCache, comm.MsgUIBack:
			m.updateData(false)
		}

		// 仅在未输入状态下，响应按键事件
		if !m.TableFilter.Input.Focused() {
			switch msgType := msg.(type) {
			// 按键事件
			case tea.KeyMsg:
				switch msgType.String() {
				// 返回上一级
				case "esc":
					return m.GoBack()
				// 进入容器Shell
				case "enter", "s":
					return m, ui.NewCmd(k8s.ContainerShell(m.PodName, row[0]))
				// 查看日志
				case "l":
					containerLog := k8s.ContainerLog(m.PodName, row[0])
					podLogs, err := containerLog.Stream(comm.Context.Context)
					if err != nil {
						logrus.Fatalf("Error getting logs: %s\n", err.Error())
					}
					defer podLogs.Close()

					buf := new(bytes.Buffer)
					_, err = io.Copy(buf, podLogs)
					if err != nil {
						logrus.Fatalf("Error copying logs: %s\n", err.Error())
					}
					return ui.PageViewContent(m.PodName, buf.String(), m), nil
				}
			}
		}
		return nil, nil
	}

	return m, nil

}

// 获取字符串输出的Bool值
func getBoolString(val bool) string {
	if val {
		return "✔️"
	}
	return "❌️"
}

// 获取容器的状态名称
func getContainerState(status v1.ContainerStatus) string {
	if status.State.Waiting != nil {
		return "♾️"
	} else if status.State.Terminated != nil {
		return "✴️"
	} else if status.State.Running != nil {
		return "✔️"
	}
	return "❓️"
}
