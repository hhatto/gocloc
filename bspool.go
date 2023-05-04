package gocloc

import (
	"bytes"
	"sync"
)

var bsPool = sync.Pool{New: func() any { return new(bytes.Buffer) }}

func getByteSlice() *bytes.Buffer {
	v := bsPool.Get().(*bytes.Buffer)
	v.Reset()
	return v
}

func putByteSlice(bs *bytes.Buffer) {
	bsPool.Put(bs)
}
