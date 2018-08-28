package bc

import (
	"io"
	"github.com/golang/protobuf/proto"

)

// Entry is the interface implemented by each addressable unit in a
// blockchain: transaction components such as spends, issuances,
// outputs, and retirements (among others), plus blockheaders.
type Entry interface {
	proto.Message

	// type produces a short human-readable string uniquely identifying
	// the type of this entry.
	typ() string

	// writeForHash writes the entry's body for hashing.
	writeForHash(w io.Writer)
}

