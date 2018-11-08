package mempool

import (
	"sync"
	"time"

	"github.com/srchain/srcd/common/common"
	"github.com/srchain/srcd/core/transaction"
	"github.com/srchain/srcd/core/types"
)

// blockChain provides the state of blockchain and current gas limit to do
// some pre checks in tx pool and event subscribers.
type blockChain interface {
	CurrentBlock() *types.Block
	GetBlock(hash common.Hash, number uint64) *types.Block

	// SubscribeChainHeadEvent(ch chan<- ChainHeadEvent) event.Subscription
}

// TxPoolConfig are the configuration parameters of the transaction pool.
type TxPoolConfig struct {
	NoLocals  bool          // Whether local transaction handling should be disabled
	Journal   string        // Journal of local transactions to survive node restarts
	Rejournal time.Duration // Time interval to regenerate the local transaction journal

	PriceLimit uint64 // Minimum price to enforce for acceptance into the pool
	// PriceBump  uint64 // Minimum price bump percentage to replace an already existing transaction (nonce)

	AccountSlots uint64 // Number of executable transaction slots guaranteed per account
	GlobalSlots  uint64 // Maximum number of executable transaction slots for all accounts
	AccountQueue uint64 // Maximum number of non-executable transaction slots permitted per account
	GlobalQueue  uint64 // Maximum number of non-executable transaction slots for all accounts

	Lifetime time.Duration // Maximum amount of time non-executable transaction are queued
}

// TxPool contains all currently known transactions. Transactions
// enter the pool when they are received from the network or submitted
// locally. They exit the pool when they are included in the blockchain.
//
// The pool separates processable transactions (which can be applied to the
// current state) and future transactions. Transactions move between those
// two states over time as they are received and processed.
type TxPool struct {
	config  TxPoolConfig
	chain   blockChain
	pool    *transaction.TxPool
	pending types.Transactions

	mu sync.RWMutex
	wg sync.WaitGroup // for shutdown sync
}

// NewTxPool creates a new transaction pool to gather, sort and filter inbound
// transactions from the network.
func NewTxPool(config TxPoolConfig, chain blockChain) *TxPool {
	// Create the transaction pool with its initial settings
	pool := &TxPool{
		pool:    transaction.NewTxPool(),
		config:  config,
		chain:   chain,
		pending: make(types.Transactions, 1024),
	}

	// Start the event loop and return
	pool.wg.Add(1)
	go pool.loop()

	return pool
}

// loop is the transaction pool's main event loop, waiting for and reacting to
// outside blockchain events as well as for various reporting and transaction
// eviction events.
func (pool *TxPool) loop() {
	defer pool.wg.Done()

	txCh := pool.pool.GetMsgCh()

	for {
		select {
		case ev := <-txCh:
			tx := &types.Transaction{
				Tx: ev.Tx.TxData,
			}
			pool.enqueueTx(tx)
		}
	}
}

func (pool *TxPool) enqueueTx(tx *types.Transaction) {
	pool.pending = append(pool.pending, tx)
}

func (pool *TxPool) Pending() (types.Transactions, error) {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	return pool.pending, nil
}
