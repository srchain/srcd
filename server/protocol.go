package server

import (
	"github.com/srchain/srcd/common/common"
	"github.com/srchain/srcd/core/types"
	"github.com/srchain/srcd/core"
	"github.com/srchain/srcd/event"
)

// Constants to match up protocol versions and messages
const (
	src63 = 63
)

// ProtocolName is the official short name of the protocol used during capability negotiation.
var ProtocolName = "src"

// ProtocolVersions are the upported versions of the eth protocol (first is primary).
var ProtocolVersions = []uint{src63}

// ProtocolLengths are the number of implemented message corresponding to different protocol versions.
var ProtocolLengths = []uint64{17, 8}

const ProtocolMaxMsgSize = 10 * 1024 * 1024 // Maximum cap on the size of a protocol message

// eth protocol message codes
const (

	// Protocol messages belonging to src/63
	StatusMsg          = 0x00
	NewBlockHashesMsg  = 0x01
	TxMsg              = 0x02
	GetBlockHeadersMsg = 0x03
	BlockHeadersMsg    = 0x04
	GetBlockBodiesMsg  = 0x05
	BlockBodiesMsg     = 0x06
	NewBlockMsg        = 0x07
	GetNodeDataMsg 	   = 0x08
	NodeDataMsg        = 0x09
	GetReceiptsMsg     = 0x0A
	ReceiptsMsg        = 0x0B
)



type txPool interface {
	// AddRemotes should add the given transactions to the pool.
	AddRemotes([]*types.Transaction) []error

	// Pending should return pending transactions.
	// The slice should be modifiable by the caller.
	Pending() (map[common.Address]types.Transactions, error)

	// SubscribeNewTxsEvent should return an event subscription of
	// NewTxsEvent and send events to the given channel.
	SubscribeNewTxsEvent(chan<- core.NewTxsEvent) event.Subscription
}


// newBlockHashesData is the network packet for the block announcements.
type newBlockHashesData []struct {
	Hash   common.Hash // Hash of one particular block being announced
	Number uint64      // Number of one particular block being announced
}

