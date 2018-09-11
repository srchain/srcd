package core

import (
	"srcd/common/common"
	"srcd/core/types"
)

// NewTxsEvent is posted when a batch of transactions enter the transaction pool.
type NewTxsEvent struct{ Txs []*types.Transaction }

// PendingStateEvent is posted pre mining and notifies of pending state changes.
type PendingStateEvent struct{}

// NewMinedBlockEvent is posted when a block has been imported.
type NewMinedBlockEvent struct{ Block *types.Block }

type ChainEvent struct {
	Block *types.Block
	Hash  common.Hash
}

type ChainHeadEvent struct{ Block *types.Block }
