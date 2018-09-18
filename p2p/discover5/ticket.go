package discover5

import (
	"go-ethereum/common/mclock"
	"github.com/srchain/srcd/common/common"
	"time"
	"github.com/srchain/srcd/log"
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
	log.Trace("Removing tickets from unknown topic","topic",topic)
}