package vm

import (
	"encoding/binary"
)

type  Builder struct {
	program     []byte
	jumpCounter int

	// Maps a jump target number to its absolute address.
	jumpAddr map[int]uint32

	// Maps a jump target number to the list of places where its
	// absolute address must be filled in once known.
	jumpPlaceholders map[int][]int
}

func NewBuilder() *Builder {
	return &Builder{
		jumpAddr:         make(map[int]uint32),
		jumpPlaceholders: make(map[int][]int),
	}
}

func (b *Builder) Build() ([]byte, error) {
	for target, placeholders := range b.jumpPlaceholders {
		addr, ok := b.jumpAddr[target]
		if !ok {
			return nil, nil
		}
		for _, placeholder := range placeholders {
			binary.LittleEndian.PutUint32(b.program[placeholder:placeholder+4], addr)
		}
	}
	return b.program, nil
}


// AddInt64 adds a pushdata instruction for an integer value.
func (b *Builder) AddInt64(n int64) *Builder {
	b.program = append(b.program, PushdataInt64(n)...)
	return b
}
// AddData adds a pushdata instruction for a given byte string.
func (b *Builder) AddData(data []byte) *Builder {
	b.program = append(b.program, PushdataBytes(data)...)
	return b
}

