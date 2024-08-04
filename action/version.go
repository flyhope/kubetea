package action

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"runtime"
	"strconv"
	"time"
)

// 展示版本信息
var (
	Ver       string
	BuildTime string
	GitCommit string
)

// Version 获取版本信息
func Version(c *cli.Context) error {
	timestamp, err := strconv.ParseInt(BuildTime, 10, 64)
	if err != nil {
		logrus.WithFields(logrus.Fields{"build-time": BuildTime}).Warnln(err)
	}
	buildTime := time.Unix(timestamp, 0)

	fmt.Printf("Version: %s\n", Ver)
	fmt.Printf("Build Time: %s\n", buildTime.Format(time.DateTime))
	fmt.Printf("Git Commit: %s\n", GitCommit)
	fmt.Printf("Go Version: %s\n", runtime.Version())
	return nil
}
