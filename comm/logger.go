package comm

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"time"
)

// LogSetStdout 设置日志输出为标准输出
func LogSetStdout() {
	logSetDefault()
	logrus.StandardLogger().ExitFunc = os.Exit
	logrus.SetOutput(os.Stdout)
}

// 写入日志的内容
type logFileData struct {
	Data    [][]byte
	MaxSize int
}

func (l *logFileData) Write(p []byte) (n int, err error) {
	l.Data = append(l.Data, p)
	diffSize := len(l.Data) - l.MaxSize
	if diffSize > 0 {
		l.Data = l.Data[diffSize:]
	}
	return len(p), nil
}

func (l *logFileData) String() string {
	var str string
	for _, v := range l.Data {
		str += string(v)
	}
	return str
}

var LogFileData = &logFileData{
	MaxSize: 100,
}

// LogSetFile 设置日志输出为文件
func LogSetFile() {
	logrus.SetLevel(logrus.Level(ShowKubeteaConfig().Log.Level))

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

	// 退出时看有没有日志，如果有，输出
	logrus.StandardLogger().ExitFunc = func(code int) {
		if errRelease := Program.ReleaseTerminal(); errRelease != nil {
			logrus.Warnln(errRelease)
		}

		LogSetStdout()
		if len(LogFileData.Data) > 0 {
			if _, errFprint := fmt.Fprint(os.Stderr, LogFileData.String()); errFprint != nil {
				logrus.Errorln(errFprint)
			}
		}
		os.Exit(code)
	}

	// 组合多个IO写入器，写入文件的同时记录下来
	writer := io.MultiWriter(fo, LogFileData)

	// 设置日志输出路径
	logSetDefault()
	logrus.SetOutput(writer)
}

// 设置默认的日志配置
func logSetDefault() {
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: time.DateTime,
	})
}
