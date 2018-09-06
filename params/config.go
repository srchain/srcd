package params

import (
	"math/big"
	"strings"
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




type NetParams struct {
	// Name defines a human-readable identifier for the network.
	Name            string
	Bech32HRPSegwit string
}

var TestNetParams  = NetParams{
	Name:            "testnet",
	Bech32HRPSegwit: "sr",
}

func IsBech32SegwitPrefix(prefix string, params NetParams) bool {
	prefix = strings.ToLower(prefix)
	return prefix == params.Bech32HRPSegwit+"1"
}
