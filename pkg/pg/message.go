package pg

import (
	"errors"
)

// Message is a PostgreSQL message that a client and a DB sending to each other.
type Message interface {
	// Frame serializes the message into a network frame.
	Frame() Frame
}

var (
	// ErrMalformedMessage is returned when message cannot be parsed.
	ErrMalformedMessage = errors.New("pg: malformed message")
)
