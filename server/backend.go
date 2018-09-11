package server

import (
	"fmt"

	"srcd/core/mempool"
	"srcd/core/blockchain"
	"srcd/database"
	"srcd/consensus"
	"srcd/consensus/pow"
)

type ProtocolManager struct {
	// networkID uint64

	// fastSync  uint32 // Flag whether fast sync is enabled (gets disabled if we already have blocks)
	acceptTxs uint32 // Flag whether we're considered synchronised (enables transaction processing)

	txpool      mempool.txPool
	blockchain  *blockchain.BlockChain
	chainconfig *params.ChainConfig
	maxPeers    int

	downloader *downloader.Downloader
	fetcher    *fetcher.Fetcher
	peers      *peerSet

	SubProtocols []p2p.Protocol

	eventMux      *event.TypeMux
	txsCh         chan core.NewTxsEvent
	txsSub        event.Subscription
	minedBlockSub *event.TypeMuxSubscription

	// channels for fetcher, syncer, txsyncLoop
	newPeerCh   chan *peer
	txsyncCh    chan *txsync
	quitSync    chan struct{}
	noMorePeers chan struct{}

	// wait group is used for graceful shutdowns during downloading
	// and processing
	wg sync.WaitGroup
}

// Server implements the full node service.
type Server struct {
	config          *Config
	// chainConfig *params.ChainConfig

	// Channel for shutting down the service
	shutdownChan chan bool

	// Handlers
	// txPool          *core.TxPool
	blockchain      *blockchain.BlockChain
	protocolManager *ProtocolManager

	// DB interfaces
	chainDb         database.Database // Block chain database

	// eventMux       *event.TypeMux
	engine          consensus.Engine
	// accountManager *accounts.Manager
	wallet		*wallet.Wallet

	// bloomRequests chan chan *bloombits.Retrieval // Channel receiving bloom data retrieval requests
	// bloomIndexer  *core.ChainIndexer             // Bloom indexer operating during block imports

	// APIBackend *EthAPIBackend

	miner           *miner.Miner
	coinbase        common.Address

	// networkID     uint64
	// netRPCService *ethapi.PublicNetAPI

	lock            sync.RWMutex
}

// New creates a new Server object
func New(ctx *node.ServiceContext, config *Config) (*Server, error) {
	chainDb, err := CreateDB(ctx, config, "chaindata")
	if err != nil {
		return nil, err
	}

	genesisHash, genesisErr := blockchain.SetupGenesisBlock(chainDb, config.Genesis)
	if genesisErr != nil {
		return nil, genesisErr
	}

	server := &Server{
		config:         config,
		chainDb:        chainDb,
		wallet:		ctx.Wallet,
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
	server.txPool = core.NewTxPool(config.TxPool, server.blockchain)

	if server.protocolManager, err = NewProtocolManager(eth.chainConfig, config.SyncMode, config.NetworkId, eth.eventMux, eth.txPool, eth.engine, eth.blockchain, chainDb); err != nil {
		return nil, err
	}
	server.miner = miner.New(server, server.chainConfig, server.EventMux(), server.engine)
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
	engine := pow.new()
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

func (s *Server) StartMining(local bool) error {
	cb, err := s.Coinbase()
	if err != nil {
		log.Error("Cannot start mining without coinbase", "err", err)
		return fmt.Errorf("coinbase missing: %v", err)
	}

	if local {
		// If local (CPU) mining is started, we can disable the transaction rejection
		// mechanism introduced to speed sync times. CPU mining on mainnet is ludicrous
		// so none will ever hit this path, whereas marking sync done on CPU mining
		// will ensure that private networks work in single miner mode too.
		atomic.StoreUint32(&s.protocolManager.acceptTxs, 1)
	}
	go s.miner.Start(cb)
	return nil
}

func (s *Server) BlockChain() *blockchain.BlockChain	{ return s.blockchain }
func (s *Server) Engine() consensus.Engine		{ return s.engine }

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
	s.protocolManager.Start(10)

	return nil
}

// Stop implements node.Service, terminating all internal goroutines used by the
// Server protocol.
func (s *Server) Stop() error {
	// s.bloomIndexer.Close()
	s.blockchain.Stop()
	s.protocolManager.Stop()
	// s.txPool.Stop()
	s.miner.Stop()
	// s.eventMux.Stop()

	s.chainDb.Close()
	close(s.shutdownChan)

	return nil
}
