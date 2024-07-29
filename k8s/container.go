package k8s

import (
	"github.com/flyhope/kubetea/comm"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/rest"
	"os/exec"
)

// ContainerLog 获取一个容器的日志
func ContainerLog(podName string, containerName string) *rest.Request {
	return Client().CoreV1().Pods(ShowNamespace()).GetLogs(podName, &v1.PodLogOptions{
		Container: containerName,
		TailLines: comm.Ptr[int64](1000),
	})
}

// ContainerShell 生成一个容器Shell命令
func ContainerShell(podName string, containerName string) *exec.Cmd {
	cmd := KubeCmd()
	cmd.Args = append(cmd.Args,
		"exec", "-it", "-c", containerName, podName,
		"--",
		"sh", "-c", "if command -v bash >/dev/null 2>&1; then bash; else sh; fi",
	)
	return cmd
}
