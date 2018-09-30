package params

import "math/big"

const (
	MaximumExtraDataSize  uint64 = 32    // Maximum size extra data may be after Genesis.
)

var GenesisDifficulty      = big.NewInt(4) // Difficulty of the Genesis block.
