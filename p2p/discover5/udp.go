package discover5

import "net"

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