package protocol

import (
	"sync"
	"srcd/protocol/bc"
	"time"
	"srcd/protocol/bc/types"
)

//TxPool is use for store the unconfirmed transaction
type TxPool struct {
	lastUpdated int64
	mtx         sync.RWMutex
	pool        map[bc.Hash]*TxDesc
	utxo        map[bc.Hash]bc.Hash
	//errCache    *lru.Cache
	newTxCh     chan *types.Tx
}

// TxDesc store tx and related info for mining strategy
type TxDesc struct {
	Tx       *types.Tx
	Added    time.Time
	Height   uint64
	Weight   uint64
	Fee      uint64
	FeePerKB uint64
}