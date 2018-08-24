package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"sort"

	"srcd/cmd/utils"
	"srcd/console"

	"gopkg.in/urfave/cli.v1"
)

// The app that holds all commands and flags.
var app = utils.NewApp("", "the SilkRoad command line interface")

func init() {
	// Initialize the CLI app
	app.Action = entry
	app.HideVersion = true
	app.Copyright = "Copyright 2018 The Shilling Authors"
	app.Commands = []cli.Command{
		initCommand,
	}
	sort.Sort(cli.CommandsByName(app.Commands))

	// app.Flags = append(app.Flags, nodeFlags...)
	// app.Flags = append(app.Flags, rpcFlags...)

	app.Before = func(ctx *cli.Context) error {
		// Use all processor cores.
		runtime.GOMAXPROCS(runtime.NumCPU())

		// Block and transaction processing can cause bursty allocations.  This
		// limits the garbage collector from excessively overallocating during
		// bursts.  This value was arrived at with the help of profiling live
		// usage.
		debug.SetGCPercent(20)

		return nil
	}

	app.After = func(ctx *cli.Context) error {
		// Resets terminal mode.
		console.Stdin.Close()
		return nil
	}
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
