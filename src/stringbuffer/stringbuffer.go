package stringbuffer

import (
	"io"
)

// Write writes len(p) bytes from p to the underlying data stream. It returns the
// number of bytes written from p (0 <= n <= len(p)) and any error encountered
// that caused the write to stop early. Write must return a non-nil error if it
// returns n < len(p). Write must not modify the slice data, even temporarily.

type Buffer struct {
    buffer []byte
}

func (p *Buffer) Write(b []byte) (int, error) {
    p.buffer = append(p.buffer, b...)
    return len(b), nil
}

func (p *Buffer) String() (string) {
    return string(p.buffer)
}

// Read reads up to len(dst) bytes into dst. It returns the number of bytes read
// (0 <= n <= len(dst)) and any error encountered. Even if Read returns n < len(dst),
// it may  use all of dst as scratch space during the call. If some data is available
// but not len(dst) bytes, Read conventionally returns what is available instead of
// waiting for more. 

func (p *Buffer) Read(dst []byte) (int, error) {
	max := len(dst)
	n := 0
	for idx, b := range p.buffer {
		if idx < max {
			dst[idx] = b
			n++
		} else {
			break
		}
	}
	p.buffer = p.buffer[n:]
	var err error = nil
	if len(p.buffer) == 0 {
		err = io.EOF
	}
	return n, err
}

func (p *Buffer) Store(str string) {
	p.buffer = []byte(str)
}
