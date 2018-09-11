package mempool

import (
	"srcd/core/blockchain"
	"srcd/consensus"
)

// StateProcessor is a basic Processor, which takes care of transitioning
// state from one point to another.
type StateProcessor struct {
	// config *params.ChainConfig // Chain configuration options
	bc     *blockchain.BlockChain         // Canonical block chain
	engine consensus.Engine    // Consensus engine used for block rewards
}

// NewStateProcessor initialises a new StateProcessor.
func NewStateProcessor(bc *blockchain.BlockChain, engine consensus.Engine) *StateProcessor {
	return &StateProcessor{
		// config: config,
		bc:     bc,
		engine: engine,
	}
}
