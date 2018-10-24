package server

import (
	"srcd/common/common"
	"srcd/core/blockchain"
	"srcd/core/mempool"
)

// DefaultConfig contains default settings for use on main net.
var DefaultConfig = Config{
	// SyncMode: downloader.FastSync,
	// NetworkId:     1,
	// LightPeers:    100,
	DatabaseCache: 768,
	// TrieCache:     256,
	// TrieTimeout:   60 * time.Minute,

	// TxPool: core.DefaultTxPoolConfig,
}

type Config struct {
	// The genesis block, which is inserted if the database is empty.
	// If nil, main net block is used.
	Genesis         *blockchain.Genesis `toml:",omitempty"`

	// Protocol options
	// NetworkId uint64 // Network ID to use for selecting peers to connect to
	// SyncMode  downloader.SyncMode
	// NoPruning bool

	// Database options
	DatabaseHandles int `toml:"-"`
	DatabaseCache   int
	// TrieCache     int
	// TrieTimeout   time.Duration

	// Mining-related options
	Coinbase        common.Address `toml:",omitempty"`
	MinerThreads    int            `toml:",omitempty"`
	ExtraData       []byte         `toml:",omitempty"`

	// Transaction pool options
	TxPool mempool.TxPoolConfig

	// Gas Price Oracle options
	// GPO gasprice.Config
}
