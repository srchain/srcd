package rawdb

import (
	"srcd/common/misc"
)

// ReadCanonicalHash retrieves the hash assigned to a canonical block number.
func ReadCanonicalHash(db DatabaseReader, number uint64) misc.Hash {
	data, _ := db.Get(headerHashKey(number))
	if len(data) == 0 {
		return misc.Hash{}
	}
	return misc.BytesToHash(data)
}
