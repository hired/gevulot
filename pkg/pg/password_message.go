package pg

// PasswordMessageType identifies PasswordMessage message.
const PasswordMessageType = 'p'

// PasswordMessage is sent by a frontend in response to authentication request.
type PasswordMessage struct {
	// The password (encrypted, if requested).
	Password string
}

// Compile time check to make sure that QueryMessage implements the Message interface.
var _ Message = &PasswordMessage{}

// PasswordMessage parses PasswordMessage from a network frame.
func ParsePasswordMessage(frame Frame) (*PasswordMessage, error) {
	// Assert the message type
	if frame.MessageType() != PasswordMessageType {
		return nil, ErrMalformedMessage
	}

	messageData := ReadBuffer(frame.MessageBody())

	password, err := messageData.ReadString()

	if err != nil {
		return nil, err
	}

	return &PasswordMessage{password}, nil
}

// Frame serializes the message into a network frame.
func (m *PasswordMessage) Frame() Frame {
	var messageBuffer WriteBuffer
	messageBuffer.WriteString(m.Password)

	return NewStandardFrame(PasswordMessageType, messageBuffer)
}
