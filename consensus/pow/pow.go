package pow

import (
	"sync"
	"math/rand"
)

// Config are the configuration parameters of the ethash.
type Config struct {
	CacheDir       string
	CachesInMem    int
	CachesOnDisk   int
	DatasetDir     string
	DatasetsInMem  int
	DatasetsOnDisk int
}

// Pow is a consensus engine based on proot-of-work
type Pow struct {
	config Config

	// caches   *lru // In memory caches to avoid regenerating too often
	// datasets *lru // In memory datasets to avoid regenerating too often

	// Mining related fields
	rand     *rand.Rand    // Properly seeded random source for nonces
	threads  int           // Number of threads to mine on if mining
	update   chan struct{} // Notification channel to update mining parameters
	// hashrate metrics.Meter // Meter tracking the average hashrate

	// The fields below are hooks for testing
	// shared    *Ethash       // Shared PoW verifier to avoid cache regeneration
	// fakeFail  uint64        // Block number which fails PoW check even in fake mode
	// fakeDelay time.Duration // Time delay to sleep for before returning from verify

	lock sync.Mutex // Ensures thread safety for the in-memory caches and mining fields
}

// New creates a full sized PoW scheme.
func New(config Config) *Pow {
	if config.CachesInMem <= 0 {
		log.Warn("One Pow cache must always be in memory", "requested", config.CachesInMem)
		config.CachesInMem = 1
	}
	if config.CacheDir != "" && config.CachesOnDisk > 0 {
		log.Info("Disk storage enabled for Pow caches", "dir", config.CacheDir, "count", config.CachesOnDisk)
	}
	if config.DatasetDir != "" && config.DatasetsOnDisk > 0 {
		log.Info("Disk storage enabled for Pow DAGs", "dir", config.DatasetDir, "count", config.DatasetsOnDisk)
	}
	return &Pow{
		config:   config,
		// caches:   newlru("cache", config.CachesInMem, newCache),
		// datasets: newlru("dataset", config.DatasetsInMem, newDataset),
		update:   make(chan struct{}),
		// hashrate: metrics.NewMeter(),
	}
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
