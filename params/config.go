package params

import (
	"math/big"
)

var (
	// MainnetChainConfig is the chain parameters to run a node on the main network.
	MainnetChainConfig = &ChainConfig{
		ChainID:        big.NewInt(1),
		Pow:		new(PowConfig),
	}
)

// ChainConfig is the core config which determines the blockchain settings.
type ChainConfig struct {
	// chainId identifies the current chain and is used for replay protection
	ChainID *big.Int `json:"chainId"`

	// Various consensus engines
	Pow *PowConfig
}

// PowConfig is the consensus engine configs for proof-of-work based sealing.
type PowConfig struct{}

// String implements the stringer interface, returning the consensus engine details.
func (c *PowConfig) String() string {
	return "pow"
}

// ConfigCompatError is raised if the locally-stored blockchain is initialised with a
// ChainConfig that would alter the past.
type ConfigCompatError struct {
	What string
	// block numbers of the stored and new configurations
	StoredConfig, NewConfig *big.Int
	// the block number to which the local chain must be rewound to correct the error
	RewindTo uint64
}
