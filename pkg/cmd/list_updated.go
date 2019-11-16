package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/scylladb/go-set/strset"
	"github.com/suzuki-shunsuke/go-cliutil"
	"github.com/urfave/cli"

	"github.com/suzuki-shunsuke/candy/pkg/config"
)

type (
	Params struct {
		Args    cli.Args
		IsDebug bool
	}
)

// ListUpdated is the sub command "list-updated".
func ListUpdated(c *cli.Context) error {
	return cliutil.ConvErrToExitError(listUpdated(
		Params{
			Args:    c.Args(),
			IsDebug: c.Bool("debug"),
		},
	))
}

func listUpdated(params Params) error {
	printD := func(args ...interface{}) {
		if params.IsDebug {
			fmt.Fprintln(os.Stderr, args...)
		}
	}
	cfg := config.Config{}
	if err := config.Read(".candy.yaml", &cfg); err != nil {
		return err
	}

	envs := os.Environ()

	for _, target := range cfg.Targets {
		printD("target: " + target)
		srvCfg := config.ServiceConfig{}
		if err := config.ReadService(filepath.Join(target, ".candy.yaml"), &srvCfg); err != nil {
			return err
		}
		for _, task := range srvCfg.Tasks {
			printD("target: " + target + ", task: " + task.Name)
			paths := strset.New()
			if len(task.Files) != 0 {
				for _, file := range task.Files {
					if len(file.Paths) != 0 {
						if file.Excluded {
							paths.Remove(file.Paths...)
						} else {
							paths.Add(file.Paths...)
						}
						continue
					}
					if file.Command != "" {
						cmd := exec.Command("sh", "-c", file.Command)
						cmd.Env = envs
						var stdout bytes.Buffer
						if params.IsDebug {
							cmd.Stdout = io.MultiWriter(&stdout, os.Stderr)
							cmd.Stderr = os.Stderr
						} else {
							cmd.Stdout = &stdout
						}
						printD("+ "+cmd.String(), "# list files which the task depends on")
						if err := cmd.Run(); err != nil {
							return err
						}
						if file.Excluded {
							paths.Remove(strings.Split(stdout.String(), "\n")...)
						} else {
							paths.Add(strings.Split(stdout.String(), "\n")...)
						}
					}
				}
			}
			if paths.Size() == 0 {
				paths.Add(target)
			}
			updated := false
			switch {
			case task.Change.IsFilesChanged.Command != "":
				cmd := exec.Command("sh", "-c", task.Change.IsFilesChanged.Command+" "+strings.Join(paths.List(), " "))
				cmd.Env = envs
				if params.IsDebug {
					cmd.Stdout = os.Stderr
					cmd.Stderr = os.Stderr
				}
				printD("+ "+cmd.String(), "# list changed files")
				if err := cmd.Run(); err != nil {
					updated = true
				}
			default:
				cmd := exec.Command("sh", "-c", "git diff --quiet origin/master HEAD "+strings.Join(paths.List(), " "))
				cmd.Env = envs
				if params.IsDebug {
					cmd.Stdout = os.Stderr
					cmd.Stderr = os.Stderr
				}
				printD("+ "+cmd.String(), "# list changed files")
				if err := cmd.Run(); err != nil {
					updated = true
				}
			}
			if updated {
				fmt.Println(target + ":" + task.Name)
			}
		}
	}

	return nil
}
