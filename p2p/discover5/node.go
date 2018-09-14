package discover5

import (
	"net"
	"encoding/hex"
	"strings"
	"fmt"
	"github.com/srchain/srcd/crypto/crypto"
)

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

func NewNode(id NodeID, ip net.IP, udpPort, tcpPort uint16) *Node {
	if ipv4 := ip.To4(); ipv4 != nil {
		ip = ipv4
	}
	return &Node {
		IP: ip,
		UDP: udpPort,
		TCP: tcpPort,
		ID:	id,
		nodeNetGuts: nodeNetGuts{sha:crypto.Keccak256Hash(id[:])},
	}
}

func MustHexID(in string) NodeID {
	id, err := HexID(in)
	if err != nil {
		panic(err)
	}
	return id
}

func HexID(in string) (NodeID, error) {
	var id NodeID
	b, err := hex.DecodeString(strings.TrimPrefix(in,"0x"))
	if err != nil {
		return id, err
	} else if len(b) != len(id) {
		return id, fmt.Errorf("wrong length, want %d hex chars",len(id) * 2)
	}
	copy(id[:],b)
	return id, nil
}
