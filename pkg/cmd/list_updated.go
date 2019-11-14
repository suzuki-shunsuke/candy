package cmd

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/suzuki-shunsuke/go-cliutil"
	"github.com/urfave/cli"

	"github.com/suzuki-shunsuke/candy/pkg/config"
)

type (
	Params struct {
		Args cli.Args
	}
)

// ListUpdated is the sub command "list-updated".
func Run(c *cli.Context) error {
	return cliutil.ConvErrToExitError(listUpdated(
		Params{
			Args: c.Args(),
		},
	))
}

func listUpdated(params Params) error {
	cfg := config.Config{}
	if err := config.Read(".candy.yaml", &cfg); err != nil {
		return err
	}

	for _, target := range cfg.Targets {
		srvCfg := config.ServiceConfig{}
		if err := config.ReadService(filepath.Join(target, ".candy.yaml"), &srvCfg); err != nil {
			return err
		}
		for _, task := range srvCfg.Tasks {
			// fmt.Println("+ git diff --quiet origin/master HEAD " + target)
			cmd := exec.Command("git", "diff", "--quiet", "origin/master", "HEAD", target)
			if err := cmd.Run(); err != nil {
				// updated
				fmt.Println(target + ":" + task.Name)
			}
		}
	}

	return nil
}
