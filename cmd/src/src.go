package src

import (
	"fmt"
	"github.com/urfave/cli"
	"os"
)

var (
	app = cli.NewApp()
)

func init() {
	app.Action = srccoin
	app.HideVersion = true // we have a command to print the version
	app.Copyright = " Copyright(c)2017-2020 SRCChain"

	app.Commands = []cli.Command{
		consoleCommand,
	}
	app.Before = func(ctx *cli.Context) error {
		return nil
	}
}

func srccoin(ctx *cli.Context) {
	fmt.Println(111111)
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
