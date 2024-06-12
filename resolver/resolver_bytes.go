package resolver

import (
	"bytes"
	"io"
	"time"
)

type resolverBytes struct {
	r Resolver[[]byte]
	b *bytes.Buffer
}

func NewBytes(timeout time.Duration) ResolverBytes {
	return &resolverBytes{
		r: NewSingle[[]byte](timeout),
		b: bytes.NewBuffer(nil),
	}
}

func (r *resolverBytes) Read(buf []byte) (int, error) {
	if r.b.Len() > 0 {
		return r.b.Read(buf)
	}

	var bytes, err = r.r.Read()
	if err != nil {
		if err == ErrClosed {
			return 0, io.EOF
		}
		return 0, err
	}

	r.b.Reset()
	r.b.Write(bytes)

	return r.b.Read(buf)
}

func (r *resolverBytes) Write(buf []byte) (int, error) {
	if err := r.r.Write(buf); err != nil {
		return 0, err
	}
	return len(buf), nil
}

func (r *resolverBytes) Close() error {
	return r.r.Close()
}