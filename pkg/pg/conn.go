package pg

import (
	"net"
)

// Conn is a wrapper around net.Conn that allows us to send/receive PostgreSQL messages.
type Conn struct {
	conn net.Conn
}

// NewConn initializes new pg.Conn.
func NewConn(conn net.Conn) *Conn {
	return &Conn{conn: conn}
}

// RecvStartupMessage receives StartupMessage from the underlying network connection.
func (h *Conn) RecvStartupMessage() (*StartupMessage, error) {
	frame, err := ReadStartupFrame(h.conn)

	if err != nil {
		return nil, err
	}

	return ParseStartupMessage(frame)
}

// RecvMessage receives PostgreSQL message from the underlying network connection.
func (h *Conn) RecvMessage() (Message, error) {
	frame, err := ReadStandardFrame(h.conn)

	if err != nil {
		return nil, err
	}

	switch frame.MessageType() {
	default:
		return ParseGenericMessage(frame)
	}
}

// SendMessage sends given message over the network.
func (h *Conn) SendMessage(msg Message) error {
	_, err := h.conn.Write(msg.Frame().Bytes())
	return err
}

// SendByte sends given byte over the network.
func (h *Conn) SendByte(c byte) error {
	_, err := h.conn.Write([]byte{c})
	return err
}

// Close closes the underlying network connection.
func (h *Conn) Close() error {
	return h.conn.Close()
}
