package downloader

import (

)

// peerConnection represents an active peer from which hashes and blocks are retrieved.
type peerConnection struct {
	id string // Unique identifier of the peer

	headerIdle  int32 // Current header activity state of the peer (idle = 0, active = 1)
	blockIdle   int32 // Current block activity state of the peer (idle = 0, active = 1)
	receiptIdle int32 // Current receipt activity state of the peer (idle = 0, active = 1)
	stateIdle   int32 // Current node data activity state of the peer (idle = 0, active = 1)

	headerThroughput  float64 // Number of headers measured to be retrievable per second
	blockThroughput   float64 // Number of blocks (bodies) measured to be retrievable per second
	receiptThroughput float64 // Number of receipts measured to be retrievable per second
	stateThroughput   float64 // Number of node data pieces measured to be retrievable per second

	rtt time.Duration // Request round trip time to track responsiveness (QoS)

	headerStarted  time.Time // Time instance when the last header fetch was started
	blockStarted   time.Time // Time instance when the last block (body) fetch was started
	receiptStarted time.Time // Time instance when the last receipt fetch was started
	stateStarted   time.Time // Time instance when the last node data fetch was started

	lacking map[common.Hash]struct{} // Set of hashes not to request (didn't have previously)

	peer Peer

	version int        // Eth protocol version number to switch strategies
	log     log.Logger // Contextual logger to add extra infos to peer logs
	lock    sync.RWMutex
}

// LightPeer encapsulates the methods required to synchronise with a remote light peer.
type LightPeer interface {
	Head() (common.Hash, *big.Int)
	RequestHeadersByHash(common.Hash, int, int, bool) error
	RequestHeadersByNumber(uint64, int, int, bool) error
}

// Peer encapsulates the methods required to synchronise with a remote full peer.
type Peer interface {
	LightPeer
	RequestBodies([]common.Hash) error
	RequestReceipts([]common.Hash) error
	RequestNodeData([]common.Hash) error
}

// peerSet represents the collection of active peer participating in the chain
// download procedure.
type peerSet struct {
	peers        map[string]*peerConnection
	// newPeerFeed  event.Feed
	// peerDropFeed event.Feed
	lock         sync.RWMutex
}

// newPeerSet creates a new peer set top track the active download sources.
func newPeerSet() *peerSet {
	return &peerSet{
		peers: make(map[string]*peerConnection),
	}
}
