package gstauto

import (
	"bufio"
	"io"
	"os"
)

// Blank assertions to ensure interfaces are implemented.
var _ io.ReadCloser = &readCloser{}
var _ io.WriteCloser = &writeCloser{}
var _ io.ReadWriteCloser = &readWriteCloser{}

// readCloser is a struct that provides a read buffer that can also be written to
// internally.
type readCloser struct {
	rReader, rWriter *os.File
	rBuf             *bufio.Reader
}

// newReadCloser returns a new readCloser.
func newReadCloser() (*readCloser, error) {
	r, w, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	return &readCloser{
		rReader: r,
		rWriter: w,
		rBuf:    bufio.NewReader(r),
	}, nil
}

// Read implements a Reader for objects embdedding this struct.
func (r *readCloser) Read(p []byte) (int, error) { return r.rBuf.Read(p) }

// Close implements a Closer for objects embedding this struct.
func (r *readCloser) Close() error {
	if err := r.rWriter.Close(); err != nil {
		return err
	}
	return r.rReader.Close()
}

// writeCloser is a struct that provides a read buffer that can also be
// read from internally.
type writeCloser struct {
	wReader, wWriter *os.File
	wBuf             *bufio.Writer
}

// newWriteCloser returns a new writeCloser.
func newWriteCloser() (*writeCloser, error) {
	r, w, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	return &writeCloser{
		wReader: r,
		wWriter: w,
		wBuf:    bufio.NewWriter(w),
	}, nil
}

// Write implements a Writer for objects embedding this struct.
func (w *writeCloser) Write(p []byte) (int, error) { return w.wBuf.Write(p) }

// Close implements a Closer for objects embedding this struct.
func (w *writeCloser) Close() error {
	if err := w.wWriter.Close(); err != nil {
		return err
	}
	return w.wReader.Close()
}

// readWriteCloser is a struct that provides both read and write buffers.
type readWriteCloser struct {
	*readCloser
	*writeCloser
}

// newReadWriteCloser returns a new readWriteCloser.
func newReadWriteCloser() (*readWriteCloser, error) {
	rCloser, err := newReadCloser()
	if err != nil {
		return nil, err
	}
	wCloser, err := newWriteCloser()
	if err != nil {
		return nil, err
	}
	return &readWriteCloser{readCloser: rCloser, writeCloser: wCloser}, nil
}

// Close implements a Closer for objects embedding this struct.
func (rw *readWriteCloser) Close() error {
	if err := rw.writeCloser.Close(); err != nil {
		return err
	}
	return rw.readCloser.Close()
}
