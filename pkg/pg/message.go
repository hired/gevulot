package pg

// Message is a generic PostgreSQL message that a client and a DB sending to each other.
type Message interface {
	// Marshal serializes the message to send it over the network.
	Marshal() ([]byte, error)
}
