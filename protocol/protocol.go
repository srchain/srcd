package protocol

import (
	"srcd/protocol/state"
	"sync"
)

//provide functions for working with the Src blockchain.
type Chain struct {
	index *state.BlockIndex
	//orphanManage   *OrphanManage
	txPool         *TxPool
	store          Store
	processBlockCh chan *processBlockMsg

	cond     sync.Cond
	bestNode *state.BlockNode
}