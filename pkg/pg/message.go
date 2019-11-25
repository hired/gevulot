package pg

// Message is a generic PostgreSQL message that a client and a DB sending to each other.
type Message interface {
	// Marshall serializes the message to send it over the network.
	Marshall() ([]byte, error)
}
