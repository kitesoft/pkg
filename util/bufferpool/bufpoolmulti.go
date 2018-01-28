// Copyright 2015-2016, Cyrill @ Schumacher.fm and the CoreStore contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bufferpool

import (
	"bytes"
	"io"
	"sync"
)

var twinBufferPool = NewTwin(256) // estimated *cough* average size

type twinBuffer struct {
	First  *bytes.Buffer
	Second *bytes.Buffer
}

func (tw twinBuffer) Write(p []byte) (n int, err error) {
	n, err = tw.First.Write(p)
	if err != nil {
		return
	}
	if n != len(p) {
		return 0, io.ErrShortWrite
	}
	n, err = tw.Second.Write(p)
	return n, err
}

// GetTwin returns a buffer containing two buffers, `First` and `Second`, from the pool.
func GetTwin() *twinBuffer {
	return twinBufferPool.Get()
}

// PutTwin returns a twin buffer to the pool. The buffers get reset before they
// are put back into circulation.
func PutTwin(buf *twinBuffer) {
	twinBufferPool.Put(buf)
}

// PutTwinCallBack same as PutTwin but executes fn after buf has been returned
// into the pool.
//		buf := twinBuf.Get()
//		defer twinBuf.PutCallBack(buf, wg.Done)
func PutTwinCallBack(buf *twinBuffer, fn func()) {
	twinBufferPool.PutCallBack(buf, fn)
}

// twinTank implements a sync.Pool for twinBuffer
type twinTank struct {
	p *sync.Pool
}

// Get returns type safe a buffer
func (t twinTank) Get() *twinBuffer {
	return t.p.Get().(*twinBuffer)
}

// Put empties the buffer and returns it back to the pool.
//
//		bp := NewTwin(512)
//		buf := bp.Get()
//		defer bp.Put(buf)
//		// your code
//		return buf.String()
//
// If you use Bytes() function to return bytes make sure you copy the data
// away otherwise your returned byte slice will be empty.
// For using String() no copying is required.
func (t twinTank) Put(buf *twinBuffer) {
	buf.First.Reset()
	buf.Second.Reset()
	t.p.Put(buf)
}

// PutCallBack same as Put but executes fn after buf has been returned into the
// pool. Good use case when you might have multiple defers in your code.
func (t twinTank) PutCallBack(buf *twinBuffer, fn func()) {
	buf.First.Reset()
	buf.Second.Reset()
	t.p.Put(buf)
	fn()
}

// NewTwin instantiates a new twinBuffer pool with a custom pre-allocated
// buffer size. The fields `First` and `Second` will have the same size.
func NewTwin(size int) twinTank {
	return twinTank{
		p: &sync.Pool{
			New: func() interface{} {
				return &twinBuffer{
					First:  bytes.NewBuffer(make([]byte, 0, size)),
					Second: bytes.NewBuffer(make([]byte, 0, size)),
				}
			},
		},
	}
}