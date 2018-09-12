package discover5

import (
	"github.com/srchain/srcd/common/common"
	"net"
	"go-ethereum/p2p/netutil"
	"time"
)

type nodeNetGuts struct {
	sha common.Hash
	state *nodeState
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

type Network struct {
	db *nodeDB
	conn transport
	nettrestict *netutil.Netlist

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
	topciSearchReq	chan topicSearchReq

	// State of the main loop
	tab *Table
	topictab *topicTable
	ticketStore	*ticketStore



}

type nodeState struct {
	name string
	handle func(*Network, *Node, nodeEvent, *ingressPacket) (next *nodeState, err error)
	enter func(*Network, *Node)
}

type nodeEvent uint