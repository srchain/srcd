package discover5

import (
	"go-ethereum/common/mclock"
	"github.com/srchain/srcd/common/common"
)

type topicRadiusEvent int

type timeBucket int

type topicTickets struct {
	bucket map[timeBucket] []ticketRef
}

type ticket struct {
	topics  []Topic
	regTime []mclock.AbsTime // Per-topic local absolute time when the ticket can be used.

	// The serial number that was issued by the server.
	serial uint32
	// Used by registrar, tracks absolute time when the ticket was created.
	issueTime mclock.AbsTime

	// Fields used only by registrants
	node   *Node  // the registrar node that signed this ticket
	refCnt int    // tracks number of topics that will be registered using this ticket
	pong   []byte // encoded pong packet signed by the registrar
}

type ticketRef struct {
	t   *ticket
	idx int // index of the topic in t.topics and t.regTime
}



type lookupInfo struct {
	target       common.Hash
	topic        Topic
	radiusLookup bool
}


const (
	trOutside topicRadiusEvent = iota
	trInside
	trNoAdjust
	trCount
)

type sentQueru struct {
	sent mclock.AbsTime
	lookup lookupInfo
}

type ticketStore struct {
	radius map[Topic]*topicRadius
	tickets map[Topic]*topicTickets


	regQueue []Topic            // Topic registration queue for round robin attempts
	regSet   map[Topic]struct{} // Topic registration queue contents for fast filling

	nodes map[*Node]*ticket
	nodeLastReq map[Topic]struct{}

	lastBucketFetched timeBucket
	nextTicketCached	*ticketRef
	nextTicketReg	mclock.AbsTime

	searchTopicMap	map[Topic]searchTopic
	nextTopicQueryCleanup mclock.AbsTime
	queriesSent map[*Node]map[common.Hash]sentQueru


}

type searchTopic struct {
	foundChn chan<- *Node
}


type topicRadiusBucket struct {
	weights	[trCount] float64
	lastTime mclock.AbsTime
	value float64
	lookupSent map[common.Hash] mclock.AbsTime
}

type topicRadius struct {
	topic Topic
	topicHashPrefix	uint64
	radius, minRadius uint64
	buckets []topicRadiusBucket
	converged bool
	radiusLookuoCnt	int
}
