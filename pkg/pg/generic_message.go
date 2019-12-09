package pg

// GenericMessage represents any PostgreSQL message that we don't want to process.
type GenericMessage struct {
	Type byte   // Message type
	Body []byte // Message raw bytes
}

// Compile time check to make sure that GenericMessage implements the Message interface.
var _ Message = &GenericMessage{}

// ParseGenericMessage parses raw network frame and returns GenericMessage.
func ParseGenericMessage(frame Frame) (*GenericMessage, error) {
	return &GenericMessage{
		Type: frame.MessageType(),
		Body: frame.MessageBody(),
	}, nil
}

// Frame serializes the message to send it over the network.
func (m *GenericMessage) Frame() (Frame, error) {
	// XXX: we assume that the frame is always standard
	return NewStandardFrame(m.Type, m.Body), nil
}
