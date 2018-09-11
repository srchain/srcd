package blockchain

import (
	mrand "math/rand"
	"sync/atomic"

	"srcd/core/types"
	"srcd/core/rawdb"
	"srcd/database"
	"srcd/common/common"
	"srcd/consensus"
)

// HeaderChain implements the basic block header chain logic that is shared by
// core.BlockChain and light.LightChain. It is not usable in itself, only as
// a part of either structure.
type HeaderChain struct {
	chainDb       database.Database
	genesisHeader *types.Header

	currentHeader     atomic.Value // Current head of the header chain (may be above the block chain!)
	currentHeaderHash common.Hash  // Hash of the current head of the header chain (prevent recomputing all the time)

	// headerCache *lru.Cache // Cache for the most recent block headers
	// tdCache     *lru.Cache // Cache for the most recent block total difficulties
	// numberCache *lru.Cache // Cache for the most recent block numbers

	procInterrupt func() bool

	rand   *mrand.Rand
	engine consensus.Engine
}

// NewHeaderChain creates a new HeaderChain structure.
func NewHeaderChain(chainDb database.Database, engine consensus.Engine, procInterrupt func() bool) (*HeaderChain, error) {
	// Seed a fast but crypto originating random generator
	seed, err := crand.Int(crand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		return nil, err
	}

	hc := &HeaderChain{
		chainDb:       chainDb,
		procInterrupt: procInterrupt,
		rand:          mrand.New(mrand.NewSource(seed.Int64())),
		engine:        engine,
	}

	hc.genesisHeader = hc.GetHeaderByNumber(0)
	if hc.genesisHeader == nil {
		return nil, ErrNoGenesis
	}

	hc.currentHeader.Store(hc.genesisHeader)
	if head := rawdb.ReadHeadBlockHash(chainDb); head != (common.Hash{}) {
		if chead := hc.GetHeaderByHash(head); chead != nil {
			hc.currentHeader.Store(chead)
		}
	}
	hc.currentHeaderHash = hc.CurrentHeader().Hash()

	return hc, nil
}

// GetBlockNumber retrieves the block number belonging to the given hash
// from the cache or database
func (hc *HeaderChain) GetBlockNumber(hash common.Hash) *uint64 {
	if cached, ok := hc.numberCache.Get(hash); ok {
		number := cached.(uint64)
		return &number
	}
	number := rawdb.ReadHeaderNumber(hc.chainDb, hash)
	if number != nil {
		hc.numberCache.Add(hash, *number)
	}
	return number
}

// GetHeader retrieves a block header from the database by hash and number.
func (hc *HeaderChain) GetHeader(hash common.Hash, number uint64) *types.Header {
	header := rawdb.ReadHeader(hc.chainDb, hash, number)
	if header == nil {
		return nil
	}

	return header
}

// GetHeaderByHash retrieves a block header from the database by hash.
func (hc *HeaderChain) GetHeaderByHash(hash common.Hash) *types.Header {
	number := hc.GetBlockNumber(hash)
	if number == nil {
		return nil
	}
	return hc.GetHeader(hash, *number)
}

// GetHeaderByNumber retrieves a block header from the database by number.
func (hc *HeaderChain) GetHeaderByNumber(number uint64) *types.Header {
	hash := rawdb.ReadCanonicalHash(hc.chainDb, number)
	if hash == (common.Hash{}) {
		return nil
	}
	return hc.GetHeader(hash, number)
}

// CurrentHeader retrieves the current head header of the canonical chain. The
// header is retrieved from the HeaderChain's internal cache.
func (hc *HeaderChain) CurrentHeader() *types.Header {
	return hc.currentHeader.Load().(*types.Header)
}

// SetCurrentHeader sets the current head header of the canonical chain.
func (hc *HeaderChain) SetCurrentHeader(head *types.Header) {
	rawdb.WriteHeadHeaderHash(hc.chainDb, head.Hash())

	hc.currentHeader.Store(head)
	hc.currentHeaderHash = head.Hash()
}

// DeleteCallback is a callback function that is called by SetHead before
// each header is deleted.
type DeleteCallback func(rawdb.DatabaseDeleter, common.Hash, uint64)

// SetHead rewinds the local chain to a new head. Everything above the new head
// will be deleted and the new one set.
func (hc *HeaderChain) SetHead(head uint64, delFn DeleteCallback) {
	height := uint64(0)
	if hdr := hc.CurrentHeader(); hdr != nil {
		height = hdr.Number.Uint64()
	}

	batch := hc.chainDb.NewBatch()
	for hdr := hc.CurrentHeader(); hdr != nil && hdr.Number.Uint64() > head; hdr = hc.CurrentHeader() {
		hash := hdr.Hash()
		num := hdr.Number.Uint64()
		if delFn != nil {
			delFn(batch, hash, num)
		}
		rawdb.DeleteHeader(batch, hash, num)

		hc.currentHeader.Store(hc.GetHeader(hdr.ParentHash, hdr.Number.Uint64()-1))
	}

	for i := height; i > head; i-- {
		rawdb.DeleteCanonicalHash(batch, i)
	}
	batch.Write()

	if hc.CurrentHeader() == nil {
		hc.currentHeader.Store(hc.genesisHeader)
	}
	hc.currentHeaderHash = hc.CurrentHeader().Hash()

	rawdb.WriteHeadHeaderHash(hc.chainDb, hc.currentHeaderHash)
}

// SetGenesis sets a new genesis block header for the chain
func (hc *HeaderChain) SetGenesis(head *types.Header) {
	hc.genesisHeader = head
}
