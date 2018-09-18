package discover5

import (
	"github.com/srchain/srcd/common/common"
	"net"
	"go-ethereum/p2p/netutil"
	"time"
	"crypto/ecdsa"
	"github.com/srchain/srcd/log"
)


const (
	autoRefreshInterval   = 1 * time.Hour
	bucketRefreshInterval = 1 * time.Minute
	seedCount             = 30
	seedMaxAge            = 5 * 24 * time.Hour
	lowPort               = 1024
)


const (
	printTestImgLogs = false
)

var (
	unknown          *nodeState
	verifyinit       *nodeState
	verifywait       *nodeState
	remoteverifywait *nodeState
	known            *nodeState
	contested        *nodeState
	unresponsive     *nodeState
)


type nodeNetGuts struct {
	sha common.Hash
	state *nodeState
	pingEcho []byte
	pingTopics []Topic
	deferredQueries 	[]*findnodeQuery
	pendingNeighbours	*findnodeQuery
	queryTimeouts	int
}

type timeoutEvent struct {
	ev nodeEvent
	node *Node
}

type findnodeQuery struct {
	remote *Node
	target common.Hash
	reply chan<- []*Node
	nresults	int
}

type topicRegisterReq struct {
	add bool
	topic Topic
}

type topicSearchReq struct {
	topic Topic
	found chan<- *Node
	lookup chan<- bool
	delay time.Duration
}

type Network struct {
	db *nodeDB
	conn transport
	netrestrict *netutil.Netlist

	closed	chan struct{}
	closeReq	chan struct{}
	refreshReq	chan []*Node
	refreshRsp 	chan (<-chan struct{})
	read 	chan ingressPacket
	timeout chan timeoutEvent
	queryReq chan *findnodeQuery
	tableOpReq	chan func()
	tableOpResp	chan func()
	topicRegisterReq	chan topicRegisterReq
	topicSearchReq	chan topicSearchReq

	// State of the main loop
	tab *Table
	topictab *topicTable
	ticketStore	*ticketStore
	nursery []*Node
	nodes 	map[NodeID]*Node
	timeoutTimers map[timeoutEvent]*time.Timer

	slowRevalidateQueue []*Node
	fastRevalidateQueue []*Node

	sendBuf []*ingressPacket



}



type transport interface {
	sendPing(remote *Node, remoteAddr *net.UDPAddr, topics []Topic) (hash []byte)
	sendNeighbours(remote *Node, nodes []*Node)
	sendFindnodeHash(remote *Node, nodes []*Node)
	sendTopicRegister(remote *Node, topics []Topic, topicIdx int, pong []byte)
	sendTopicNodes(remote *Node, queryHash common.Hash, nodes []*Node)

	send(remote *Node, ptype nodeEvent, p interface{}) (hash []byte)

	localAddr() *net.UDPAddr
	Close()

}


type nodeState struct {
	name string
	handle func(*Network, *Node, nodeEvent, *ingressPacket) (next *nodeState, err error)
	enter func(*Network, *Node)
}

type nodeEvent uint


type topicSearchResult {
	target lookupInfo
	nodes []*Node
}



func newNetwork(conn transport, ourPubkey ecdsa.PublicKey, dbPath string , netrestrict *netutil.Netlist) (*Network, error) {
	ourID := PubKeyID(ourPubkey)
	var db *nodeDB
	if dbPath != "<no database>" {
		var err error
		if db, err = newNodeDB(dbPath,Version,ourID); err != nil {
			return nil, err
		}
	}

	tab := newTable(ourID,conn.localAddr())
	net := &Network{
		db:db,
		conn: conn,
		netrestrict: netrestrict,
		tab: tab,
		topictab: newTopicTable(db,tab.self),
		ticketStore: newTicketStore(),
		refreshReq: make(chan []*Node),
		refreshRsp: make(chan (<- chan struct{})),
		closed: make(chan struct{}),
		closeReq: make(chan struct{}),
		read: make(chan ingressPacket, 100),
		timeout: make(chan timeoutEvent),
		timeoutTimers: make(map[timeoutEvent]*time.Timer),
		tableOpReq: make(chan func()),
		tableOpResp: make(chan struct{}),
		queryReq: make(chan *findnodeQuery),
		topicRegisterReq: make(chan topicRegisterReq),
		topicSearchReq: make(chan topicSearchReq),
		nodes: make(map[NodeID]*Node),
	}
	go net.loop()
	return net , nil
}

func (net *Network) loop() {
	var (
		refreshTimer = time.NewTicker(autoRefreshInterval)
		bucketRefreshTimer = time.NewTimer(bucketRefreshInterval)
		refreshDone chan struct{}
	)

	var (
		nextTicket *ticketRef
		nextRegisterTimer *time.Timer
		nextRegisterTime <-chan time.Time
	)

	defer func () {
		if nextRegisterTimer != nil {
			nextRegisterTimer.Stop()
		}
	}()

	resetNextTicket := func() {
		ticket, timeout := net.ticketStore.nextFilteredTicket()
		if nextTicket != ticket {
			nextTicket = ticket
			if nextRegisterTimer != nil {
				nextRegisterTimer.Stop()
				nextRegisterTime = nil
			}
			if ticket != nil {
				nextRegisterTimer = time.NewTimer(timeout)
				nextRegisterTime = nextRegisterTimer.C
			}
		}
	}

	var (
		topicRegisterLookupTarget lookupInfo
		topicRegisterLookupDone chan []*Node
		topicRegisterLookupTick = time.NewTimer(0)
		searchReqWhenRefreshDone []topicSearchReq
		searchInfo = make(map[Topic]topicSearchInfo)
		activeSearchCount int
	)

	topicSearchLookupDone := make(chan topicSearchResult, 100)
	topicSearch := make(chan Topic, 100)
	<- topicRegisterLookupTick.C

	statsDump := time.NewTicker(10 * time.Second)

	loop:
		for {
			resetNextTicket()
			select {
				case <- net.closeReq:
					log.Trace("<- net.closeReq")
					break loop
				case pkt := <-net.read:
					log.Trace("<- net.read")
				n := net.internNode(&pkt)
				prestate := n.state
				status := "ok"
				if err := net.hand

			}
		}
}

func (net *Network) internNode(pkt *ingressPacket) *Node {
	if n := net.nodes[pkt.remoteID]; n != nil {
		n.IP = pkt.remoteAddr.IP
		n.UDP = uint16(pkt.remoteAddr.Port)
		n.TCP = uint16(pkt.remoteAddr.Port)
		return n
	}
	n := NewNode(pkt.remoteID,pkt.remoteAddr.IP,uint16(pkt.remoteAddr.Port),uint16(pkt.remoteAddr.Port))
	n.state = unknown
	net.nodes[n.ID] = n
	return n

}