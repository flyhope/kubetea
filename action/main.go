package action

import (
	"errors"
	"github.com/flyhope/kubetea/comm"
	"github.com/flyhope/kubetea/k8s"
	"github.com/flyhope/kubetea/ui"
	"github.com/flyhope/kubetea/view"
	"github.com/urfave/cli/v2"
)

// Main 主入口
func Main(c *cli.Context) error {
	// 设置日志
	comm.LogSetFile()

	// 直接通过POD IP或名称进入POD的container列表
	args := c.Args()
	if args.Len() > 0 {
		return ActionPod(args.Get(0))
	}

	// 展示首屏内容
	m, err := view.ShowCluster()
	if err != nil {
		return err
	}

	_, err = ui.RunProgram(m)
	return err
}

// ActionPod 根据IP或Name直接打开Container
func ActionPod(ipOrName string) error {
	podList, err := k8s.PodCache().ShowList(false)
	if err != nil {
		return err
	}

	var podName string
	for _, pod := range podList.Items {
		// 匹配名称
		if pod.Name == ipOrName {
			podName = pod.Name
			break
		}
		// 匹配IP
		for _, ip := range pod.Status.PodIPs {
			if ip.IP == ipOrName {
				podName = pod.Name
				break
			}
		}
	}

	// not found pod
	if podName == "" {
		return errors.New("pod not found")
	}

	uiModel, errModel := view.ShowContainer(podName, nil)
	if errModel != nil {
		return errModel
	}

	_, err = ui.RunProgram(uiModel)
	return err
}
