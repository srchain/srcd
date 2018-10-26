package server

import (
	"fmt"
	"runtime"
	"sync"

	"srcd/accounts"
	"srcd/common/common"
	"srcd/common/hexutil"
	"srcd/consensus"
	"srcd/consensus/pow"
	"srcd/core/blockchain"
	"srcd/core/mempool"
	"srcd/database"
	"srcd/log"
	"srcd/miner"
	"srcd/node"
	"srcd/params"
	"srcd/rlp"
)

// Server implements the full node service.
type Server struct {
	config *Config
	// chainConfig *params.ChainConfig

	// Channel for shutting down the service
	shutdownChan chan bool

	// Handlers
	txPool          *mempool.TxPool
	blockchain      *blockchain.BlockChain
	protocolManager *ProtocolManager

	// DB interfaces
	chainDb database.Database // Block chain database

	// eventMux       *event.TypeMux
	engine         consensus.Engine
	accountManager *accounts.Manager

	// bloomRequests chan chan *bloombits.Retrieval // Channel receiving bloom data retrieval requests
	// bloomIndexer  *core.ChainIndexer             // Bloom indexer operating during block imports

	// APIBackend *EthAPIBackend

	miner    *miner.Miner
	coinbase common.Address

	// networkID     uint64
	// netRPCService *ethapi.PublicNetAPI

	lock sync.RWMutex
}

// New creates a new Server object
func New(ctx *node.ServiceContext, config *Config) (*Server, error) {
	chainDb, err := CreateDB(ctx, config, "chaindata")
	if err != nil {
		return nil, err
	}

	if _, genesisErr := blockchain.SetupGenesisBlock(chainDb, config.Genesis); genesisErr != nil {
		return nil, genesisErr
	}

	server := &Server{
		config:         config,
		chainDb:        chainDb,
		accountManager: ctx.AccountManager,
		engine:         CreateConsensusEngine(),
		shutdownChan:   make(chan bool),
		coinbase:       config.Coinbase,
	}

	server.blockchain, err = blockchain.NewBlockChain(chainDb, server.engine)
	if err != nil {
		return nil, err
	}

	// server.bloomIndexer.Start(eth.blockchain)

	// if config.TxPool.Journal != "" {
	// config.TxPool.Journal = ctx.ResolvePath(config.TxPool.Journal)
	// }
	server.txPool = mempool.NewTxPool(config.TxPool, server.blockchain)

	// if server.protocolManager, err = NewProtocolManager(eth.chainConfig, config.SyncMode, config.NetworkId, eth.eventMux, eth.txPool, eth.engine, eth.blockchain, chainDb); err != nil {
	// return nil, err
	// }

	server.miner = miner.New(server, server.engine)
	server.miner.SetExtra(makeExtraData(config.ExtraData))

	return server, nil
}

func makeExtraData(extra []byte) []byte {
	if len(extra) == 0 {
		// create default extradata
		extra, _ = rlp.EncodeToBytes([]interface{}{
			uint(params.VersionMajor<<16 | params.VersionMinor<<8 | params.VersionPatch),
			"srcd",
			runtime.Version(),
			runtime.GOOS,
		})
	}
	if uint64(len(extra)) > params.MaximumExtraDataSize {
		log.Warn("Miner extra data exceed limit", "extra", hexutil.Bytes(extra), "limit", params.MaximumExtraDataSize)
		extra = nil
	}
	return extra
}

// CreateDB creates the chain database.
func CreateDB(ctx *node.ServiceContext, config *Config, name string) (database.Database, error) {
	return ctx.OpenDatabase(name, config.DatabaseCache, config.DatabaseHandles)
}

// CreateConsensusEngine creates the required type of consensus engine instance for Server
func CreateConsensusEngine() consensus.Engine {
	engine := pow.New()
	engine.SetThreads(1)

	return engine
}

func (s *Server) Coinbase() (cb common.Address, err error) {
	s.lock.RLock()
	coinbase := s.coinbase
	s.lock.RUnlock()

	if coinbase != (common.Address{}) {
		return coinbase, nil
	}
	if wallets := s.AccountManager().Wallets(); len(wallets) > 0 {
		if accounts := wallets[0].Accounts(); len(accounts) > 0 {
			coinbase := accounts[0].Address

			s.lock.Lock()
			s.coinbase = coinbase
			s.lock.Unlock()

			log.Info("Coinbase automatically configured", "address", coinbase)
			return coinbase, nil
		}
	}
	return common.Address{}, fmt.Errorf("coinbase must be explicitly specified")
}

// StartMining starts the miner with the given number of CPU threads. If mining
// is already running, this method adjust the number of threads allowed to use.
func (s *Server) StartMining(threads int) error {
	// Update the thread count within the consensus engine
	type threaded interface {
		SetThreads(threads int)
	}
	if th, ok := s.engine.(threaded); ok {
		log.Info("Updated mining threads", "threads", threads)
		if threads == 0 {
			threads = -1 // Disable the miner from within
		}
		th.SetThreads(threads)
	}
	// If the miner was not running, initialize it
	if !s.IsMining() {
		// Configure the local mining address
		cb, err := s.Coinbase()
		if err != nil {
			log.Error("Cannot start mining without coinbase", "err", err)
			return fmt.Errorf("coinbase missing: %v", err)
		}

		// If mining is started, we can disable the transaction rejection mechanism
		// introduced to speed sync times.
		// atomic.StoreUint32(&s.protocolManager.acceptTxs, 1)

		go s.miner.Start(cb)
	}
	return nil
}

func (s *Server) IsMining() bool      { return s.miner.Mining() }

func (s *Server) AccountManager() *accounts.Manager  { return s.accountManager }
func (s *Server) BlockChain() *blockchain.BlockChain { return s.blockchain }
func (s *Server) TxPool() *mempool.TxPool            { return s.txPool }
func (s *Server) Engine() consensus.Engine           { return s.engine }
func (s *Server) ChainDb() database.Database         { return s.chainDb }

// Start implements node.Service, starting all internal goroutines needed by the
// Server protocol implementation.
// func (s *Server) Start(srvr *p2p.Server) error {
func (s *Server) Start() error {
	// // Start the bloom bits servicing goroutines
	// s.startBloomHandlers()

	// // Start the RPC service
	// s.netRPCService = ethapi.NewPublicNetAPI(srvr, s.NetVersion())

	// // Figure out a max peers count based on the server limits
	// maxPeers := srvr.MaxPeers
	// if s.config.LightServ > 0 {
	// if s.config.LightPeers >= srvr.MaxPeers {
	// return fmt.Errorf("invalid peer config: light peer count (%d) >= total peer count (%d)", s.config.LightPeers, srvr.MaxPeers)
	// }
	// maxPeers -= s.config.LightPeers
	// }

	// Start the networking layer and the light server if requested
	// s.protocolManager.Start(10)

	return nil
}

// Stop implements node.Service, terminating all internal goroutines used by the
// Server protocol.
func (s *Server) Stop() error {
	// s.bloomIndexer.Close()
	s.blockchain.Stop()
	// s.engine.Close()
	// s.protocolManager.Stop()
	// s.txPool.Stop()
	s.miner.Stop()
	// s.eventMux.Stop()

	s.chainDb.Close()
	close(s.shutdownChan)
	return nil
}
