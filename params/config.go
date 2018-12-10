package params

import (
	"fmt"
	"github.com/srchain/srcd/common/common"
	"math/big"
	"strings"
)


var (
	MainnetGenesisHash = common.HexToHash("0xd4e56740f876aef8c010b86a40d5f56745a118d0906a34e69aec8c0db1cb8fa3")
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

var (
	TestChainConfig = &ChainConfig{ChainID:big.NewInt(9527),Pow:new(PowConfig)}

	// AllEthashProtocolChanges contains every protocol change (EIPs) introduced
	// and accepted by the Ethereum core developers into the Ethash consensus.
	//
	// This configuration is intentionally not using keyed fields to force anyone
	// adding flags to the config to also have to set these fields.
	AllPowProtocolChanges = &ChainConfig{
		ChainID:big.NewInt(10086),
		Pow: new(PowConfig),
	}
)


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


func (err *ConfigCompatError) Error() string {
	return fmt.Sprintf("mismatching %s in database (have %d, want %d, rewindto %d)", err.What, err.StoredConfig, err.NewConfig, err.RewindTo)
}


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



// CheckCompatible checks whether scheduled fork transitions have been imported
// with a mismatching chain configuration.
func (c *ChainConfig) CheckCompatible(newcfg *ChainConfig, height uint64) *ConfigCompatError {
	bhead := new(big.Int).SetUint64(height)

	// Iterate checkCompatible to find the lowest conflict.
	var lasterr *ConfigCompatError
	for {
		err := c.checkCompatible(newcfg, bhead)
		if err == nil || (lasterr != nil && err.RewindTo == lasterr.RewindTo) {
			break
		}
		lasterr = err
		bhead.SetUint64(err.RewindTo)
	}
	return lasterr
}

func (c *ChainConfig) checkCompatible(config *ChainConfig, bhead *big.Int) *ConfigCompatError {

	return nil
}

