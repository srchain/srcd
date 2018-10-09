package mempool

import (
)

// // TxPoolConfig are the configuration parameters of the transaction pool.
// type TxPoolConfig struct {
	// NoLocals  bool          // Whether local transaction handling should be disabled
	// Journal   string        // Journal of local transactions to survive node restarts
	// Rejournal time.Duration // Time interval to regenerate the local transaction journal

	// PriceLimit uint64 // Minimum gas price to enforce for acceptance into the pool
	// PriceBump  uint64 // Minimum price bump percentage to replace an already existing transaction (nonce)

	// AccountSlots uint64 // Minimum number of executable transaction slots guaranteed per account
	// GlobalSlots  uint64 // Maximum number of executable transaction slots for all accounts
	// AccountQueue uint64 // Maximum number of non-executable transaction slots permitted per account
	// GlobalQueue  uint64 // Maximum number of non-executable transaction slots for all accounts

	// Lifetime time.Duration // Maximum amount of time non-executable transaction are queued
// }

// // TxPool contains all currently known transactions. Transactions
// // enter the pool when they are received from the network or submitted
// // locally. They exit the pool when they are included in the blockchain.
// //
// // The pool separates processable transactions (which can be applied to the
// // current state) and future transactions. Transactions move between those
// // two states over time as they are received and processed.
// type TxPool struct {
	// config       TxPoolConfig
	// chainconfig  *params.ChainConfig
	// chain        blockChain
	// txFeed       event.Feed
	// scope        event.SubscriptionScope
	// chainHeadCh  chan ChainHeadEvent
	// chainHeadSub event.Subscription
	// signer       types.Signer
	// mu           sync.RWMutex

	// currentState  *state.StateDB      // Current state in the blockchain head
	// pendingState  *state.ManagedState // Pending state tracking virtual nonces
	// currentMaxGas uint64              // Current gas limit for transaction caps

	// locals  *accountSet // Set of local transaction to exempt from eviction rules
	// journal *txJournal  // Journal of local transaction to back up to disk

	// pending map[common.Address]*txList   // All currently processable transactions
	// queue   map[common.Address]*txList   // Queued but non-processable transactions
	// beats   map[common.Address]time.Time // Last heartbeat from each known account
	// all     *txLookup                    // All transactions to allow lookups
	// priced  *txPricedList                // All transactions sorted by price

	// wg sync.WaitGroup // for shutdown sync

	// homestead bool
// }

// // NewTxPool creates a new transaction pool to gather, sort and filter inbound
// // transactions from the network.
// func NewTxPool(config TxPoolConfig, chainconfig *params.ChainConfig, chain blockChain) *TxPool {
	// // Sanitize the input to ensure no vulnerable gas prices are set
	// config = (&config).sanitize()

	// // Create the transaction pool with its initial settings
	// pool := &TxPool{
		// config:      config,
		// chainconfig: chainconfig,
		// chain:       chain,
		// signer:      types.NewEIP155Signer(chainconfig.ChainID),
		// pending:     make(map[common.Address]*txList),
		// queue:       make(map[common.Address]*txList),
		// beats:       make(map[common.Address]time.Time),
		// all:         newTxLookup(),
		// chainHeadCh: make(chan ChainHeadEvent, chainHeadChanSize),
		// gasPrice:    new(big.Int).SetUint64(config.PriceLimit),
	// }
	// pool.locals = newAccountSet(pool.signer)
	// pool.priced = newTxPricedList(pool.all)
	// pool.reset(nil, chain.CurrentBlock().Header())

	// // If local transactions and journaling is enabled, load from disk
	// if !config.NoLocals && config.Journal != "" {
		// pool.journal = newTxJournal(config.Journal)

		// if err := pool.journal.load(pool.AddLocals); err != nil {
			// log.Warn("Failed to load transaction journal", "err", err)
		// }
		// if err := pool.journal.rotate(pool.local()); err != nil {
			// log.Warn("Failed to rotate transaction journal", "err", err)
		// }
	// }
	// // Subscribe events from blockchain
	// pool.chainHeadSub = pool.chain.SubscribeChainHeadEvent(pool.chainHeadCh)

	// // Start the event loop and return
	// pool.wg.Add(1)
	// go pool.loop()

	// return pool
// }
