package params

import "math/big"

const (
	MaximumExtraDataSize  uint64 = 32    // Maximum size extra data may be after Genesis.
	EpochDuration    uint64 = 30000 // Duration between proof-of-work epochs.
)

var GenesisDifficulty      = big.NewInt(1000) // Difficulty of the Genesis block.
