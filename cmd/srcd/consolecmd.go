package tmp

import (
	"github.com/urfave/cli"
	"srcd/cmd/util"
	"srcd/cmd/console"
)
var(

	consoleFlags = []cli.Flag{util.JSpathFlag, util.ExecFlag, util.PreloadJSFlag}

	consoleCommand = cli.Command{
		Name:      "console",
		Usage:     "Start an interactive JavaScript environment",
		Action:    localConsole,
		ArgsUsage: " ",
	}

)

func localConsole(ctx *cli.Context) error {
	config := console.Config{
		//Name:    "ucoin console",
		//DataDir: util.MakeDataDir(ctx),
		//DocRoot: ctx.GlobalString(util.JSpathFlag.Name),
		//Client:  nil,
		//Preload: util.MakeConsolePreloads(ctx),
	}

	nodeConfig := makeNodeConfig(ctx)
	console, _ := console.New(config)
	console.Welcome()
	console.Interactive()

	return nil
}