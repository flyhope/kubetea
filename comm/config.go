package comm

import (
	_ "embed"
	"github.com/charmbracelet/bubbles/table"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"sync"
)

//go:embed kubetea.yaml
var kubeteaDefaultYaml []byte

// ShowKubeteaConfig 获取k.yaml配置
var ShowKubeteaConfig = sync.OnceValue(func() *KubeteaConfig {
	config := new(KubeteaConfig)

	// 文件件在才继续加载配置
	configFilePath := FixPath(Context.String("config"))

	yamlContent := kubeteaDefaultYaml
	_, errStat := os.Stat(configFilePath)
	if os.IsNotExist(errStat) {
		// 文件不存在，加载并写入一份默认配置
		configFileDir := filepath.Dir(configFilePath)
		if errMkdir := os.MkdirAll(configFileDir, os.ModePerm); errMkdir != nil {
			logrus.Warnln(errMkdir)
		} else {
			// 写入一份默认配置
			errWrite := os.WriteFile(configFilePath, kubeteaDefaultYaml, 0664)
			if errWrite != nil {
				logrus.WithFields(logrus.Fields{"path": configFilePath}).Warnln(errWrite)
			}
		}
	} else {
		// 文件存在，读取YAML文件内容
		if fileContent, err := os.ReadFile(configFilePath); err != nil {
			logrus.WithFields(logrus.Fields{"config": configFilePath}).Fatalf("Failed to read YAML file: %v", err)
		} else {
			yamlContent = fileContent
		}
	}

	// 解析YAML文件
	err := yaml.Unmarshal(yamlContent, &config)
	if err != nil {
		logrus.Fatalf("Failed to unmarshal YAML: %v", err)
	}

	return config
})

// KubeteaConfig YAML配置定义
type KubeteaConfig struct {
	PodCacheLivetime uint32            `yaml:"pod_cache_livetime_second"` // 缓存Pod过期时间，过期自动刷新
	Log              KubeteaConfigLog  `yaml:"log"`                       // 日志配置
	ClusterByLabel   string            `yaml:"cluster_by_label"`          // 筛选显示cluster的Label的名称
	ClusterFilters   []string          `yaml:"cluster_filters"`           // 筛选显示cluster的Label的值，支持glob
	Sort             KubeteaConfigSort `yaml:"sort"`                      // 自定义排序
	Template         KubeteaTemplate   `yaml:"template"`                  // 显示模板定义
}

type KubeteaConfigLog struct {
	Dir          string `yaml:"dir"`
	FileTotalMax int    `yaml:"file_total_max"`
	Level        uint32 `yaml:"level"`
}

type KubeteaConfigSort struct {
	Container SortMap `yaml:"container"`
}

type KubeteaTemplate struct {
	Pod KubeteaTemplateTable `yaml:"pod"`
}

type KubeteaTemplateTable struct {
	Column []table.Column `yaml:"column"`
	Body   []string       `yaml:"body"`
}
