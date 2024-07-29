package view

import (
	"github.com/flyhope/kubetea/ui"
	"github.com/urfave/cli/v2"
)

// Action 主入口
func Action(c *cli.Context) error {
	m, err := ShowCluster()
	if err != nil {
		return err
	}

	_, err = ui.RunProgram(m)
	return err
}
