package comm

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type FileInfo struct {
	Name    string
	ModTime time.Time
}

// DeleteFilesWhtiKeep 删除一个目录下的文件，保留指定个数
func DeleteFilesWhtiKeep(directory string, keep int) error {
	// 读取目录中的所有文件信息
	entries, err := os.ReadDir(directory)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	// 将文件信息存储在一个切片中
	var fileInfos []FileInfo
	for _, entry := range entries {
		if !entry.IsDir() {
			info, errInfo := entry.Info()
			if errInfo != nil {
				return fmt.Errorf("failed to get file info: %w", errInfo)
			}
			fileInfos = append(fileInfos, FileInfo{Name: entry.Name(), ModTime: info.ModTime()})
		}
	}

	// 如果文件数量小于或等于要保留的数量，则无需删除任何文件
	if len(fileInfos) <= keep {
		return nil
	}

	// 按修改时间升序排序
	sort.Slice(fileInfos, func(i, j int) bool {
		return fileInfos[i].ModTime.Before(fileInfos[j].ModTime)
	})

	// 删除最早的文件，直到只剩下保留的数量
	for i := 0; i < len(fileInfos)-keep; i++ {
		errRemove := os.Remove(directory + "/" + fileInfos[i].Name)
		if errRemove != nil {
			return fmt.Errorf("failed to delete file %s: %w", fileInfos[i].Name, errRemove)
		}
	}

	return nil
}

// FixPath 修正路径，支持 ~
func FixPath(configFilePath string) string {
	if strings.HasPrefix(configFilePath, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			logrus.Warnln(err)
		} else {
			configFilePath = filepath.Join(homeDir, configFilePath[1:])
		}
	}
	return configFilePath
}
