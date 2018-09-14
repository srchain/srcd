package discover5

import "net"


const Version = 4

type (
	rpcEndpoint struct {
		IP net.IP
		UDP uint16
		TCP uint16
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