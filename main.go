package main

import (
	"github.com/flyhope/kubetea/comm"
	"github.com/flyhope/kubetea/view"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	// 日志配置
	comm.LogSetStdout()

	// 定义入口
	app := &cli.App{
		Name: "kubernetes simple cli ui client",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "config", Aliases: []string{"c"}, Value: "~/.kubetea/config.yaml", Usage: "(optional) path to the kubetea.yaml config file"},
			&cli.StringFlag{Name: "namespace", Aliases: []string{"n"}},
			&cli.StringFlag{Name: "context"},
			&cli.StringFlag{Name: "kubeconfig", Aliases: []string{"k"}, Usage: "(optional) absolute path to the kubeconfig file"},
		},
		Commands: router(),
		Before: func(context *cli.Context) error {
			comm.Context = context
			return nil
		},
		Action: view.Action,
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}
