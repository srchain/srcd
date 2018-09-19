package discover5

import (
	"net"
	"github.com/srchain/srcd/common/common"
	"github.com/srchain/srcd/rlp"
	"time"
)


const Version = 4


const (
	respTimeout = 500 * time.Millisecond
	expiration  = 20 * time.Second

	driftThreshold = 10 * time.Second // Allowed clock drift before warning user
)


type (
	rpcEndpoint struct {
		IP net.IP
		UDP uint16
		TCP uint16
	}

	pong struct {
		To rpcEndpoint
		ReplyTok []byte
		Expiration uint64
		TopicHash common.Hash
		TicketSerial uint32
		WaitPeriods []uint32
		Rest []rlp.RawValue `rlp:"tail"`
	}
)


type ingressPacket struct {
	remoteID 	NodeID
	remoteAddr	*net.UDPAddr
	ev 	nodeEvent
	hash 	[]byte
	data 	interface{}
	rawData	[]byte
}

type conn interface {
	 ReadFromUDP(b []byte) (n int, addr *net.UDPAddr, err error)
	 WriteToUDP(b []byte, addr *net.UDPAddr) (n int, err error)
	 Close() error
	 LocalAddr() net.Addr
}