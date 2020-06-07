package packages

import (
	"bytes"
	"os"
	"sync"
)

type bufferPool struct {
	pool sync.Pool
}

func (b *bufferPool) Get() *bytes.Buffer {
	return b.pool.Get().(*bytes.Buffer)
}

func (b *bufferPool) Put(buf *bytes.Buffer) {
	buf.Reset()
	b.pool.Put(buf)
}

var memoryPageSize = os.Getpagesize()

var bufPool = &bufferPool{
	sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 0, memoryPageSize))
		},
	},
}
