package utils

import (
	"os"
	"path/filepath"

	"srcd/params"

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
		Value: DirectoryString{"~/.gosr/data"},
	}
	IdentityFlag = cli.StringFlag{
		Name:  "identity",
		Usage: "Custom node name",
	}

	// RPC settings
	// RPCEnabledFlag = cli.BoolFlag{
		// Name:  "rpc",
		// Usage: "Enable the HTTP-RPC server",
	// }
	// RPCListenAddrFlag = cli.StringFlag{
		// Name:  "rpcaddr",
		// Usage: "HTTP-RPC server listening interface",
		// Value: "localhost",
	// }
	// RPCPortFlag = cli.IntFlag{
		// Name:  "rpcport",
		// Usage: "HTTP-RPC server listening port",
		// Value: 8545,
	// }
)

// // SetPeerConfig applies peer-related command line flags to the config.
// func SetPeerConfig(ctx *cli.Context, cfg *node.Config) {
	// // SetP2PConfig(ctx, &cfg.P2P)
	// // setIPC(ctx, cfg)
	// // setHTTP(ctx, cfg)
	// // setWS(ctx, cfg)
	// // setNodeUserIdent(ctx, cfg)

	// switch {
	// case ctx.GlobalIsSet(DataDirFlag.Name):
		// cfg.DataDir = ctx.GlobalString(DataDirFlag.Name)
	// case ctx.GlobalBool(DeveloperFlag.Name):
		// cfg.DataDir = "" // unless explicitly requested, use memory databases
	// case ctx.GlobalBool(TestnetFlag.Name):
		// cfg.DataDir = filepath.Join(node.DefaultDataDir(), "testnet")
	// case ctx.GlobalBool(RinkebyFlag.Name):
		// cfg.DataDir = filepath.Join(node.DefaultDataDir(), "rinkeby")
	// }

	// if ctx.GlobalIsSet(KeyStoreDirFlag.Name) {
		// cfg.KeyStoreDir = ctx.GlobalString(KeyStoreDirFlag.Name)
	// }
	// if ctx.GlobalIsSet(LightKDFFlag.Name) {
		// cfg.UseLightweightKDF = ctx.GlobalBool(LightKDFFlag.Name)
	// }
	// if ctx.GlobalIsSet(NoUSBFlag.Name) {
		// cfg.NoUSB = ctx.GlobalBool(NoUSBFlag.Name)
	// }
// }

