package discover5

import (
	"go-ethereum/common/mclock"
	"github.com/srchain/srcd/common/common"
	"time"
	"github.com/srchain/srcd/log"
	"github.com/srchain/srcd/crypto/crypto"
	"encoding/binary"
	"math"
)



const (
	ticketTimeBucketLen = time.Minute
	timeWindow          = 10 // * ticketTimeBucketLen
	wantTicketsInWindow = 10
	collectFrequency    = time.Second * 30
	registerFrequency   = time.Second * 60
	maxCollectDebt      = 10
	maxRegisterDebt     = 5
	keepTicketConst     = time.Minute * 10
	keepTicketExp       = time.Minute * 5
	targetWaitTime      = time.Minute * 10
	topicQueryTimeout   = time.Second * 5
	topicQueryResend    = time.Minute
	// topic radius detection
	maxRadius           = 0xffffffffffffffff
	radiusTC            = time.Minute * 20
	radiusBucketsPerBit = 8
	minSlope            = 1
	minPeakSize         = 40
	maxNoAdjust         = 20
	lookupWidth         = 8
	minRightSum         = 20
	searchForceQuery    = 4
)


type topicRadiusEvent int

type timeBucket int

type topicTickets struct {

	buckets map[timeBucket] []ticketRef
	nextLookup mclock.AbsTime
	nextReg mclock.AbsTime
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

func (ref ticketRef) topicRegTime() mclock.AbsTime {
	return ref.t.regTime[ref.idx]

}
func (ref ticketRef) topic() Topic {
	return ref.t.topics[ref.idx]
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

type reqInfo struct {
	pingHash []byte
	lookup lookupInfo
	time mclock.AbsTime
}

type ticketStore struct {
	radius map[Topic]*topicRadius
	tickets map[Topic]*topicTickets


	regQueue []Topic            // Topic registration queue for round robin attempts
	regSet   map[Topic]struct{} // Topic registration queue contents for fast filling

	nodes map[*Node]*ticket
	nodeLastReq map[*Node] reqInfo

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

func (s *ticketStore) nextFilteredTicket() (*ticketRef, time.Duration) {
	now := mclock.Now()
	for {
		ticket, wait := s.nextRegisterableTicket()
		if ticket == nil {
			return ticket, wait
		}
		log.Trace("Found discovery ticket to register","node",ticket.t.node,"serial",ticket.t.serial,"wait",wait)
		regTime := now + mclock.AbsTime(wait)
		topic := ticket.t.topics[ticket.idx]
		if s.tickets[topic] != nil && regTime >= s.tickets[topic].nextReg {
			return ticket, wait
		}
		s.removeTicketRef(*ticket)
	}
}

func (s *ticketStore) nextRegisterableTicket() (*ticketRef, time.Duration) {
	now := mclock.Now()
	if s.nextTicketCached != nil {
		return s.nextTicketCached, time.Duration(s.nextTicketCached.topicRegTime() - now)
	}
	for bucket := s.lastBucketFetched; ; bucket++ {
		var (
			empty = true
			nextTicket ticketRef
		)
		for _, tickets := range s.tickets {
			if len(tickets.buckets) != 0 {
				empty = false

				list := tickets.buckets[bucket]
				for _, ref := range list {
					if nextTicket.t == nil || ref.topicRegTime() < nextTicket.topicRegTime() {
						nextTicket = ref
					}
				}
			}
		}
		if empty {
			return nil, 0
		}
		if nextTicket.t != nil {
			s.nextTicketCached = &nextTicket
			return &nextTicket, time.Duration(nextTicket.topicRegTime() - now)
		}
		s.lastBucketFetched = bucket
	}
}

func (s *ticketStore) removeTicketRef(ref ticketRef) {
	log.Trace("Removing discovery ticket reference", "node", ref.t.node.ID, "serial", ref.t.serial)
	s.nextTicketCached = nil
	topic := ref.topic()
	tickets := s.tickets[topic]
	if tickets == nil {
		log.Trace("Removing tickets from unknown topic","topic",topic)
		return
	}
	bucket := timeBucket(ref.t.regTime[ref.idx] / mclock.AbsTime(ticketTimeBucketLen))
	list := tickets.buckets[bucket]
	idx := -1
	for i, bt := range list {
		if bt.t == ref.t {
			idx = i
			break
		}
	}
	if idx == -1 {
		panic(nil)
	}
	list = append(list[:idx],list[idx+1:]...)
	if len(list) != 0 {
		tickets.buckets[bucket] = list
	} else {
		delete(tickets.buckets,bucket)
	}
	ref.t.refCnt--
	if ref.t.refCnt == 0 {
		delete(s.nodes,ref.t.node)
		delete(s.nodeLastReq, ref.t.node)
	}

}

// removeRegisterTopic deletes all tickets for the given topic.
func (s *ticketStore) removeRegisterTopic(topic Topic) {
	log.Trace("Removing discovery topic", "topic",topic)
	if s.tickets[topic] == nil {
		log.Warn("Removing non-existent discovery topic","topic",topic)
		return
	}
	for _, list := range s.tickets[topic].buckets {
		for _, ref := range list {
			ref.t.refCnt--
			if ref.t.refCnt == 0 {
				delete(s.nodes,ref.t.node)
				delete(s.nodeLastReq,ref.t.node)
			}
		}
	}
	delete(s.tickets,topic)
}
func (s *ticketStore) addTopic(topic Topic, register bool) {
	log.Trace("Adding discovery topic","topic","register",register)
	if s.radius[topic] == nil {
		s.radius[topic] = newTopicRadius(topic)
	}
	
}
func newTopicRadius(topic Topic) *topicRadius {
	topicHash := crypto.Keccak256Hash([]byte(topic))
	topicHashPrefix := binary.BigEndian.Uint64(topicHash[0:8])
	return &topicRadius{
		topic:topic,
		topicHashPrefix: topicHashPrefix ,
		radius: maxRadius,
		minRadius: maxRadius,

	}
}

func (r *topicRadius) getBucketIdx(addrHash common.Hash) int {
	prefix := binary.BigEndian.Uint64(addrHash[0:8])
	var log2 float64
	if prefix != r.topicHashPrefix {
		log2 = math.Log2(float64(prefix ^ r.topicHashPrefix))
	}
	bucket := int((64 - log2) * radiusBucketsPerBit)
	max := 64*radiusBucketsPerBit - 1
	if bucket > max {
		return max
	}
	if bucket < 0 {
		return 0
	}
	return bucket
}