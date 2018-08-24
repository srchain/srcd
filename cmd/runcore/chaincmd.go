package main

import (
	"os"
	"encoding/json"

	"srcd/cmd/utils"
	"srcd/core/blockchain"
	"srcd/node"

	"gopkg.in/urfave/cli.v1"
)

var (
	initCommand = cli.Command{
		Action:    utils.MigrateFlags(initGenesis),
		Name:      "init",
		Usage:     "Bootstrap and initialize a new genesis block",
		ArgsUsage: "<genesisPath>",
		Flags: []cli.Flag{
			utils.DataDirFlag,
			// utils.LightModeFlag,
		},
		Category: "BLOCKCHAIN COMMANDS",
		Description: `
The init command initializes a new genesis block and definition for the network.
This is a destructive action and changes the network in which you will be
participating.

It expects the genesis file as argument.`,
	}
}

// initGenesis will initialise the given JSON format genesis file and writes it as
// the zero'd block or will fail hard if it can't succeed.
func initGenesis(ctx *cli.Context) error {
	// Make sure we have a valid genesis JSON
	genesisPath := ctx.Args().First()
	if len(genesisPath) == 0 {
		utils.Fatalf("Must supply path to genesis JSON file")
	}
	file, err := os.Open(genesisPath)
	if err != nil {
		utils.Fatalf("Failed to read genesis file: %v", err)
	}
	defer file.Close()

	genesis := new(blockchain.Genesis)
	if err := json.NewDecoder(file).Decode(genesis); err != nil {
		utils.Fatalf("invalid genesis file: %v", err)
	}

	// Open an initialise databases
	node := makeNode(ctx)

	chaindb, err := peer.OpenDatabase("chaindata", 0, 0)
	if err != nil {
		utils.Fatalf("Failed to open database: %v", err)
	}
	_, hash, err := core.SetupGenesisBlock(chaindb, genesis)
	if err != nil {
		utils.Fatalf("Failed to write genesis block: %v", err)
	}
	// log.Info("Successfully wrote genesis state", "database", name, "hash", hash)

	return nil
}
