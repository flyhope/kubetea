package comm

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// ShowKConfig 获取k.yaml配置
var ShowKConfig = sync.OnceValue(func() *KConfig {
	config := &KConfig{
		ClusterByLabel:   "app",
		ClusterFilters:   []string{"*"},
		PodCacheLivetime: 10,
	}

	// 文件件在才继续加载配置
	configFilePath := Context.String("config")
	if strings.HasPrefix(configFilePath, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			logrus.Warnln(err)
		} else {
			configFilePath = filepath.Join(homeDir, configFilePath[1:])
		}
	}
	_, errStat := os.Stat(configFilePath)
	if !os.IsNotExist(errStat) {
		// 读取YAML文件内容
		yamlFile, err := os.ReadFile(configFilePath)
		if err != nil {
			logrus.WithFields(logrus.Fields{"config": configFilePath}).Fatalf("Failed to read YAML file: %v", err)
		}

		// 解析YAML文件
		err = yaml.Unmarshal(yamlFile, &config)
		if err != nil {
			logrus.Fatalf("Failed to unmarshal YAML: %v", err)
		}
	}

	return config
})

type KConfig struct {
	ClusterByLabel   string   `yaml:"cluster_by_label"`
	ClusterFilters   []string `yaml:"cluster_filters"`
	PodCacheLivetime uint32   `yaml:"pod_cache_livetime_second"`
}
