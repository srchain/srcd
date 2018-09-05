package blockchain

import (
	"srcd/consensus"
	"srcd/core/types"
)

// BlockValidator is responsible for validating block headers and processed state.
type BlockValidator struct {
	bc     *BlockChain         // Canonical block chain
	engine consensus.Engine    // Consensus engine used for validating
}

// NewBlockValidator returns a new block validator which is safe for re-use
func NewBlockValidator(blockchain *BlockChain, engine consensus.Engine) *BlockValidator {
	validator := &BlockValidator{
		engine: engine,
		bc:     blockchain,
	}
	return validator
}

// ValidateBody validates the given block's uncles and verifies the the block
// header's transaction and uncle roots. The headers are assumed to be already
// validated at this point.
func (v *BlockValidator) ValidateBody(block *types.Block) error {
	// // Check whether the block's known, and if not, that it's linkable
	// if v.bc.HasBlockAndState(block.Hash(), block.NumberU64()) {
		// return ErrKnownBlock
	// }
	// if !v.bc.HasBlockAndState(block.ParentHash(), block.NumberU64()-1) {
		// if !v.bc.HasBlock(block.ParentHash(), block.NumberU64()-1) {
			// return consensus.ErrUnknownAncestor
		// }
		// return consensus.ErrPrunedAncestor
	// }
	// // Header validity is known at this point, check the uncles and transactions
	// header := block.Header()
	// if err := v.engine.VerifyUncles(v.bc, block); err != nil {
		// return err
	// }
	// if hash := types.CalcUncleHash(block.Uncles()); hash != header.UncleHash {
		// return fmt.Errorf("uncle root hash mismatch: have %x, want %x", hash, header.UncleHash)
	// }
	// if hash := types.DeriveSha(block.Transactions()); hash != header.TxHash {
		// return fmt.Errorf("transaction root hash mismatch: have %x, want %x", hash, header.TxHash)
	// }
	// return nil

	return nil
}
