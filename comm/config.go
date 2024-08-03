package comm

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"sync"
)

// ShowKubeteaConfig 获取k.yaml配置
var ShowKubeteaConfig = sync.OnceValue(func() *KubeteaConfig {
	config := &KubeteaConfig{
		ClusterByLabel:   "app",
		PodCacheLivetime: 10,
		Log: KubeteaConfigLog{
			Dir:          "~/.kubetea/logs",
			FileTotalMax: 10,
			Level:        logrus.InfoLevel,
		},
		ClusterFilters: []string{"*"},
	}

	// 文件件在才继续加载配置
	configFilePath := FixPath(Context.String("config"))

	_, errStat := os.Stat(configFilePath)
	if os.IsNotExist(errStat) {
		configFileDir := filepath.Dir(configFilePath)
		if errMkdir := os.MkdirAll(configFileDir, os.ModePerm); errMkdir != nil {
			logrus.Warnln(errMkdir)
		} else {
			// 文件不存在，写入一份默认配置
			yamlData, err := yaml.Marshal(&config)
			if err != nil {
				logrus.Warnln(err)
			} else {
				errWrite := os.WriteFile(configFilePath, yamlData, 0664)
				if errWrite != nil {
					logrus.WithFields(logrus.Fields{"path": configFilePath}).Warnln(errWrite)
				}
			}
		}

	} else {
		// 文件存在，读取YAML文件内容
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

// KubeteaConfig YAML配置定义
type KubeteaConfig struct {
	ClusterByLabel   string           `yaml:"cluster_by_label"`
	PodCacheLivetime uint32           `yaml:"pod_cache_livetime_second"`
	Log              KubeteaConfigLog `yaml:"log"`
	ClusterFilters   []string         `yaml:"cluster_filters"`
}

type KubeteaConfigLog struct {
	Dir          string       `yaml:"dir"`
	FileTotalMax int          `yaml:"file_total_max"`
	Level        logrus.Level `yaml:"level"`
}
