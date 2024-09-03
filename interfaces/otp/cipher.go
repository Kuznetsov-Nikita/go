//go:build !solution

package otp

import (
	"io"
)

type CipherReader struct {
	input  io.Reader
	prng   io.Reader
	buffer []byte
}

func (r *CipherReader) Read(p []byte) (n int, err error) {
	r.buffer = make([]byte, len(p))
	cnt, err := r.input.Read(r.buffer)
	if err != nil {
		if cnt != 0 {
			prngByte := make([]byte, cnt)
			_, _ = r.prng.Read(prngByte)

			for i := 0; i < cnt; i++ {
				p[i] = r.buffer[i] ^ prngByte[i]
			}

			r.buffer = r.buffer[cnt:]

			return cnt, nil
		}
		return 0, err
	}

	prngByte := make([]byte, cnt)
	_, _ = r.prng.Read(prngByte)

	for i := 0; i < cnt; i++ {
		p[i] = r.buffer[i] ^ prngByte[i]
	}

	r.buffer = r.buffer[cnt:]

	return cnt, nil
}

func NewReader(r io.Reader, prng io.Reader) io.Reader {
	return &CipherReader{input: r, prng: prng}
}

type CipherWriter struct {
	output io.Writer
	prng   io.Reader
}

func (w *CipherWriter) Write(p []byte) (n int, err error) {
	data := make([]byte, len(p))

	prngByte := make([]byte, len(p))
	_, _ = w.prng.Read(prngByte)

	for i := 0; i < len(p); i++ {
		data[i] = p[i] ^ prngByte[i]
	}

	cnt, err := w.output.Write(data)
	if err != nil {
		return cnt, err
	}

	return len(p), nil
}

func NewWriter(w io.Writer, prng io.Reader) io.Writer {
	return &CipherWriter{output: w, prng: prng}
}
