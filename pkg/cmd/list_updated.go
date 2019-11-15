package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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

	envs := os.Environ()

	for _, target := range cfg.Targets {
		srvCfg := config.ServiceConfig{}
		if err := config.ReadService(filepath.Join(target, ".candy.yaml"), &srvCfg); err != nil {
			return err
		}
		for _, task := range srvCfg.Tasks {
			// fmt.Println("+ git diff --quiet origin/master HEAD " + target)
			t := target
			if len(task.Files) != 0 {
				paths := []string{}
				for _, file := range task.Files {
					if len(file.Paths) != 0 {
						paths = append(paths, file.Paths...)
						continue
					}
					if file.Command != "" {
						cmd := exec.Command("sh", "-c", file.Command)
						cmd.Env = envs
						var stdout bytes.Buffer
						cmd.Stdout = &stdout
						if err := cmd.Run(); err != nil {
							return err
						}
						paths = append(paths, strings.Split(stdout.String(), "\n")...)
					}
				}
				t = strings.Join(paths, " ")
			}
			cmd := exec.Command("sh", "-c", "git diff --quiet origin/master HEAD "+t)
			cmd.Env = envs
			if err := cmd.Run(); err != nil {
				// updated
				fmt.Println(target + ":" + task.Name)
			}
		}
	}

	return nil
}
