package algorithm

import (
	"errors"
	"io"
)

type XORReader struct {
	src    io.Reader
	key    []byte
	pIndex *uint64

	keyLen uint64
}

func NewXORReader(src io.Reader, key []byte) io.Reader {
	readerIndex := uint64(0)
	return NewXORReaderWithOffset(src, key, &readerIndex)
}

func NewXORReaderWithOffset(src io.Reader, key []byte, pIndex *uint64) io.Reader {
	return &XORReader{
		src:    src,
		key:    key,
		pIndex: pIndex,
		keyLen: uint64(len(key)),
	}
}

func (r *XORReader) Read(p []byte) (int, error) {
	n, err := r.src.Read(p)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return n, err
	}

	for i := 0; i < n; i++ {
		p[i] = p[i] ^ r.key[((*r.pIndex)+uint64(i))%r.keyLen]
	}

	*r.pIndex = *r.pIndex + uint64(n)
	return n, err
}

type XORWriter struct {
	dst    io.Writer
	key    []byte
	pIndex *uint64

	keyLen uint64
}

func NewXORWriter(dst io.Writer, key []byte) io.Writer {
	writerIndex := uint64(0)
	return NewXORWriterWithOffset(dst, key, &writerIndex)
}

func NewXORWriterWithOffset(dst io.Writer, key []byte, pIndex *uint64) io.Writer {
	return &XORWriter{
		dst:    dst,
		key:    key,
		pIndex: pIndex,
		keyLen: uint64(len(key)),
	}
}

func (w *XORWriter) Write(p []byte) (int, error) {
	xorp := make([]byte, len(p))
	n := copy(xorp, p)
	if n != len(p) {
		return 0, errors.New("copy failed")
	}

	for i := 0; i < n; i++ {
		xorp[i] = xorp[i] ^ w.key[((*w.pIndex)+uint64(i))%w.keyLen]
	}

	writeN, err := w.dst.Write(xorp)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return writeN, err
	}

	*w.pIndex = *w.pIndex + uint64(writeN)
	return writeN, err
}

type XORReaderAt struct {
	src    io.ReaderAt
	key    []byte
	keyLen uint64
}

func NewXORReaderAt(src io.ReaderAt, key []byte) io.ReaderAt {
	return &XORReaderAt{
		src:    src,
		key:    key,
		keyLen: uint64(len(key)),
	}
}

func (r *XORReaderAt) ReadAt(p []byte, off int64) (int, error) {
	n, err := r.src.ReadAt(p, off)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return n, err
	}

	for i := 0; i < n; i++ {
		p[i] = p[i] ^ r.key[(uint64(off)+uint64(i))%r.keyLen]
	}

	return n, err
}
