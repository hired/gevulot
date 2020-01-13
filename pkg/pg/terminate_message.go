package pg

// TerminateMessageType identifies TerminateMessage message.
const TerminateMessageType = 'X'

// TerminateMessage is sent by a frontend to terminate the session.
type TerminateMessage struct{}

// Compile time check to make sure that TerminateMessage implements the Message interface.
var _ Message = &TerminateMessage{}

// ParseQueryMessage parses QueryMessage from a network frame.
func ParseTerminateMessage(frame Frame) (*TerminateMessage, error) {
	// Assert the message type
	if frame.MessageType() != TerminateMessageType {
		return nil, ErrMalformedMessage
	}

	// Just in case assert that there is no message body
	if len(frame.MessageBody()) > 0 {
		return nil, ErrMalformedMessage
	}

	return &TerminateMessage{}, nil
}

// Frame serializes the message into a network frame.
func (m *TerminateMessage) Frame() Frame {
	return NewStandardFrame(TerminateMessageType, nil)
}
