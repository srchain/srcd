package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"unicode"
	"path/filepath"

	"srcd/cmd/utils"
	"srcd/node"
	"srcd/entity"

	"github.com/naoina/toml"

	"gopkg.in/urfave/cli.v1"
)

var (
	configFileFlag = cli.StringFlag{
		Name:  "config",
		Usage: "TOML configuration file",
	}
)

type config struct {
	Server    server.Config
	Node      node.Config
}

func loadConfig(file string, cfg *config) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	err = tomlSettings.NewDecoder(bufio.NewReader(f)).Decode(cfg)
	// Add file name to errors that have a line number.
	if _, ok := err.(*toml.LineError); ok {
		err = errors.New(file + ", " + err.Error())
	}
	return err
}

func defaultNodeConfig() node.Config {
	cfg := node.DefaultConfig
	cfg.Name = "srcd"
	cfg.Version = params.Version

	return cfg
}

func makeConfig(ctx *cli.Context) *config {
	// Default config.
	cfg := config{
		Server:    server.DefaultConfig,
		Node:      defaultNodeConfig,
	}

	// Load config file.
	if file := ctx.GlobalString(configFileFlag.Name); file != "" {
		if err := loadConfig(file, &cfg); err != nil {
			utils.Fatalf("%v", err)
		}
	}

	// Apply flags.
	// utils.SetServerConfig(ctx, &cfg.Server)
	utils.SetNodeConfig(ctx, &cfg.Node)

	return &cfg
}

func makeNode(ctx *cli.Context) *node.Node {
	cfg := makeConfig(ctx)

	node := node.New(&cfg.Node)

	return node
}
