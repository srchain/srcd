package transaction

import (
	"github.com/srchain/srcd/common/common"
	"github.com/srchain/srcd/core"
	"github.com/srchain/srcd/core/types"
	"github.com/srchain/srcd/event"
	"sync"
	"time"

	"github.com/srchain/srcd/errors"
	"github.com/srchain/srcd/log"
)

const (
	MsgNewTx = iota
)

type TxPool struct {
	Utxo   map[Hash]Tx
	Tx     Tx
	Weight uint64
	Height uint64
	Fee    uint64
	Mtx    sync.RWMutex
	MsgCh  chan *TxPoolMsg
	Pool   map[Hash]*TxPoolMsg
}

//TODO: 有请张先生现身说法
func (tp *TxPool) AddRemotes([]*types.Transaction) []error {
	panic("implement me")
}

func (tp *TxPool) Pending() (map[common.Address]types.Transactions, error) {
	panic("implement me")
}

func (tp *TxPool) SubscribeNewTxsEvent(chan<- core.NewTxsEvent) event.Subscription {
	panic("implement me")
}

type TxPoolMsg struct {
	Tx      Tx
	Added   time.Time
	Weight  uint64
	Fee     uint64
	MsgType int
}

func NewTxPool() *TxPool {
	return &TxPool{
		Utxo:   make(map[Hash]Tx),
		Tx:     Tx{},
		Weight: uint64(0),
		Height: uint64(0),
		Fee:    uint64(0),
		MsgCh:  make(chan *TxPoolMsg, 1000), //chan to broadcast tx msg,
		Pool:   make(map[Hash]*TxPoolMsg),
	}
}

func (tp *TxPool) GetMsgCh() <-chan *TxPoolMsg {
	return tp.MsgCh
}

func (tp *TxPool) AddTransaction(tx Tx, fee uint64) error {
	tp.Mtx.Lock()
	defer tp.Mtx.Unlock()

	msg := &TxPoolMsg{tx, time.Now(), tx.SerializedSize, fee, MsgNewTx}
	for _, id := range tx.ResultIds {
		tp.Utxo[*id] = tx
	}
	tp.Pool[tx.ID] = msg
	tp.MsgCh <- msg
	log.Info("add txpool ", "tx_id", tx.ID.String())
	return nil
}

func (tp *TxPool) GetTransaction(hash *Hash) (*TxPoolMsg, error) {
	tp.Mtx.RLock()
	defer tp.Mtx.RUnlock()

	msg := tp.Pool[*hash]
	if msg != nil {
		return msg, nil
	}
	return &TxPoolMsg{}, errors.New("txpool has no this tx")
}
