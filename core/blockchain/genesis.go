package blockchain

import (
	"fmt"
	"github.com/srchain/srcd/errors"
	"math/big"
	"time"

	"github.com/srchain/srcd/common/common"
	"github.com/srchain/srcd/common/hexutil"
	"github.com/srchain/srcd/core/rawdb"
	"github.com/srchain/srcd/core/types"
	"github.com/srchain/srcd/database"
	"github.com/srchain/srcd/log"
	"github.com/srchain/srcd/params"
)



var errGenesisNoConfig = errors.New("genesis has no chain configuration")


// Genesis specifies the header fields, state of a genesis block.
type Genesis struct {
	Config     *params.ChainConfig `json:"config"`
	Nonce      uint64         `json:"nonce"`
	Timestamp  uint64         `json:"timestamp"`
	ExtraData  []byte         `json:"extraData"`
	Difficulty *big.Int       `json:"difficulty" gencodec:"required"`
	Coinbase   common.Address `json:"coinbase"`
	// Alloc      GenesisAlloc        `json:"alloc"      gencodec:"required"`

	// These fields are used for consensus tests. Please don't use them
	// in actual genesis blocks.
	Number     uint64      `json:"number"`
	ParentHash common.Hash `json:"parentHash"`
}

// GenesisAlloc specifies the initial state that is part of the genesis block.
// type GenesisAlloc map[common.Address]GenesisAccount

// GenesisAccount is an account in the state of the genesis block.
// type GenesisAccount struct {
// Code       []byte                      `json:"code,omitempty"`
// Storage    map[common.Hash]common.Hash `json:"storage,omitempty"`
// Balance    *big.Int                    `json:"balance" gencodec:"required"`
// Nonce      uint64                      `json:"nonce,omitempty"`
// PrivateKey []byte                      `json:"secretKey,omitempty"` // for tests
// }

// GenesisMismatchError is raised when trying to overwrite an existing
// genesis block with an incompatible one.
type GenesisMismatchError struct {
	Stored, New common.Hash
}

func (e *GenesisMismatchError) Error() string {
	return fmt.Sprintf("database already contains an incompatible genesis block (have %x, new %x)", e.Stored[:8], e.New[:8])
}


// SetupGenesisBlock writes or updates the genesis block in db.
// The block that will be used is:
//
//                          genesis == nil       genesis != nil
//                       +------------------------------------------
//     db has no genesis |  main-net default  |  genesis
//     db has genesis    |  from DB           |  genesis (if compatible)
//
// The stored chain configuration will be updated if it is compatible (i.e. does not
// specify a fork block below the local head block). In case of a conflict, the
// error is a *params.ConfigCompatError and the new, unwritten config is returned.
//
// The returned chain configuration is never nil.

func SetupGenesisBlock(db database.Database, genesis *Genesis) (*params.ChainConfig, common.Hash, error) {
	if genesis != nil && genesis.Config == nil {
		return params.AllPowProtocolChanges, common.Hash{}, errGenesisNoConfig
	}

	// Just commit the new block if there is no stored genesis block.
	stored := rawdb.ReadCanonicalHash(db, 0)
	if (stored == common.Hash{}) {
		if genesis == nil {
			log.Info("Writing default main-net genesis block")
			genesis = DefaultGenesisBlock()
		} else {
			log.Info("Writing custom genesis block")
		}
		block, err := genesis.Commit(db)
		return genesis.Config, block.Hash(), err
	}

	// Check whether the genesis block is already written.
	if genesis != nil {
		hash := genesis.ToBlock().Hash()
		if hash != stored {
			return genesis.Config, hash, &GenesisMismatchError{stored, hash}
		}
	}

	// Get the existing chain configuration.
	newcfg := genesis.configOrDefault(stored)
	storedcfg := rawdb.ReadChainConfig(db, stored)
	if storedcfg == nil {
		log.Warn("Found genesis block without chain config")
		rawdb.WriteChainConfig(db, stored, newcfg)
		return newcfg, stored, nil
	}
	// Special case: don't change the existing config of a non-mainnet chain if no new
	// config is supplied. These chains would get AllProtocolChanges (and a compat error)
	// if we just continued here.
	if genesis == nil && stored != params.MainnetGenesisHash {
		return storedcfg, stored, nil
	}

	// Check config compatibility and write the config. Compatibility errors
	// are returned to the caller unless we're already at block zero.
	height := rawdb.ReadHeaderNumber(db, rawdb.ReadHeadHeaderHash(db))
	if height == nil {
		return newcfg, stored, fmt.Errorf("missing block number for head header hash")
	}
	compatErr := storedcfg.CheckCompatible(newcfg, *height)
	if compatErr != nil && *height != 0 && compatErr.RewindTo != 0 {
		return newcfg, stored, compatErr
	}
	rawdb.WriteChainConfig(db, stored, newcfg)
	return newcfg, stored, nil
}



// ToBlock creates the genesis block.
func (g *Genesis) ToBlock() *types.Block {
	head := &types.Header{
		Number:     new(big.Int).SetUint64(g.Number),
		Nonce:      types.EncodeNonce(g.Nonce),
		Time:       new(big.Int).SetUint64(g.Timestamp),
		ParentHash: g.ParentHash,
		Extra:      g.ExtraData,
		Difficulty: g.Difficulty,
		Coinbase:   g.Coinbase,
	}
	if g.Difficulty == nil {
		head.Difficulty = params.GenesisDifficulty
	}

	return types.NewBlock(head, nil)
}

func GenesisBlockForTesting(db database.Database) *types.Block {
	g := Genesis{}
	return g.MustCommit(db)
}

func (g *Genesis) MustCommit(db database.Database) *types.Block {
	block, err := g.Commit(db)
	if err != nil {
		panic(err)
	}
	return block
}


// Commit writes the block and state of a genesis specification to the database.
func (g *Genesis) Commit(db database.Database) (*types.Block, error) {
	block := g.ToBlock()
	if block.Number().Sign() != 0 {
		return nil, fmt.Errorf("can't commit genesis block with number > 0")
	}

	rawdb.WriteBlock(db, block)
	rawdb.WriteCanonicalHash(db, block.Hash(), block.NumberU64())
	rawdb.WriteHeadBlockHash(db, block.Hash())
	rawdb.WriteHeadHeaderHash(db, block.Hash())

	return block, nil
}

func (g *Genesis) configOrDefault(ghash common.Hash) *params.ChainConfig {
	switch {
	case g != nil:
		return g.Config
	case ghash == params.MainnetChainConfig:
		return params.MainnetChainConfig
	case ghash == params.TestnetGenesisHash:
		return params.TestnetChainConfig
	default:
		return params.AllEthashProtocolChanges
	}
}

// DefaultGenesisBlock returns main net genesis block.
func DefaultGenesisBlock() *Genesis {
	return &Genesis{
		Nonce:      2505,
		Timestamp:  uint64(time.Now().Unix()),
		ExtraData:  hexutil.MustDecode("0xf7f480febb057fb7176fabad3fc28b602052a4e76043a5d7cffe066a62daa84b"),
		Difficulty: big.NewInt(1000),
		//Alloc:      decodePrealloc(mainnetAllocData),
	}
}

// func decodePrealloc(data string) GenesisAlloc {
// var p []struct{ Addr, Balance *big.Int }
// if err := rlp.NewStream(strings.NewReader(data), 0).Decode(&p); err != nil {
// panic(err)
// }
// ga := make(GenesisAlloc, len(p))
// for _, account := range p {
// ga[common.BigToAddress(account.Addr)] = GenesisAccount{Balance: account.Balance}
// }
// return ga
// }
