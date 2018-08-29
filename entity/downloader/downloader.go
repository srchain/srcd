package downloader

import (
	"time"
)

type Downloader struct {
	mode SyncMode       // Synchronisation mode defining the strategy used (per sync cycle)
	// mux  *event.TypeMux // Event multiplexer to announce sync operation events

	queue   *queue   // Scheduler for selecting the hashes to download
	peers   *peerSet // Set of active peers from which download can proceed
	stateDB ethdb.Database

	// rttEstimate   uint64 // Round trip time to target for download requests
	// rttConfidence uint64 // Confidence in the estimated RTT (unit: millionths to allow atomic ops)

	// Statistics
	syncStatsChainOrigin uint64 // Origin block number where syncing started at
	syncStatsChainHeight uint64 // Highest block number known when syncing started
	syncStatsState       stateSyncStats
	syncStatsLock        sync.RWMutex // Lock protecting the sync stats fields

	blockchain BlockChain

	// Callbacks
	dropPeer peerDropFn // Drops a peer for misbehaving

	// Status
	synchroniseMock func(id string, hash common.Hash) error // Replacement for synchronise during testing
	synchronising   int32
	notified        int32
	committed       int32

	// Channels
	headerCh      chan dataPack        // [eth/62] Channel receiving inbound block headers
	bodyCh        chan dataPack        // [eth/62] Channel receiving inbound block bodies
	receiptCh     chan dataPack        // [eth/63] Channel receiving inbound receipts
	bodyWakeCh    chan bool            // [eth/62] Channel to signal the block body fetcher of new tasks
	receiptWakeCh chan bool            // [eth/63] Channel to signal the receipt fetcher of new tasks
	headerProcCh  chan []*types.Header // [eth/62] Channel to feed the header processor new tasks

	// for stateFetcher
	stateSyncStart chan *stateSync
	trackStateReq  chan *stateReq
	stateCh        chan dataPack // [eth/63] Channel receiving inbound node state data

	// Cancellation and termination
	cancelPeer string         // Identifier of the peer currently being used as the master (cancel on drop)
	cancelCh   chan struct{}  // Channel to cancel mid-flight syncs
	cancelLock sync.RWMutex   // Lock to protect the cancel channel and peer in delivers
	cancelWg   sync.WaitGroup // Make sure all fetcher goroutines have exited.

	quitCh   chan struct{} // Quit channel to signal termination
	quitLock sync.RWMutex  // Lock to prevent double closes

	// Testing hooks
	syncInitHook     func(uint64, uint64)  // Method to call upon initiating a new sync run
	bodyFetchHook    func([]*types.Header) // Method to call upon starting a block body fetch
	receiptFetchHook func([]*types.Header) // Method to call upon starting a receipt fetch
	chainInsertHook  func([]*fetchResult)  // Method to call upon inserting a chain of blocks (possibly in multiple invocations)
}

// LightChain encapsulates functions required to synchronise a light chain.
type LightChain interface {
	// HasHeader verifies a header's presence in the local chain.
	HasHeader(common.Hash, uint64) bool

	// GetHeaderByHash retrieves a header from the local chain.
	GetHeaderByHash(common.Hash) *types.Header

	// CurrentHeader retrieves the head header from the local chain.
	CurrentHeader() *types.Header

	// GetTd returns the total difficulty of a local block.
	GetTd(common.Hash, uint64) *big.Int

	// InsertHeaderChain inserts a batch of headers into the local chain.
	InsertHeaderChain([]*types.Header, int) (int, error)

	// Rollback removes a few recently added elements from the local chain.
	Rollback([]common.Hash)
}

// BlockChain encapsulates functions required to sync a (full or fast) blockchain.
type BlockChain interface {
	LightChain

	// HasBlock verifies a block's presence in the local chain.
	HasBlock(common.Hash, uint64) bool

	// GetBlockByHash retrieves a block from the local chain.
	GetBlockByHash(common.Hash) *types.Block

	// CurrentBlock retrieves the head block from the local chain.
	CurrentBlock() *types.Block

	// CurrentFastBlock retrieves the head fast block from the local chain.
	CurrentFastBlock() *types.Block

	// FastSyncCommitHead directly commits the head block to a certain entity.
	FastSyncCommitHead(common.Hash) error

	// InsertChain inserts a batch of blocks into the local chain.
	InsertChain(types.Blocks) (int, error)

	// InsertReceiptChain inserts a batch of receipts into the local chain.
	InsertReceiptChain(types.Blocks, []types.Receipts) (int, error)
}

// New creates a new downloader to fetch hashes and blocks from remote peers.
func New(mode SyncMode, stateDb ethdb.Database, mux *event.TypeMux, chain BlockChain, dropPeer peerDropFn) *Downloader {
	dl := &Downloader{
		mode:           mode,
		stateDB:        stateDb,
		mux:            mux,
		queue:          newQueue(),
		peers:          newPeerSet(),
		blockchain:     chain,
		dropPeer:       dropPeer,
		headerCh:       make(chan dataPack, 1),
		bodyCh:         make(chan dataPack, 1),
		receiptCh:      make(chan dataPack, 1),
		bodyWakeCh:     make(chan bool, 1),
		receiptWakeCh:  make(chan bool, 1),
		headerProcCh:   make(chan []*types.Header, 1),
		quitCh:         make(chan struct{}),
		stateCh:        make(chan dataPack),
		stateSyncStart: make(chan *stateSync),
		syncStatsState: stateSyncStats{
			processed: rawdb.ReadFastTrieProgress(stateDb),
		},
		trackStateReq: make(chan *stateReq),
	}
	go dl.stateFetcher()
	return dl
}
