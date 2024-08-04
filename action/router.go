package action

import "github.com/urfave/cli/v2"

func Router() []*cli.Command {
	return []*cli.Command{
		{
			Name:   "version",
			Usage:  "show version",
			Action: Version,
		},
	}
}
