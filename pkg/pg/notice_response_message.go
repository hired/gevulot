package pg

// NoticeResponseMessageType identifies NoticeResponseMessageType message.
const NoticeResponseMessageType = 'N'

// NoticeResponseMessage is sent by a backend when some kind of non-critical warning occurs.
type NoticeResponseMessage struct {
	// One of more fields with notice info
	Fields []*MessageField
}

// Compile time check to make sure that NoticeResponseMessage implements the Message interface.
var _ Message = &NoticeResponseMessage{}

// PnoticerrorResponseMessage parses NoticeResponseMessage from a network frame
func ParseNoticeResponseMessage(frame Frame) (*NoticeResponseMessage, error) {
	// Assert the message type
	if frame.MessageType() != NoticeResponseMessageType {
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

	return &NoticeResponseMessage{fields}, nil
}

// Frame serializes the message into a network frame.
func (m *NoticeResponseMessage) Frame() Frame {
	var messageBuffer WriteBuffer

	for _, field := range m.Fields {
		messageBuffer.WriteByte(byte(field.Type))
		messageBuffer.WriteString(field.Value)
	}

	// Terminator
	messageBuffer.WriteByte(0)

	return NewStandardFrame(NoticeResponseMessageType, messageBuffer)
}
