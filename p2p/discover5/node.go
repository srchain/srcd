package discover5

import "net"

const nodeIDBits  = 512

type NodeID [nodeIDBits / 8]byte

type Node struct {
	IP       net.IP // len 4 for IPv4 or 16 for IPv6
	UDP, TCP uint16 // port numbers
	ID       NodeID // the node's public key

	// Network-related fields are contained in nodeNetGuts.
	// These fields are not supposed to be used off the
	// Network.loop goroutine.
	nodeNetGuts
}
