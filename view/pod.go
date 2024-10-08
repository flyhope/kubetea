package view

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/flyhope/kubetea/comm"
	"github.com/flyhope/kubetea/k8s"
	"github.com/flyhope/kubetea/lang"
	"github.com/flyhope/kubetea/ui"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"sort"
	"time"
)

// POD列表页面
type podModel struct {
	ui.Abstract
	*ui.TableFilter
	app string
}

// 更新数据
func (c *podModel) updateData(force bool) {
	pods, err := k8s.PodCache().ShowList(force)
	if err != nil {
		logrus.Warnln(err)
		return
	}

	rows := make([]table.Row, 0, len(pods.Items))
	for _, pod := range pods.Items {
		if pod.Labels["app"] == c.app {
			rows = append(rows, TemplateRenderBody(comm.ConfigTemplatePod, pod))

			//name := pod.Name
			////if strings.Index(name, app) == 0 {
			////	name = name[len(app):]
			////	name = strings.TrimLeft(name, "-_.")
			////}
		}
	}

	sort.Slice(rows, func(i, j int) bool {
		return rows[i][0] < rows[j][0]
	})

	c.Table.SetRows(rows)
	c.SubDescs = []string{
		lang.Data(langTotalWithNumber, lang.Map{"number": len(rows)}),
		lang.Data(langUpdateTime, lang.Map{"UpdateTime": k8s.PodCache().CreatedAt.Format(time.DateTime)}),
	}
}

// ShowPod 获取POD列表
func ShowPod(app string, lastModel tea.Model) (tea.Model, error) {
	// 渲染UI
	m := &podModel{
		Abstract:    ui.Abstract{LastModel: lastModel},
		TableFilter: ui.FetchTableFilter(),
		app:         app,
	}
	m.Abstract.Model = m

	m.TableFilter.Table.SetColumns(TemplateRenderColumn(comm.ConfigTemplatePod))
	m.Table.SetSortIndex(comm.ShowKubeteaConfig().Sort.Cluster)
	m.TableFilter.Focus()
	m.updateData(false)

	m.UpdateEvent = func(msg tea.Msg) (tea.Model, tea.Cmd) {
		switch msgType := msg.(type) {
		// 按键事件
		case tea.KeyMsg:
			switch msgType.String() {
			case "alt+left", "ctrl+left":
				return m.GoBack()
			case "f5", "ctrl+r":
				m.updateData(true)
			}

			// 仅在未输入状态下，响应按键事件
			if !m.TableFilter.Input.Focused() {
				switch msgType.String() {

				// 返回上一级
				case "esc":
					return m.GoBack()

				// 打开容列表
				case "enter":
					row := m.Table.SelectedRow()
					if len(row) == 0 {
						break
					}
					model, err := ShowContainer(row[0], m)
					if err != nil {
						logrus.Fatal(err)
					}
					return model, nil

				// 查看JSON数据
				case "i":
					row := m.Table.SelectedRow()
					pod, _, err := k8s.PodCache().Show(row[0], false)
					if err != nil {
						logrus.Fatal(err)
					}
					return ui.ViewModel(ui.PageViewJson(row[0], pod, m.TableFilter))

				// 查看 Describe
				case "e":
					return ui.NewCmdPause(m, k8s.KubeCmdArgs("describe", "pod", m.Table.SelectedRow()[0]))
				}
			}

		case comm.MsgPodCache, comm.MsgUIBack:
			m.updateData(false)
		}

		return nil, nil
	}

	return m, nil
}

var phaseAlias = map[v1.PodPhase]string{
	v1.PodPending:   "♾️",
	v1.PodRunning:   "✔️",
	v1.PodSucceeded: "🔅",
	v1.PodFailed:    "❌️",
	v1.PodUnknown:   "❓️",
	"Terminating":   "✴️",
}

// podPhaseView 友好显示POD状态
func podPhaseView(pod v1.Pod) string {
	phase := pod.Status.Phase
	if pod.DeletionTimestamp != nil {
		phase = "Terminating"
	}

	result := phaseAlias[phase]
	if result == "" {
		result = string(phase)
	}
	return result
}

// podReadyView 友好显示POD的Ready状态
func podReadyView(pod v1.Pod) string {
	for _, condition := range pod.Status.Conditions {
		if condition.Type == v1.PodReady && condition.Status == v1.ConditionTrue {
			return "✔️"
		}
	}
	return "❌️"
}
