package pg

// CommandCompleteMessageType identifies CommandCompleteMessage message.
const CommandCompleteMessageType = 'C'

// CommandCompleteMessage sent by a backend to notify frontend about successful command execution.
type CommandCompleteMessage struct {
	// The command tag. This is usually a single word that identifies which SQL command was completed.
	Tag string
}

// Compile time check to make sure that CommandCompleteMessage implements the Message interface.
var _ Message = &CommandCompleteMessage{}

// ParseCommandCompleteMessage parses CommandCompleteMessage from a network frame.
func ParseCommandCompleteMessage(frame Frame) (*CommandCompleteMessage, error) {
	// Assert the message type
	if frame.MessageType() != CommandCompleteMessageType {
		return nil, ErrMalformedMessage
	}

	messageData := ReadBuffer(frame.MessageBody())

	tag, err := messageData.ReadString()

	if err != nil {
		return nil, err
	}

	return &CommandCompleteMessage{Tag: tag}, nil
}

// Frame serializes the message into a network frame.
func (m *CommandCompleteMessage) Frame() Frame {
	var messageBuffer WriteBuffer

	messageBuffer.WriteString(m.Tag)

	return NewStandardFrame(CommandCompleteMessageType, messageBuffer)
}
