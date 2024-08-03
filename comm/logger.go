package comm

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

// LogSetStdout 设置日志输出为标准输出
func LogSetStdout() {
	logSetDefault()
	logrus.SetOutput(os.Stdout)
}

// LogSetFile 设置日志输出为文件
func LogSetFile() {
	logrus.SetLevel(ShowKubeteaConfig().Log.Level)

	// 创建日志目录
	dir := FixPath(ShowKubeteaConfig().Log.Dir)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		logrus.Fatal(err)
	}

	// 扫描目录，删除超过指定数量的日志文件
	errDel := DeleteFilesWhtiKeep(dir, ShowKubeteaConfig().Log.FileTotalMax)

	// 打开日志文件
	logFilePath := fmt.Sprintf("%s/%s.log", dir, time.Now().Format(time.DateOnly))
	fo, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0664)
	if err != nil {
		if errDel != nil {
			logrus.WithFields(logrus.Fields{"dir": dir}).Errorln(errDel)
		}
		logrus.WithFields(logrus.Fields{"dir": dir}).Fatal(err)
	}

	// 设置日志输出路径
	logSetDefault()
	logrus.SetOutput(fo)
}

// 设置默认的日志配置
func logSetDefault() {
	logrus.SetReportCaller(true)
}
