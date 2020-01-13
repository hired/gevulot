package pg

// GenericMessage represents any PostgreSQL message that we don't want to process.
type GenericMessage struct {
	Type byte   // Message type
	Body []byte // Message raw bytes
}

// Compile time check to make sure that GenericMessage implements the Message interface.
var _ Message = &GenericMessage{}

// ParseGenericMessage parses GenericMessage from the network frame.
func ParseGenericMessage(frame Frame) (*GenericMessage, error) {
	return &GenericMessage{
		Type: frame.MessageType(),
		Body: frame.MessageBody(),
	}, nil
}

// Frame serializes the message into a network frame.
func (m *GenericMessage) Frame() Frame {
	// XXX: we assume that the frame is always standard
	return NewStandardFrame(m.Type, m.Body)
}
