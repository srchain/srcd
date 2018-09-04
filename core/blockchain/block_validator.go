package blockchain

import (

)

// BlockValidator is responsible for validating block headers and processed state.
type BlockValidator struct {
	// config *params.ChainConfig // Chain configuration options
	bc     *BlockChain         // Canonical block chain
	engine consensus.Engine    // Consensus engine used for validating
}

// NewBlockValidator returns a new block validator which is safe for re-use
func NewBlockValidator(blockchain *BlockChain, engine consensus.Engine) *BlockValidator {
	validator := &BlockValidator{
		// config: config,
		engine: engine,
		bc:     blockchain,
	}
	return validator
}
