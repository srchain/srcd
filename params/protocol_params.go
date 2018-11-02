package params

import "math/big"

const (
	MaximumExtraDataSize  uint64 = 32    // Maximum size extra data may be after Genesis.
)

var GenesisDifficulty      = big.NewInt(1000) // Difficulty of the Genesis block.
