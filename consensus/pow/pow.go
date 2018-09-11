package pow

import (
	"sync"
	"math/rand"
	"srcd/log"
)

// Pow is a consensus engine based on proot-of-work
type Pow struct {
	rand     *rand.Rand    // Properly seeded random source for nonces
	threads  int           // Number of threads to mine on if mining
	update   chan struct{} // Notification channel to update mining parameters
	lock     sync.Mutex    // Ensures thread safety for the in-memory caches and mining fields
}

// New creates a full sized PoW scheme.
func New() *Pow {
	return &Pow{ update:   make(chan struct{}) }
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
