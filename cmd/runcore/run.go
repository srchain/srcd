package main

import (
	"fmt"
	"srcd/cmd/utils"
	"srcd/node"
	"srcd/server"

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

	// wallet op ...

	// Start auxiliary services if enabled
	if ctx.GlobalBool(utils.MiningEnabledFlag.Name) {
		fmt.Println("ssss")
		var server *server.Server

		if err := stack.Service(&server); err != nil {
			utils.Fatalf("Srcd service not running: %v", err)
		}

		// Start mining
		if err := server.StartMining(true); err != nil {
			utils.Fatalf("Failed to start mining: %v", err)
		}
	}
}
