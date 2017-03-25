package mcquery

import (
	"net"
	"time"
)

type TimeoutConn struct {
	conn    net.Conn
	timeout time.Duration
}

func (i TimeoutConn) Read(buf []byte) (int, error) {
	if err := i.SetDeadline(time.Now().Add(i.timeout)); err != nil {
		return 0, err
	}
	return i.conn.Read(buf)
}

func (i TimeoutConn) Write(buf []byte) (int, error) {
	if err := i.SetDeadline(time.Now().Add(i.timeout)); err != nil {
		return 0, err
	}
	return i.conn.Write(buf)
}

func (i TimeoutConn) Close() error {
	return i.conn.Close()
}

func (i TimeoutConn) LocalAddr() net.Addr {
	return i.conn.LocalAddr()
}

func (i TimeoutConn) RemoteAddr() net.Addr {
	return i.conn.RemoteAddr()
}

func (i TimeoutConn) SetDeadline(t time.Time) error {
	return i.conn.SetDeadline(t)
}

func (i TimeoutConn) SetReadDeadline(t time.Time) error {
	return i.conn.SetReadDeadline(t)
}

func (i TimeoutConn) SetWriteDeadline(t time.Time) error {
	return i.conn.SetWriteDeadline(t)
}
