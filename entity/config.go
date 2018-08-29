package node

import (

)

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