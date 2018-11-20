package pow

import (
	"math/big"
	"math/rand"
	"sync"
	"github.com/srchain/srcd/core/types"
)

// two256 is a big integer representing 2^256
var two256 = new(big.Int).Exp(big.NewInt(2), big.NewInt(256), big.NewInt(0))

// Pow is a consensus engine based on proot-of-work
type Pow struct {
	rand     *rand.Rand        // Properly seeded random source for nonces
	threads  int               // Number of threads to mine on if mining
	update   chan struct{}     // Notification channel to update mining parameters

	resultCh chan *types.Block // Channel used by mining threads to return result

	// workCh       chan *types.Block // Notification channel to push new work to remote sealer
	// fetchWorkCh  chan *sealWork    // Channel used for remote sealer to fetch mining work
	// submitWorkCh chan *mineResult  // Channel used for remote sealer to submit their mining result

	lock      sync.Mutex      // Ensures thread safety for the in-memory caches and mining fields
}

// New creates a full sized PoW scheme.
func New() *Pow {
	pow := &Pow{
		update:   make(chan struct{}),
		resultCh: make(chan *types.Block),
	}

	return pow
}

// Threads returns the number of mining threads currently enabled. This doesn't
// necessarily mean that mining is running!
func (pow *Pow) Threads() int {
	pow.lock.Lock()
	defer pow.lock.Unlock()

	return pow.threads
}

// SetThreads updates the number of mining threads currently enabled. Calling
// this method does not start mining, only sets the thread count. If zero is
// specified, the miner will use all cores of the machine. Setting a thread
// count below zero is allowed and will cause the miner to idle, without any
// work being done.
func (pow *Pow) SetThreads(threads int) {
	pow.lock.Lock()
	defer pow.lock.Unlock()

	// Update the threads and ping any running seal to pull in any changes
	pow.threads = threads
	select {
	case pow.update <- struct{}{}:
	default:
	}
}
