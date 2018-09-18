package discover5

import (
	"github.com/srchain/srcd/common/common"
	"net"
)

const (
	alpha      = 3  // Kademlia concurrency factor
	bucketSize = 16 // Kademlia bucket size
	hashBits   = len(common.Hash{}) * 8
	nBuckets   = hashBits + 1 // Number of buckets

	maxFindnodeFailures = 5
)

type Table struct {
	count int
	buckets [nBuckets]*bucket
	nodeAddedHook func(*Node)
	self	*Node
}

type bucket struct {
	entries      []*Node
	replacements []*Node
}


func newTable(ourID NodeID, ourAddr *net.UDPAddr) *Table {
	self := NewNode(ourID,ourAddr.IP,uint16(ourAddr.Port),uint16(ourAddr.Port))
	tab := &Table{self: self}
	for i := range tab.buckets {
		tab.buckets[i] = new(bucket)

	}
	return tab
}
