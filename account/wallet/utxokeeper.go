package wallet

import (
	"sync"
	"srcd/protocol/transaction"
)

type utxoKeeper struct {
	// `sync/atomic` expects the first word in an allocated struct to be 64-bit
	// aligned on both ARM and x86-32. See https://goo.gl/zW7dgq for more details.
	mtx           sync.RWMutex
	unconfirmed  map[transaction.Hash]*Utxo
}

func (uk *utxoKeeper) AddUnconfirmedTx(utxos []*Utxo)  {
	uk.mtx.Lock()
	defer uk.mtx.Unlock()

	for _, utxo := range utxos {
		uk.unconfirmed[utxo.OutputID] = utxo
	}
}
