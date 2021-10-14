package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/vvatanabe/tf2b/cli"
	"github.com/vvatanabe/tfnotify/errors"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	app := cli.New()
	err := app.RunContext(ctx, os.Args)
	code := errors.HandleExit(err)
	os.Exit(code)
}
