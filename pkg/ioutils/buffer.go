package ioutils

import (
	"errors"
	"io"
)

var (
	errBufferFull             = errors.New("buffer is full")
	errBufferDataWillTruncate = errors.New("data wiill be truncated")
)

type fixedBuffer struct {
	buf      []byte
	pos      int
	lastRead int
}

func (b *fixedBuffer) Write(p []byte) (int, error) {
	n := copy(b.buf[b.pos:cap(b.buf)], p)
	b.pos += n

	if n < len(p) {
		if b.pos == cap(b.buf) {
			return n, errBufferFull
		}
		return n, io.ErrShortWrite
	}
	return n, nil
}

func (b *fixedBuffer) Read(p []byte) (int, error) {
	n := copy(p, b.buf[b.lastRead:b.pos])
	b.lastRead += n
	return n, nil
}

func (b *fixedBuffer) Len() int {
	return b.pos - b.lastRead
}

func (b *fixedBuffer) Cap() int {
	return cap(b.buf)
}

func (b *fixedBuffer) Reset() {
	b.pos = 0
	b.lastRead = 0
	b.buf = b.buf[:0]
}

func (b *fixedBuffer) String() string {
	return string(b.buf[b.lastRead:b.pos])
}

func (b *fixedBuffer) InitBuf(n int) {
	b.buf = make([]byte, n, n)
	b.lastRead = 0
	b.pos = 0
}

// 整理空间
func (b *fixedBuffer) Arrange() {
	if b.lastRead != 0 {
		copy(b.buf[0:cap(b.buf)], b.buf[b.lastRead:b.pos])
		b.pos = b.pos - b.lastRead
		b.lastRead = 0
	}
}

func (b *fixedBuffer) ReSieze(n int) error {
	if n == cap(b.buf) {
		return nil
	}
	if n < b.Len() {
		return errBufferDataWillTruncate
	}
	nbuf := make([]byte, 0, n)
	copy(nbuf[0:n], b.buf[b.lastRead:b.pos])

	b.pos = b.Len()
	b.lastRead = 0
	b.buf = nbuf
	return nil
}

func (b *fixedBuffer) Slice() []byte {
	return b.buf[b.lastRead:b.pos]
}
