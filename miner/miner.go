package miner

import (
	"fmt"
	"sync/atomic"

	"github.com/srchain/srcd/common/common"
	"github.com/srchain/srcd/consensus"
	"github.com/srchain/srcd/core/blockchain"
	"github.com/srchain/srcd/core/mempool"
	"github.com/srchain/srcd/log"
	"github.com/srchain/srcd/params"
)

// Backend wraps all methods required for mining.
type Backend interface {
	BlockChain() *blockchain.BlockChain
	TxPool()     *mempool.TxPool
}

// Miner creates blocks and searches for proof-of-work values.
type Miner struct {
	// mux      *event.TypeMux
	worker   *worker
	coinbase common.Address
	server   Backend
	engine   consensus.Engine
	exitCh   chan struct{}

	canStart    int32 // can start indicates whether we can start the mining operation
	shouldStart int32 // should start indicates whether we should start after sync
}

func New(server Backend, engine consensus.Engine) *Miner {
	miner := &Miner{
		server: server,
		// mux:      mux,
		engine:   engine,
		worker:   newWorker(engine, server),
		canStart: 1,
	}
	go miner.update()

	return miner
}

// update keeps track of the downloader events. Please be aware that this is a one shot type of update loop.
// It's entered once and as soon as `Done` or `Failed` has been broadcasted the events are unregistered and
// the loop is exited. This to prevent a major security vuln where external parties can DOS you with blocks
// and halt your mining operation for as long as the DOS continues.
func (self *Miner) update() {
	// events := self.mux.Subscribe(downloader.StartEvent{}, downloader.DoneEvent{}, downloader.FailedEvent{})
	// defer events.Unsubscribe()

	// for {
	// select {
	// case ev := <-events.Chan():
	// if ev == nil {
	// return
	// }
	// switch ev.Data.(type) {
	// case downloader.StartEvent:
	// atomic.StoreInt32(&self.canStart, 0)
	// if self.Mining() {
	// self.Stop()
	// atomic.StoreInt32(&self.shouldStart, 1)
	// log.Info("Mining aborted due to sync")
	// }
	// case downloader.DoneEvent, downloader.FailedEvent:
	// shouldStart := atomic.LoadInt32(&self.shouldStart) == 1

	// atomic.StoreInt32(&self.canStart, 1)
	// atomic.StoreInt32(&self.shouldStart, 0)
	// if shouldStart {
	// self.Start(self.coinbase)
	// }
	// // stop immediately and ignore all further pending events
	// return
	// }
	// case <-self.exitCh:
	// return
	// }
	// }

	for {
		select {
		case <-self.exitCh:
			return
		default:
			shouldStart := atomic.LoadInt32(&self.shouldStart) == 1

			atomic.StoreInt32(&self.canStart, 1)
			atomic.StoreInt32(&self.shouldStart, 0)
			if shouldStart {
				self.Start(self.coinbase)
			}
			// stop immediately and ignore all further pending events
			return
		}
	}
}

func (self *Miner) Start(coinbase common.Address) {
	atomic.StoreInt32(&self.shouldStart, 1)
	self.SetCoinbase(coinbase)

	if atomic.LoadInt32(&self.canStart) == 0 {
		log.Info("Network syncing, will start miner afterwards")
		return
	}
	self.worker.start()
}

func (self *Miner) Stop() {
	self.worker.stop()
	atomic.StoreInt32(&self.shouldStart, 0)
}

func (self *Miner) Close() {
	self.worker.close()
	close(self.exitCh)
}

func (self *Miner) Mining() bool {
	return self.worker.isRunning()
}

func (self *Miner) SetExtra(extra []byte) error {
	if uint64(len(extra)) > params.MaximumExtraDataSize {
		return fmt.Errorf("Extra exceeds max length. %d > %v", len(extra), params.MaximumExtraDataSize)
	}
	self.worker.setExtra(extra)
	return nil
}

func (self *Miner) SetCoinbase(addr common.Address) {
	self.coinbase = addr
	self.worker.setCoinbase(addr)
}
