package k8s

import (
	"github.com/flyhope/kubetea/comm"
	"github.com/sirupsen/logrus"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	"os"
	"path/filepath"
	"sync"
)

// ShowKubeConfigPath 获取kubeconfig的路径
var ShowKubeConfigPath = sync.OnceValue(func() string {
	// 定义 kubeconfig 路径参数
	kubeconfig := comm.Context.String("kubeconfig")
	if kubeconfig == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			logrus.Fatal(err)
		}
		kubeconfig = filepath.Join(home, ".kube", "config")
	}
	return kubeconfig
})

// ShowKubeConfig 获取kube的配置
var ShowKubeConfig = sync.OnceValue(func() *restclient.Config {

	configOverrides := &clientcmd.ConfigOverrides{
		CurrentContext: ShowContext(),
	}
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: ShowKubeConfigPath()},
		configOverrides,
	).ClientConfig()
	if err != nil {
		logrus.Fatalf("Error loading kubeconfig: %v\n", err)
	}

	return config
})

// ShowApiConfig 获取kube的配置
var ShowApiConfig = sync.OnceValue(func() *api.Config {
	configAccess := clientcmd.NewDefaultPathOptions()
	configAccess.GetLoadingPrecedence()
	config, err := configAccess.GetStartingConfig()
	if err != nil {
		logrus.Fatalf("show api config error: %v", err)
	}
	return config
})

// ShowNamespace 获取指定的的namespace
func ShowNamespace() string {
	// 手动参数指定namespace
	namespace := comm.Context.String("namespace")
	if namespace != "" {
		return namespace
	}

	config := ShowApiConfig()
	currentContext := ShowContext()
	contextDetails := config.Contexts[currentContext]
	if contextDetails == nil {
		logrus.Fatal("\"current context not found in kubeconfig")
		return ""
	}
	return contextDetails.Namespace
}

// ShowContext 获取指定的context
func ShowContext() string {
	// 手动参数指定namespace
	k8sContext := comm.Context.String("context")
	if k8sContext != "" {
		return k8sContext
	}

	config := ShowApiConfig()
	currentContext := config.CurrentContext
	return currentContext
}
