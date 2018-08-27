package params

import (
	"math/big"
)

// ChainConfig is the core config which determines the blockchain settings.
type ChainConfig struct {
	// chainId identifies the current chain and is used for replay protection
	ChainID *big.Int `json:"chainId"`

	// Various consensus engines
	hash *PowConfig `json:"ethash,omitempty"`
}

// PowConfig is the consensus engine configs for proof-of-work based sealing.
type PowConfig struct{}