// // SetNodeConfig applies node-related command line flags to the config.
// func SetEthConfig(ctx *cli.Context, stack *node.Node, cfg *eth.Config) {
	// // setEtherbase(ctx, ks, cfg)
	// // setGPO(ctx, &cfg.GPO)
	// // setTxPool(ctx, &cfg.TxPool)
	// // setEthash(ctx, cfg)

	// switch {
	// case ctx.GlobalIsSet(SyncModeFlag.Name):
		// cfg.SyncMode = *GlobalTextMarshaler(ctx, SyncModeFlag.Name).(*downloader.SyncMode)
	// case ctx.GlobalBool(FastSyncFlag.Name):
		// cfg.SyncMode = downloader.FastSync
	// case ctx.GlobalBool(LightModeFlag.Name):
		// cfg.SyncMode = downloader.LightSync
	// }
	// if ctx.GlobalIsSet(LightServFlag.Name) {
		// cfg.LightServ = ctx.GlobalInt(LightServFlag.Name)
	// }
	// if ctx.GlobalIsSet(LightPeersFlag.Name) {
		// cfg.LightPeers = ctx.GlobalInt(LightPeersFlag.Name)
	// }
	// if ctx.GlobalIsSet(NetworkIdFlag.Name) {
		// cfg.NetworkId = ctx.GlobalUint64(NetworkIdFlag.Name)
	// }

	// if ctx.GlobalIsSet(CacheFlag.Name) || ctx.GlobalIsSet(CacheDatabaseFlag.Name) {
		// cfg.DatabaseCache = ctx.GlobalInt(CacheFlag.Name) * ctx.GlobalInt(CacheDatabaseFlag.Name) / 100
	// }
	// cfg.DatabaseHandles = makeDatabaseHandles()

	// if gcmode := ctx.GlobalString(GCModeFlag.Name); gcmode != "full" && gcmode != "archive" {
		// Fatalf("--%s must be either 'full' or 'archive'", GCModeFlag.Name)
	// }
	// cfg.NoPruning = ctx.GlobalString(GCModeFlag.Name) == "archive"

	// if ctx.GlobalIsSet(CacheFlag.Name) || ctx.GlobalIsSet(CacheGCFlag.Name) {
		// cfg.TrieCache = ctx.GlobalInt(CacheFlag.Name) * ctx.GlobalInt(CacheGCFlag.Name) / 100
	// }
	// if ctx.GlobalIsSet(MinerThreadsFlag.Name) {
		// cfg.MinerThreads = ctx.GlobalInt(MinerThreadsFlag.Name)
	// }
	// if ctx.GlobalIsSet(DocRootFlag.Name) {
		// cfg.DocRoot = ctx.GlobalString(DocRootFlag.Name)
	// }
	// if ctx.GlobalIsSet(ExtraDataFlag.Name) {
		// cfg.ExtraData = []byte(ctx.GlobalString(ExtraDataFlag.Name))
	// }
	// if ctx.GlobalIsSet(GasPriceFlag.Name) {
		// cfg.GasPrice = GlobalBig(ctx, GasPriceFlag.Name)
	// }
	// if ctx.GlobalIsSet(VMEnableDebugFlag.Name) {
		// // TODO(fjl): force-enable this in --dev mode
		// cfg.EnablePreimageRecording = ctx.GlobalBool(VMEnableDebugFlag.Name)
	// }

	// // Override any default configs for hard coded networks.
	// switch {
	// case ctx.GlobalBool(TestnetFlag.Name):
		// if !ctx.GlobalIsSet(NetworkIdFlag.Name) {
			// cfg.NetworkId = 3
		// }
		// cfg.Genesis = core.DefaultTestnetGenesisBlock()
	// case ctx.GlobalBool(RinkebyFlag.Name):
		// if !ctx.GlobalIsSet(NetworkIdFlag.Name) {
			// cfg.NetworkId = 4
		// }
		// cfg.Genesis = core.DefaultRinkebyGenesisBlock()
	// case ctx.GlobalBool(DeveloperFlag.Name):
		// if !ctx.GlobalIsSet(NetworkIdFlag.Name) {
			// cfg.NetworkId = 1337
		// }
		// // Create new developer account or reuse existing one
		// var (
			// developer accounts.Account
			// err       error
		// )
		// if accs := ks.Accounts(); len(accs) > 0 {
			// developer = ks.Accounts()[0]
		// } else {
			// developer, err = ks.NewAccount("")
			// if err != nil {
				// Fatalf("Failed to create developer account: %v", err)
			// }
		// }
		// if err := ks.Unlock(developer, ""); err != nil {
			// Fatalf("Failed to unlock developer account: %v", err)
		// }
		// log.Info("Using developer account", "address", developer.Address)

		// cfg.Genesis = core.DeveloperGenesisBlock(uint64(ctx.GlobalInt(DeveloperPeriodFlag.Name)), developer.Address)
		// if !ctx.GlobalIsSet(GasPriceFlag.Name) {
			// cfg.GasPrice = big.NewInt(1)
		// }
	// }
	// // TODO(fjl): move trie cache generations into config
	// if gen := ctx.GlobalInt(TrieCacheGenFlag.Name); gen > 0 {
		// state.MaxTrieCacheGen = uint16(gen)
	// }
// }

// // SetDashboardConfig applies dashboard related command line flags to the config.
// func SetDashboardConfig(ctx *cli.Context, cfg *dashboard.Config) {
	// cfg.Host = ctx.GlobalString(DashboardAddrFlag.Name)
	// cfg.Port = ctx.GlobalInt(DashboardPortFlag.Name)
	// cfg.Refresh = ctx.GlobalDuration(DashboardRefreshFlag.Name)
// }

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
