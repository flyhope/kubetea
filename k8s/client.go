package k8s

import (
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"os/exec"
	"sync"
)

// Client 获取k8s client
var Client = sync.OnceValue(func() *kubernetes.Clientset {
	config := ShowKubeConfig()

	// 创建 Kubernetes 客户端
	clientset, errKubernetes := kubernetes.NewForConfig(config)
	if errKubernetes != nil {
		logrus.Fatal(errKubernetes)
	}

	return clientset
})

// KubeCmd 生成kubectl基础命令，带有配置文件等通用参数
func KubeCmd() *exec.Cmd {
	cmd := exec.Command("kubectl",
		"--kubeconfig", ShowKubeConfigPath(),
		"--namespace", ShowNamespace(),
		"--context", ShowContext(),
	)
	return cmd
}
