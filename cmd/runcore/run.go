package main

import (
	"gopkg.in/urfave/cli.v1"
)

// It creates a default peer based on the command line arguments and runs it in
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
func startNode(ctx *cli.Context, node *node.Node) {
	// Start up the node itself
	utils.StartNode(node)

	// wallet op ...

	// Start auxiliary services if enabled
	if ctx.GlobalBool(utils.MiningEnabledFlag.Name) || ctx.GlobalBool(utils.DeveloperFlag.Name) {
		var node node.Node

		// Start mining
		if err := node.StartMining(true); err != nil {
			utils.Fatalf("Failed to start mining: %v", err)
		}
	}
}
