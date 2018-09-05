package utils

import (
	"os"
	"path/filepath"
	"runtime"

	"srcd/params"
	"srcd/node"
	"srcd/server"

	"gopkg.in/urfave/cli.v1"
)

var (
	CommandHelpTemplate = `{{.cmd.Name}}{{if .cmd.Subcommands}} command{{end}}{{if .cmd.Flags}} [command options]{{end}} [arguments...]
{{if .cmd.Description}}{{.cmd.Description}}
{{end}}{{if .cmd.Subcommands}}
SUBCOMMANDS:
	{{range .cmd.Subcommands}}{{.Name}}{{with .ShortName}}, {{.}}{{end}}{{ "\t" }}{{.Usage}}
	{{end}}{{end}}{{if .categorizedFlags}}
{{range $idx, $categorized := .categorizedFlags}}{{$categorized.Name}} OPTIONS:
{{range $categorized.Flags}}{{"\t"}}{{.}}
{{end}}
{{end}}{{end}}`
)

func init() {
	cli.AppHelpTemplate = `{{.Name}} {{if .Flags}}[global options] {{end}}command{{if .Flags}} [command options]{{end}} [arguments...]

VERSION:
   {{.Version}}

COMMANDS:
   {{range .Commands}}{{.Name}}{{with .ShortName}}, {{.}}{{end}}{{ "\t" }}{{.Usage}}
   {{end}}{{if .Flags}}
GLOBAL OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}{{end}}
`

	cli.CommandHelpTemplate = CommandHelpTemplate
}

// NewApp creates an app with sane defaults.
func NewApp(gitCommit, usage string) *cli.App {
	app := cli.NewApp()
	app.Name = filepath.Base(os.Args[0])
	app.Author = ""
	app.Email = ""
	app.Version = params.Version
	if len(gitCommit) >= 8 {
		app.Version += "-" + gitCommit[:8]
	}
	app.Usage = usage
	return app
}

// These are all the command line flags we support.
var (
	// General settings
	DataDirFlag = DirectoryFlag{
		Name:  "datadir",
		Usage: "Data directory for the databases and keystore",
		Value: DirectoryString{node.DefaultDataDir()},
	}
	KeyStoreDirFlag = DirectoryFlag{
		Name:  "keystore",
		Usage: "Directory for the keystore (default = inside the datadir)",
	}
	IdentityFlag = cli.StringFlag{
		Name:  "identity",
		Usage: "Custom node name",
	}
	TestnetFlag = cli.BoolFlag{
		Name:  "testnet",
		Usage: "Ropsten network: pre-configured proof-of-work test network",
	}

	MinerThreadsFlag = cli.IntFlag{
		Name:  "minerthreads",
		Usage: "Number of CPU threads to use for mining",
		Value: runtime.NumCPU(),
	}

	CoinbaseFlag = cli.StringFlag{
		Name:  "coinbase",
		Usage: "Public address for block mining rewards (default = first account created)",
		Value: "0",
	}
)

// setCoinbase retrieves the etherbase either from the directly specified
// command line flags or from the keystore if CLI indexed.
func setCoinbase(ctx *cli.Context, wallet *Wallet, cfg *server.Config) {
	if ctx.GlobalIsSet(CoinbaseFlag.Name) {
		// account, err := MakeAddress(ks, ctx.GlobalString(EtherbaseFlag.Name))
		// if err != nil {
			// Fatalf("Option %q: %v", EtherbaseFlag.Name, err)
		// }
		cfg.Coinbase = wallet.GetAddress()
	}
}

// SetNodeConfig applies peer-related command line flags to the config.
func SetNodeConfig(ctx *cli.context, cfg *node.config) {
	switch {
	case ctx.GlobalIsSet(DataDirFlag.Name):
		cfg.DataDir = ctx.GlobalString(DataDirFlag.Name)
	case ctx.GlobalBool(TestnetFlag.Name):
		cfg.DataDir = filepath.Join(node.DefaultDataDir(), "testnet")
	}

	if ctx.GlobalIsSet(KeyStoreDirFlag.Name) {
		cfg.KeyStoreDir = ctx.GlobalString(KeyStoreDirFlag.Name)
	}
}

// SetServerConfig applies server-related command line flags to the config.
func SetServerConfig(ctx *cli.Context, node *node.Node, cfg *server.Config) {
	// ks := node.AccountManager().Backends(keystore.KeyStoreType)[0].(*keystore.KeyStore)
	wallet := node.wallet
	setCoinbase(ctx, wallet, cfg)

	if ctx.GlobalIsSet(MinerThreadsFlag.Name) {
		cfg.MinerThreads = ctx.GlobalInt(MinerThreadsFlag.Name)
	}
}

// RegisterService adds an srcd client to the node.
func RegisterService(node *node.Node, cfg *server.Config) {
	err := node.Register(func(ctx *node.ServiceContext) (node.Service, error) {
		fullNode, err := server.New(ctx, cfg)

		return fullNode, err
	})

	if err != nil {
		Fatalf("Failed to register the srcd service: %v", err)
	}
}

// MigrateFlags sets the global flag from a local flag when it's set.
// This is a temporary function used for migrating old command/flags to the
// new format.
func MigrateFlags(action func(ctx *cli.Context) error) func(*cli.Context) error {
	return func(ctx *cli.Context) error {
		for _, name := range ctx.FlagNames() {
			if ctx.IsSet(name) {
				ctx.GlobalSet(name, ctx.String(name))
			}
		}
		return action(ctx)
	}
}
