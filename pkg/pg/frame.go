package pg

// A Frame is a raw PostgreSQL message sent/received over the network to/from client/DB.
// See https://www.postgresql.org/docs/9.6/protocol-overview.html#PROTOCOL-MESSAGE-CONCEPTS for details.
// All frame types must implement this interface.
type Frame interface {
	// MessageType returns this frame's message type.
	MessageType() byte

	// MessageBody returns message body bytes.
	MessageBody() []byte

	// Bytes returns raw bytes of the frame ready to be sent over the network.
	Bytes() []byte
}
