package pg

// EmptyQueryResponseMessageType identifies EmptyQueryResponseMessage message.
const EmptyQueryResponseMessageType = 'I'

// EmptyQueryResponseMessage is sent by a backend as a response to an empty query string.
type EmptyQueryResponseMessage struct{}

// Compile time check to make sure that EmptyQueryResponseMessage implements the Message interface.
var _ Message = &EmptyQueryResponseMessage{}

// ParseEmptyQueryResponseMessage parses EmptyQueryResponseMessage from a network frame.
func ParseEmptyQueryResponseMessage(frame Frame) (*EmptyQueryResponseMessage, error) {
	// Assert the message type
	if frame.MessageType() != EmptyQueryResponseMessageType {
		return nil, ErrMalformedMessage
	}

	// Just in case assert that there is no message body
	if len(frame.MessageBody()) > 0 {
		return nil, ErrMalformedMessage
	}

	return &EmptyQueryResponseMessage{}, nil
}

// Frame serializes the message into a network frame.
func (m *EmptyQueryResponseMessage) Frame() Frame {
	return NewStandardFrame(EmptyQueryResponseMessageType, nil)
}
