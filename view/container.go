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
		rows = append(rows, TemplateRender(comm.ConfigTemplateContainer, container))
	}
	ui.TableRowsSort(rows, comm.ShowKubeteaConfig().Sort.Container)

	// 展示 init container
	initRows := make([]table.Row, 0, len(pod.Status.InitContainerStatuses))
	for _, container := range pod.Status.InitContainerStatuses {
		initRows = append(initRows, TemplateRender(comm.ConfigTemplateContainer, container))
	}
	ui.TableRowsSort(initRows, comm.ShowKubeteaConfig().Sort.Container)
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
	m.Abstract.Model = m

	m.Table = ui.NewTableWithData(comm.ShowKubeteaConfig().ShowTemplateColumn(comm.ConfigTemplateContainer), nil)
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
				// 查看JSON数据
				case "i":
					pod, _, err := k8s.PodCache().Show(m.PodName, false)
					if err != nil {
						logrus.Fatal(err)
					}
					return ui.ViewModel(ui.PageViewJson(pod.Name, pod, m))
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

// 获取容器的状态名称
func containerStateView(status v1.ContainerStatus) string {
	if status.State.Waiting != nil {
		return "♾️"
	} else if status.State.Terminated != nil {
		return "✴️"
	} else if status.State.Running != nil {
		return "✔️"
	}
	return "❓️"
}
