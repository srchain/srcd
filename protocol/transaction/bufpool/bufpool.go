package bufpool

import (
	"sync"
	"bytes"
)

var pool = &sync.Pool{New: func() interface{} { return bytes.NewBuffer(nil) }}

func Get() *bytes.Buffer {
	return pool.Get().(*bytes.Buffer)
}

func Put(b *bytes.Buffer) {
	b.Reset()
	pool.Put(b)
}
