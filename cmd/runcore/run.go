package main

import (
	"github.com/srchain/srcd/cmd/utils"
	"github.com/srchain/srcd/node"
	"github.com/srchain/srcd/server"

	"gopkg.in/urfave/cli.v1"
)

// It creates a default node based on the command line arguments and runs it in
// blocking mode, waiting for it to be shut down.
func entry(ctx *cli.Context) error {
	node := makeNode(ctx)
	startNode(ctx, node)
	node.Wait()

	return nil
}

// startPeer boots up the system peer and all registered protocols, after which
// it unlocks any requested accounts, and starts the RPC/IPC interfaces and the
// miner.
func startNode(ctx *cli.Context, stack *node.Node) {
	// Start up the node itself
	utils.StartNode(stack)

	// Start auxiliary services if enabled
	if ctx.GlobalBool(utils.MiningEnabledFlag.Name) {
		var server *server.Server

		if err := stack.Service(&server); err != nil {
			utils.Fatalf("Srcd service not running: %v", err)
		}

		// Set miner thread from the CLI and start mining
		threads := ctx.GlobalInt(utils.MinerThreadsFlag.Name)
		if err := server.StartMining(threads); err != nil {
			utils.Fatalf("Failed to start mining: %v", err)
		}
	}
}
