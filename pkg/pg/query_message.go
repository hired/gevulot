package pg

// QueryMessageType identifies QueryMessage message.
const QueryMessageType = 'Q'

// QueryMessage represent a simple SQL query sent by a frontend.
type QueryMessage struct {
	// SQL query
	Query string
}

// Compile time check to make sure that QueryMessage implements the Message interface.
var _ Message = &QueryMessage{}

// ParseQueryMessage parses QueryMessage from the network frame.
func ParseQueryMessage(frame Frame) (*QueryMessage, error) {
	// Assert the message type
	if frame.MessageType() != QueryMessageType {
		return nil, ErrMalformedMessage
	}

	messageData := ReadBuffer(frame.MessageBody())

	query, err := messageData.ReadString()

	if err != nil {
		return nil, err
	}

	return &QueryMessage{Query: query}, nil
}

// Frame serializes the message into a network frame.
func (m *QueryMessage) Frame() Frame {
	var messageBuffer WriteBuffer

	messageBuffer.WriteString(m.Query)

	return NewStandardFrame(QueryMessageType, messageBuffer)
}
