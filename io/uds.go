package io

import (
	"errors"
	"net"
	"time"
)

var (
	reconnectInterval               = 5 * time.Second
	defaultUDSWriterBufferSize uint = 1 << 15
)

type UDSWriter struct {
	conn            *net.UnixConn
	address         string
	buffer          chan []byte
	latestReconnect time.Time
}

// NewUDSWriter returns the poniter of a new UDSWriter, which impl io.Writer interface.
//
// The parameter address is the absolute path to unix domain socket.
// The parameter bufSize specifies the maximum number of byte slices held by the buffer,
// and each Write will increases the number of byte slices in the current buffer by one.
//
// NOTE:
//   1. If a connection cannot be established using the parameter address , Write will
// degrade to something like writing to /dev/null.
//   2. If the buffer space eventually fills up due to slow data synchronization via
// socket, subsequent Write will return an error.
func NewUDSWriter(address string, bufSize uint) *UDSWriter {
	if bufSize == 0 {
		bufSize = defaultUDSWriterBufferSize
	}
	w := &UDSWriter{
		address:         address,
		buffer:          make(chan []byte, bufSize),
		latestReconnect: time.Now(),
	}
	go w.sync()
	return w
}

func (w *UDSWriter) Write(b []byte) (n int, err error) {
	select {
	case w.buffer <- append([]byte(nil), b...):
		return len(b), nil
	default:
		return 0, errors.New("the buffer is full")
	}
}

func (w *UDSWriter) sync() {
	// initializing the connection
	w.connect()
	for b := range w.buffer {
		var err error
		if w.conn == nil {
			err = errors.New("the conn is nil")
		} else {
			_, err = w.conn.Write(b)
		}
		w.reconnect(err)
	}
}

func (w *UDSWriter) connect() error {
	if w.address == "" {
		return errors.New("missing address of uds")
	}
	conn, err := net.DialUnix("unix", nil, &net.UnixAddr{Name: w.address, Net: "unix"})
	if err != nil {
		return err
	}
	w.conn = conn
	return nil

}

func (w *UDSWriter) reconnect(err error) {
	if err == nil || time.Since(w.latestReconnect) < reconnectInterval {
		return
	}
	w.latestReconnect = time.Now()
	w.connect()
}
