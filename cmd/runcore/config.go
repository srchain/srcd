package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"reflect"
	"unicode"

	"github.com/srchain/srcd/cmd/utils"
	"github.com/srchain/srcd/node"
	"github.com/srchain/srcd/params"
	"github.com/srchain/srcd/server"
	"github.com/naoina/toml"

	"gopkg.in/urfave/cli.v1"
)

var (
	configFileFlag = cli.StringFlag{
		Name:  "config",
		Usage: "TOML configuration file",
	}
)

// These settings ensure that TOML keys use the same names as Go struct fields.
var tomlSettings = toml.Config{
	NormFieldName: func(rt reflect.Type, key string) string {
		return key
	},
	FieldToKey: func(rt reflect.Type, field string) string {
		return field
	},
	MissingField: func(rt reflect.Type, field string) error {
		link := ""
		if unicode.IsUpper(rune(rt.Name()[0])) && rt.PkgPath() != "main" {
			link = fmt.Sprintf(", see https://godoc.org/%s#%s for available fields", rt.PkgPath(), rt.Name())
		}
		return fmt.Errorf("field '%s' is not defined in %s%s", field, rt.String(), link)
	},
}

type config struct {
	Server server.Config
	Node   node.Config
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

func makeConfigNode(ctx *cli.Context) (*node.Node, config) {
	// Default config.
	cfg := config{
		Server: server.DefaultConfig,
		Node:   defaultNodeConfig(),
	}

	// Load config file.
	if file := ctx.GlobalString(configFileFlag.Name); file != "" {
		if err := loadConfig(file, &cfg); err != nil {
			utils.Fatalf("%v", err)
		}
	}

	// Apply flags.
	utils.SetNodeConfig(ctx, &cfg.Node)

	node, err := node.New(&cfg.Node)
	if err != nil {
		utils.Fatalf("Failed to create the protocol node: %v", err)
	}

	utils.SetServerConfig(ctx, node, &cfg.Server)

	return node, cfg
}

func makeNode(ctx *cli.Context) *node.Node {
	node, cfg := makeConfigNode(ctx)

	utils.RegisterService(node, &cfg.Server)

	return node
}
