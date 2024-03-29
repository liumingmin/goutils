package httpdownloader

import (
	"errors"
	"os"
	"sync"
)

// HttpFileOffsetWriter is a writer that writes data to a file at a given file offset.
type HttpFileOffsetWriter struct {
	writer       *os.File
	mutex        *sync.Mutex
	offset       int64
	end          int64
	originOffset int64
}

// NewOffsetWriter creates a new FileOffsetWriter instance to write at the given offset to the file.
func NewOffsetWriter(file *os.File, mutex *sync.Mutex, offset, end int64) *HttpFileOffsetWriter {
	return &HttpFileOffsetWriter{
		writer:       file,
		mutex:        mutex,
		offset:       offset,
		end:          end,
		originOffset: offset,
	}
}

// Write writes data to the file at the given offset.
func (w *HttpFileOffsetWriter) Write(p []byte) (n int, err error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	_, err = w.writer.Seek(w.offset, 0)
	if err != nil {
		return 0, err
	}

	n, err = w.writer.Write(p)
	if err != nil {
		return 0, err
	}

	// Increment the offset by the number of bytes written.
	w.offset += int64(n)

	if w.offset > w.end {
		return 0, errors.New("write out of range")
	}

	return n, nil
}

func (w *HttpFileOffsetWriter) ResetOffset() {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.offset = w.originOffset
}
