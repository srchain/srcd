package discover5

import "github.com/srchain/srcd/common/common"

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

