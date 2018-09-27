package pow

import (
	"fmt"
	"math/big"
	"time"

	"srcd/core/types"
	"srcd/common/common"
	"srcd/params"
)

// Max time from current time allowed for blocks, before they're considered future blocks
var allowedFutureBlockTime          = 15 * time.Second

// Various error messages to mark blocks invalid. These should be private to
// prevent engine specific errors from being referenced in the remainder of the
// codebase, inherently breaking if the engine is swapped out. Please put common
// error types into the consensus package.
var (
	errZeroBlockTime     = errors.New("timestamp equals parent's")
)

// Author implements consensus.Engine, returning the header's coinbase as the
// proof-of-work verified author of the block.
func (pow *Pow) Author(header *types.Header) (common.Address, error) {
	return header.Coinbase, nil
}

// VerifyHeader checks whether a header conforms to the consensus rules of the PoW engine.
func (pow *Pow) VerifyHeader(chain consensus.ChainReader, header *types.Header, seal bool) error {
	// Short circuit if the header is known, or it's parent not
	number := header.Number.Uint64()
	if chain.GetHeader(header.Hash(), number) != nil {
		return nil
	}
	parent := chain.GetHeader(header.ParentHash, number-1)
	if parent == nil {
		return consensus.ErrUnknownAncestor
	}

	// Sanity checks passed, do a proper verification
	return pow.verifyHeader(chain, header, parent, seal)
}

// VerifyHeaders is similar to VerifyHeader, but verifies a batch of headers
// concurrently. The method returns a quit channel to abort the operations and
// a results channel to retrieve the async verifications.
func (pow *Pow) VerifyHeaders(chain consensus.ChainReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error) {
	// Spawn as many workers as allowed threads
	workers := runtime.GOMAXPROCS(0)
	if len(headers) < workers {
		workers = len(headers)
	}

	// Create a task channel and spawn the verifiers
	var (
		inputs = make(chan int)
		done   = make(chan int, workers)
		errors = make([]error, len(headers))
		abort  = make(chan struct{})
	)
	for i := 0; i < workers; i++ {
		go func() {
			for index := range inputs {
				errors[index] = pow.verifyHeaderWorker(chain, headers, seals, index)
				done <- index
			}
		}()
	}

	errorsOut := make(chan error, len(headers))
	go func() {
		defer close(inputs)
		var (
			in, out = 0, 0
			checked = make([]bool, len(headers))
			inputs  = inputs
		)
		for {
			select {
			case inputs <- in:
				if in++; in == len(headers) {
					// Reached end of headers. Stop sending to workers.
					inputs = nil
				}
			case index := <-done:
				for checked[index] = true; checked[out]; out++ {
					errorsOut <- errors[out]
					if out == len(headers)-1 {
						return
					}
				}
			case <-abort:
				return
			}
		}
	}()

	return abort, errorsOut
}

func (pow *Pow) verifyHeaderWorker(chain consensus.ChainReader, headers []*types.Header, seals []bool, index int) error {
	var parent *types.Header
	if index == 0 {
		parent = chain.GetHeader(headers[0].ParentHash, headers[0].Number.Uint64()-1)
	} else if headers[index-1].Hash() == headers[index].ParentHash {
		parent = headers[index-1]
	}
	if parent == nil {
		return consensus.ErrUnknownAncestor
	}

	if chain.GetHeader(headers[index].Hash(), headers[index].Number.Uint64()) != nil {
		return nil
	}
	return pow.verifyHeader(chain, headers[index], parent, seals[index])
}

// verifyHeader checks whether a header conforms to the consensus rules of the PoW engine.
func (pow *Pow) verifyHeader(chain consensus.ChainReader, header, parent *types.Header, seal bool) error {
	// Ensure that the header's extra-data section is of a reasonable size
	if uint64(len(header.Extra)) > params.MaximumExtraDataSize {
		return fmt.Errorf("extra-data too long: %d > %d", len(header.Extra), params.MaximumExtraDataSize)
	}

	// Verify the header's timestamp
	if header.Time.Cmp(big.NewInt(time.Now().Add(allowedFutureBlockTime).Unix())) > 0 {
		return consensus.ErrFutureBlock
	}
	if header.Time.Cmp(parent.Time) <= 0 {
		return errZeroBlockTime
	}

	// Verify the block's difficulty based in it's timestamp and parent's difficulty
	expected := pow.CalcDifficulty(header.Time.Uint64(), parent)

	if expected.Cmp(header.Difficulty) != 0 {
		return fmt.Errorf("invalid difficulty: have %v, want %v", header.Difficulty, expected)
	}

	// Verify that the block number is parent's +1
	if diff := new(big.Int).Sub(header.Number, parent.Number); diff.Cmp(big.NewInt(1)) != 0 {
		return consensus.ErrInvalidNumber
	}

	// Verify the engine specific seal securing the block
	if seal {
		if err := pow.VerifySeal(chain, header); err != nil {
			return err
		}
	}

	return nil
}

// CalcDifficulty is the difficulty adjustment algorithm. It returns
// the difficulty that a new block should have when created at time
// given the parent block's time and difficulty.
func (pow *Pow) CalcDifficulty(time uint64, parent *types.Header) *big.Int {
	return calcDifficulty(time, parent)
}

// calcDifficulty is the difficulty adjustment algorithm.
func calcDifficulty(time uint64, parent *types.Header) *big.Int {
	return big.NewInt(4)
}

// VerifySeal implements consensus.Engine, checking whether the given block satisfies
// the PoW difficulty requirements.
func (pow *Pow) VerifySeal(chain consensus.ChainReader, header *types.Header) error {
	return nil
}

// Prepare implements consensus.Engine, initializing the difficulty field of a
// header to conform to protocol.
func (pow *Pow) Prepare(chain consensus.ChainReader, header *types.Header) error {
	parent := chain.GetHeader(header.ParentHash, header.Number.Uint64()-1)
	if parent == nil {
		return consensus.ErrUnknownAncestor
	}
	header.Difficulty = pow.CalcDifficulty(header.Time.Uint64(), parent)
	return nil
}

// Finalize implements consensus.Engine, accumulating the block rewards, assembling the block.
func (pow *Pow) Finalize(header *types.Header, txs []*types.Transaction) (*types.Block, error) {
	// Accumulate any block rewards
	accumulateRewards(header)

	// Header seems complete, assemble into a block and return
	return types.NewBlock(header, txs), nil
}

// AccumulateRewards credits the coinbase of the given block with the mining reward.
func accumulateRewards(header *types.Header) {
	// Accumulate the rewards for the miner.
	reward := new(big.Int).Set(10000)

	// AddBalance(header.Coinbase, reward)
}
