package io

import (
	"io"

	"go.uber.org/multierr"
)

type continuousMultiWriter []io.Writer

// MultiWriter creates a writer that duplicates its writes to all the
// provided writers, similar to the io.MultiWriter method in the
// standard library, with the difference that ContinuousMultiWriter
// continues to write the writers that have not yet been written to
// when it encounters an error.
//
// NOTE:
// The act of continuously writing to all io.Writer means that multiple
// write errors may be encountered at the same time, and the n returned
// by ContinuousMultiWriter will be set to 0 when an error occurs, even
// when an error of the type io.ErrShortWrite is encountered.
func NewContinuousMultiWriter(writers ...io.Writer) io.Writer {
	switch len(writers) {
	case 0:
		return new(continuousMultiWriter)
	case 1:
		return writers[0]
	default:
		return continuousMultiWriter(writers)
	}
}

func (t continuousMultiWriter) Write(p []byte) (n int, err error) {
	for _, w := range t {
		err = multierr.Append(err, doWrite(w, p))
	}
	if err == nil {
		n = len(p)
	}
	return
}

func doWrite(w io.Writer, p []byte) error {
	n, err := w.Write(p)
	if err != nil {
		return err
	}
	if n != len(p) {
		return io.ErrShortWrite
	}
	return nil
}
