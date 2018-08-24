package main

import (
	// "shilling/cmd/utils"
	// "shilling/peer"

	"gopkg.in/urfave/cli.v1"
)

// It creates a default peer based on the command line arguments and runs it in
// blocking mode, waiting for it to be shut down.
func entry(ctx *cli.Context) error {
	// peer := makePeer(ctx)
	// startPeer(ctx, peer)
	// peer.Wait()

	return nil
}

// startPeer boots up the system peer and all registered protocols, after which
// it unlocks any requested accounts, and starts the RPC/IPC interfaces and the
// miner.
// func startPeer(ctx *cli.Context, peer *peer.Peer) {
	// // Start up the node itself
	// utils.StartPeer(peer)

	// // wallet op ...

	// // Start auxiliary services if enabled
	// if ctx.GlobalBool(utils.MiningEnabledFlag.Name) || ctx.GlobalBool(utils.DeveloperFlag.Name) {
		// var node node.Node

		// // Start mining
		// if err := node.StartMining(true); err != nil {
			// utils.Fatalf("Failed to start mining: %v", err)
		// }
	// }
// }
