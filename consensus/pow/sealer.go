package pow

import (
	"encoding/binary"
	crand "crypto/rand"
	"math"
	"math/big"
	"math/rand"
	"runtime"
	"sync"

	"github.com/srchain/srcd/consensus"
	"github.com/srchain/srcd/core/types"
	"github.com/srchain/srcd/crypto/crypto"
	"github.com/srchain/srcd/log"
)

// Seal implements consensus.Engine, attempting to find a nonce that satisfies
// the block's difficulty requirements.
func (pow *Pow) Seal(chain consensus.ChainReader, block *types.Block, stop <-chan struct{}) (*types.Block, error) {
	// Create a runner and the multiple search threads it directs
	abort := make(chan struct{})

	pow.lock.Lock()
	threads := pow.threads
	if pow.rand == nil {
		seed, err := crand.Int(crand.Reader, big.NewInt(math.MaxInt64))
		if err != nil {
			pow.lock.Unlock()
			return nil, err
		}
		pow.rand = rand.New(rand.NewSource(seed.Int64()))
	}
	pow.lock.Unlock()
	if threads == 0 {
		threads = runtime.NumCPU()
	}
	if threads < 0 {
		threads = 0 // Allows disabling local mining without extra logic around local/remote
	}
	var pend sync.WaitGroup
	for i := 0; i < threads; i++ {
		pend.Add(1)
		go func(id int, nonce uint64) {
			defer pend.Done()
			pow.mine(block, id, nonce, abort, pow.resultCh)
		}(i, uint64(pow.rand.Int63()))
	}

	// Wait until sealing is terminated or a nonce is found
	var result *types.Block
	select {
	case <-stop:
		// Outside abort, stop all miner threads
		close(abort)
	case result = <-pow.resultCh:
		// One of the threads found a block, abort all others
		close(abort)
	case <-pow.update:
		// Thread count was changed on user request, restart
		close(abort)
		pend.Wait()
		return pow.Seal(chain, block, stop)
	}

	// Wait for all miners to terminate and return the block
	pend.Wait()
	return result, nil
}

// mine is the actual proof-of-work miner that searches for a nonce starting from
// seed that results in correct final block difficulty.
func (pow *Pow) mine(block *types.Block, id int, seed uint64, abort chan struct{}, found chan *types.Block) {
	// Extract some data from the header
	var (
		header = block.Header()
		hash   = header.HashNoNonce().Bytes()
		target = new(big.Int).Div(two256, header.Difficulty)
	)

	// Start generating random nonces until we abort or find a good one
	var nonce = seed

	logger := log.New("miner", id)
	logger.Trace("Started PoW search for new nonces", "seed", seed)
search:
	for {
		select {
		case <-abort:
			// Mining terminated
			logger.Trace("PoW nonce search aborted", "attempts", nonce-seed)
			break search

		default:
			// Compute the PoW value of this nonce
			result := hashimoto(hash, nonce)

			if new(big.Int).SetBytes(result).Cmp(target) <= 0 {
				// Correct nonce found, create a new header with it
				header = types.CopyHeader(header)
				header.Nonce = types.EncodeNonce(nonce)

				// Seal and return a block (if still needed)
				select {
				case found <- block.WithSeal(header):
					logger.Trace("PoW nonce found and reported", "attempts", nonce-seed, "nonce", nonce)
				case <-abort:
					logger.Trace("PoW nonce found but discarded", "attempts", nonce-seed, "nonce", nonce)
				}
				break search
			}

			nonce++
		}
	}
}

func hashimoto(hash []byte, nonce uint64) []byte {
	// Combine header+nonce into a 64 byte seed
	seed := make([]byte, 40)
	copy(seed, hash)
	binary.LittleEndian.PutUint64(seed[32:], nonce)

	return crypto.Keccak256(seed)
}
