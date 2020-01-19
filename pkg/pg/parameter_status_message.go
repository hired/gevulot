package pg

// QueryMessageType identifies ParameterStatusMessage message.
const ParameterStatusMessageType = 'S'

// ParameterStatusMessage informs a frontend about the current (initial) setting of backend parameters.
type ParameterStatusMessage struct {
	// The name of the run-time parameter being reported.
	Name string

	// The current value of the parameter.
	Value string
}

// Compile time check to make sure that ParameterStatusMessage implements the Message interface.
var _ Message = &ParameterStatusMessage{}

// ParseParameterStatusMessage parses ParameterStatusMessage from a network frame.
func ParseParameterStatusMessage(frame Frame) (*ParameterStatusMessage, error) {
	// Assert the message type
	if frame.MessageType() != ParameterStatusMessageType {
		return nil, ErrMalformedMessage
	}

	messageData := ReadBuffer(frame.MessageBody())

	name, err := messageData.ReadString()

	if err != nil {
		return nil, err
	}

	value, err := messageData.ReadString()

	if err != nil {
		return nil, err
	}

	return &ParameterStatusMessage{name, value}, nil
}

// Frame serializes the message into a network frame.
func (m *ParameterStatusMessage) Frame() Frame {
	var messageBuffer WriteBuffer

	messageBuffer.WriteString(m.Name)
	messageBuffer.WriteString(m.Value)

	return NewStandardFrame(ParameterStatusMessageType, messageBuffer)
}
