package discover5

import (
	"go-ethereum/common/mclock"
	"time"
)

type Topic string

type topicTable struct {
	db *nodeDB
	self	*Node
	nodes 	map[*Node]*nodeInfo
	topics 	map[Topic]*topicInfo
}


type waitControlLoop struct {
	lastIncoming mclock.AbsTime
	waitPerid	time.Duration
}

type nodeInfo struct {
	entries 	map[Topic]*topicEntry
	lastIssuedTicket, lastUsedTicket uint32
	noRegUntil mclock.AbsTime
}

type topicEntry struct {
	topic Topic
	fifoIdx	uint64
	node *Node
	expire mclock.AbsTime
}

type topicInfo struct {
	entries map[uint64]*topicEntry
	fifoHead, fifoTail uint64
	rqItem	*topicRequestQueueItem
	wcl	waitControlLoop
}

type topicRequestQueueItem struct {
	topic	Topic
	priority uint64
	index int
}

type topicRequestQueue []*topicRequestQueueItem



