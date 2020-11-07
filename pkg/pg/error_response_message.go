//nolint:dupl
package pg

// ErrorResponseMessageType identifies ErrorResponseMessageType message.
const ErrorResponseMessageType = 'E'

// ErrorResponseMessage is sent by a backend when an error occurs.
type ErrorResponseMessage struct {
	// One of more fields with error info
	Fields []*MessageField
}

// Compile time check to make sure that ErrorResponseMessage implements the Message interface.
var _ Message = &ErrorResponseMessage{}

// ParseErrorResponseMessage parses ErrorResponseMessage from a network frame.
func ParseErrorResponseMessage(frame Frame) (*ErrorResponseMessage, error) {
	// Assert the message type
	if frame.MessageType() != ErrorResponseMessageType {
		return nil, ErrMalformedMessage
	}

	messageData := ReadBuffer(frame.MessageBody())
	fields := []*MessageField{}

	// Read all the fields
	for {
		fieldType, err := messageData.ReadByte()

		if err != nil {
			return nil, err
		}

		// Terminator — end of message
		if fieldType == 0 {
			break
		}

		fieldValue, err := messageData.ReadString()

		if err != nil {
			return nil, err
		}

		fields = append(fields, &MessageField{
			Type:  MessageFieldType(fieldType),
			Value: fieldValue,
		})
	}

	return &ErrorResponseMessage{fields}, nil
}

// Frame serializes the message into a network frame.
func (m *ErrorResponseMessage) Frame() Frame {
	var messageBuffer WriteBuffer

	for _, field := range m.Fields {
		messageBuffer.WriteByte(byte(field.Type))
		messageBuffer.WriteString(field.Value)
	}

	// Terminator
	messageBuffer.WriteByte(0)

	return NewStandardFrame(ErrorResponseMessageType, messageBuffer)
}
