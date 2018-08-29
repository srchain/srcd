package downloader

import (
	"gopkg.in/karalabe/cookiejar.v2/collections/prque"
)

// fetchRequest is a currently running data retrieval operation.
type fetchRequest struct {
	Peer    *peerConnection // Peer to which the request was sent
	From    uint64          // [eth/62] Requested chain element index (used for skeleton fills only)
	Headers []*types.Header // [eth/62] Requested headers, sorted by request order
	Time    time.Time       // Time when the request was made
}

// fetchResult is a struct collecting partial results from data fetchers until
// all outstanding pieces complete and the result as a whole can be processed.
type fetchResult struct {
	Pending int         // Number of data fetches still pending
	Hash    common.Hash // Hash of the header to prevent recalculating

	Header       *types.Header
	Uncles       []*types.Header
	Transactions types.Transactions
	Receipts     types.Receipts
}

// queue represents hashes that are either need fetching or are being fetched
type queue struct {
	mode SyncMode // Synchronisation mode to decide on the block parts to schedule for fetching

	// Headers are "special", they download in batches, supported by a skeleton chain
	headerHead      common.Hash                    // [eth/62] Hash of the last queued header to verify order
	headerTaskPool  map[uint64]*types.Header       // [eth/62] Pending header retrieval tasks, mapping starting indexes to skeleton headers
	headerTaskQueue *prque.Prque                   // [eth/62] Priority queue of the skeleton indexes to fetch the filling headers for
	headerPeerMiss  map[string]map[uint64]struct{} // [eth/62] Set of per-peer header batches known to be unavailable
	headerPendPool  map[string]*fetchRequest       // [eth/62] Currently pending header retrieval operations
	headerResults   []*types.Header                // [eth/62] Result cache accumulating the completed headers
	headerProced    int                            // [eth/62] Number of headers already processed from the results
	headerOffset    uint64                         // [eth/62] Number of the first header in the result cache
	headerContCh    chan bool                      // [eth/62] Channel to notify when header download finishes

	// All data retrievals below are based on an already assembles header chain
	blockTaskPool  map[common.Hash]*types.Header // [eth/62] Pending block (body) retrieval tasks, mapping hashes to headers
	blockTaskQueue *prque.Prque                  // [eth/62] Priority queue of the headers to fetch the blocks (bodies) for
	blockPendPool  map[string]*fetchRequest      // [eth/62] Currently pending block (body) retrieval operations
	blockDonePool  map[common.Hash]struct{}      // [eth/62] Set of the completed block (body) fetches

	receiptTaskPool  map[common.Hash]*types.Header // [eth/63] Pending receipt retrieval tasks, mapping hashes to headers
	receiptTaskQueue *prque.Prque                  // [eth/63] Priority queue of the headers to fetch the receipts for
	receiptPendPool  map[string]*fetchRequest      // [eth/63] Currently pending receipt retrieval operations
	receiptDonePool  map[common.Hash]struct{}      // [eth/63] Set of the completed receipt fetches

	resultCache  []*fetchResult     // Downloaded but not yet delivered fetch results
	resultOffset uint64             // Offset of the first cached fetch result in the block chain
	resultSize   common.StorageSize // Approximate size of a block (exponential moving average)

	lock   *sync.Mutex
	active *sync.Cond
	closed bool
}

// newQueue creates a new download queue for scheduling block retrieval.
func newQueue() *queue {
	lock := new(sync.Mutex)
	return &queue{
		headerPendPool:   make(map[string]*fetchRequest),
		headerContCh:     make(chan bool),
		blockTaskPool:    make(map[common.Hash]*types.Header),
		blockTaskQueue:   prque.New(),
		blockPendPool:    make(map[string]*fetchRequest),
		blockDonePool:    make(map[common.Hash]struct{}),
		receiptTaskPool:  make(map[common.Hash]*types.Header),
		receiptTaskQueue: prque.New(),
		receiptPendPool:  make(map[string]*fetchRequest),
		receiptDonePool:  make(map[common.Hash]struct{}),
		resultCache:      make([]*fetchResult, blockCacheItems),
		active:           sync.NewCond(lock),
		lock:             lock,
	}
}
