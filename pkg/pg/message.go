package pg

// Message is a PostgreSQL message that a client and a DB sending to each other.
type Message interface {
	// Frame returns network frame of the message.
	Frame() (Frame, error)
}
