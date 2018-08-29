package miner

import (

)

// Agent can register themself with the worker
type Agent interface {
	Work() chan<- *Work
	SetReturnCh(chan<- *Result)
	Stop()
	Start()
	GetHashRate() int64
}

// Work is the workers current environment and holds
// all of the current state information
type Work struct {
	config *params.ChainConfig
	signer types.Signer

	state     *state.StateDB // apply state changes here
	ancestors *set.Set       // ancestor set (used for checking uncle parent validity)
	family    *set.Set       // family set (used for checking uncle invalidity)
	// uncles    *set.Set       // uncle set
	tcount    int            // tx count in cycle
	// gasPool   *core.GasPool  // available gas used to pack transactions

	Block *types.Block // the new block

	header   *types.Header
	txs      []*types.Transaction
	receipts []*types.Receipt

	createdAt time.Time
}

type Result struct {
	Work  *Work
	Block *types.Block
}

// worker is the main object which takes care of applying messages to the new state
type worker struct {
	config *params.ChainConfig
	engine consensus.Engine

	mu sync.Mutex

	// update loop
	mux          *event.TypeMux
	txsCh        chan core.NewTxsEvent
	txsSub       event.Subscription
	chainHeadCh  chan core.ChainHeadEvent
	chainHeadSub event.Subscription
	chainSideCh  chan core.ChainSideEvent
	chainSideSub event.Subscription
	wg           sync.WaitGroup

	agents map[Agent]struct{}
	recv   chan *Result

	eth     Backend
	chain   *core.BlockChain
	proc    core.Validator
	chainDb ethdb.Database

	coinbase common.Address
	extra    []byte

	currentMu sync.Mutex
	current   *Work

	// snapshotMu    sync.RWMutex
	// snapshotBlock *types.Block
	// snapshotState *state.StateDB

	// uncleMu        sync.Mutex
	// possibleUncles map[common.Hash]*types.Block

	// unconfirmed *unconfirmedBlocks // set of locally mined blocks pending canonicalness confirmations

	// atomic status counters
	mining int32
	atWork int32
}

func newWorker(config *params.ChainConfig, engine consensus.Engine, coinbase common.Address, eth Backend, mux *event.TypeMux) *worker {
	worker := &worker{
		config:         config,
		engine:         engine,
		eth:            eth,
		mux:            mux,
		txsCh:          make(chan core.NewTxsEvent, txChanSize),
		chainHeadCh:    make(chan core.ChainHeadEvent, chainHeadChanSize),
		chainSideCh:    make(chan core.ChainSideEvent, chainSideChanSize),
		chainDb:        eth.ChainDb(),
		recv:           make(chan *Result, resultQueueSize),
		chain:          eth.BlockChain(),
		proc:           eth.BlockChain().Validator(),
		coinbase:       coinbase,
		agents:         make(map[Agent]struct{}),
		unconfirmed:    newUnconfirmedBlocks(eth.BlockChain(), miningLogAtDepth),
	}
	// Subscribe NewTxsEvent for tx pool
	worker.txsSub = eth.TxPool().SubscribeNewTxsEvent(worker.txsCh)
	// Subscribe events for blockchain
	worker.chainHeadSub = eth.BlockChain().SubscribeChainHeadEvent(worker.chainHeadCh)
	worker.chainSideSub = eth.BlockChain().SubscribeChainSideEvent(worker.chainSideCh)
	go worker.update()

	go worker.wait()
	worker.commitNewWork()

	return worker
}
