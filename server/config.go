package server

import (

)

// DefaultConfig contains default settings for use on main net.
var DefaultConfig = Config{
	// SyncMode: downloader.FastSync,
	Pow: pow.Config{
		CacheDir:       "pow",
		CachesInMem:    2,
		CachesOnDisk:   3,
		DatasetsInMem:  1,
		DatasetsOnDisk: 2,
	},
	// NetworkId:     1,
	// LightPeers:    100,
	DatabaseCache: 768,
	// TrieCache:     256,
	// TrieTimeout:   60 * time.Minute,

	TxPool: core.DefaultTxPoolConfig,
}

func init() {
	home := os.Getenv("HOME")
	if home == "" {
		if user, err := user.Current(); err == nil {
			home = user.HomeDir
		}
	}
	if runtime.GOOS == "windows" {
		DefaultConfig.Pow.DatasetDir = filepath.Join(home, "AppData", "Pow")
	} else {
		DefaultConfig.Pow.DatasetDir = filepath.Join(home, ".pow")
	}
}

type Config struct {
	// The genesis block, which is inserted if the database is empty.
	// If nil, main net block is used.
	Genesis *blockchain.Genesis `toml:",omitempty"`

	// Protocol options
	// NetworkId uint64 // Network ID to use for selecting peers to connect to
	// SyncMode  downloader.SyncMode
	// NoPruning bool

	// Database options
	// SkipBcVersionCheck bool `toml:"-"`
	DatabaseHandles    int  `toml:"-"`
	DatabaseCache      int
	// TrieCache          int
	// TrieTimeout        time.Duration

	// Mining-related options
	Coinbase    common.Address `toml:",omitempty"`
	MinerThreads int            `toml:",omitempty"`
	ExtraData    []byte         `toml:",omitempty"`

	// PoW options
	Pow pow.Config

	// Transaction pool options
	TxPool core.TxPoolConfig

	// Gas Price Oracle options
	// GPO gasprice.Config
}
