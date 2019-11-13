package cmd

import (
	"context"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/suzuki-shunsuke/go-cliutil"
	"github.com/suzuki-shunsuke/go-timeout/timeout"
	"github.com/urfave/cli"
)

type (
	Params struct {
		Args cli.Args
	}
)

// Run is the sub command "run".
func Run(c *cli.Context) error {
	return cliutil.ConvErrToExitError(run(
		Params{
			Args: c.Args(),
		},
	))
}

func run(params Params) error {
	ctx := context.Background()

	cmd := exec.Command("echo", "hello")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	signalChan := make(chan os.Signal, 1)
	signal.Notify(
		signalChan, syscall.SIGHUP, syscall.SIGINT,
		syscall.SIGTERM, syscall.SIGQUIT)

	runner := timeout.NewRunner(0)

	sentSignals := map[os.Signal]struct{}{}
	exitChan := make(chan error, 1)

	go func() {
		exitChan <- runner.Run(ctx, cmd)
	}()

	for {
		select {
		case err := <-exitChan:
			return err
		case sig := <-signalChan:
			if _, ok := sentSignals[sig]; ok {
				continue
			}
			sentSignals[sig] = struct{}{}
			runner.SendSignal(sig.(syscall.Signal))
		}
	}
}
