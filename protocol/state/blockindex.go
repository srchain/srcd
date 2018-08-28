package state

import (
	"sync"
	"math/big"
	"srcd/protocol/bc"
)

//BlockIndex is the struct for help chain trace block chain as tree
type BlockIndex struct {
	sync.RWMutex
	index     map[bc.Hash]*BlockNode
	mainChain []*BlockNode
}

// BlockNode represents a block within the block chain and is primarily used to
// aid in selecting the best chain to be the main chain.
type BlockNode struct {
	Parent  *BlockNode // parent is the parent block for this node.
	Hash    bc.Hash    // hash of the block.
	Seed    *bc.Hash   // seed hash of the block
	WorkSum *big.Int   // total amount of work in the chain up to

	Version                uint64
	Height                 uint64
	Timestamp              uint64
	Nonce                  uint64
	Bits                   uint64
	TransactionsMerkleRoot bc.Hash
	TransactionStatusHash  bc.Hash
}
