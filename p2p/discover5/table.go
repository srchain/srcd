package discover5

import (
	"github.com/srchain/srcd/common/common"
	"net"
	"sort"
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

type nodesByDistance struct {
	entries []*Node
	target common.Hash
}

func (h *nodesByDistance) push(node *Node, maxElems int) {
	ix := sort.Search(len(h.entries), func (i int) bool {
		return discmp(h.target,h.entries[i].sha,node.sha) > 0
	})
	if len(h.entries) < maxElems {
		h.entries = append(h.entries,node)
	}
	if ix == len(h.entries) {

	} else {
		copy(h.entries[ix+1:],h.entries[ix:])
		h.entries[ix] = node
	}

}

func (tab *Table) closest(target common.Hash, nresults int) *nodesByDistance {
	close := &nodesByDistance{target:target}
	for _, b := range &tab.buckets {
		for _, n := range b.entries {
			close.push(n, nresults)
		}
	}
}

func newTable(ourID NodeID, ourAddr *net.UDPAddr) *Table {
	self := NewNode(ourID,ourAddr.IP,uint16(ourAddr.Port),uint16(ourAddr.Port))
	tab := &Table{self: self}
	for i := range tab.buckets {
		tab.buckets[i] = new(bucket)

	}
	return tab
}
