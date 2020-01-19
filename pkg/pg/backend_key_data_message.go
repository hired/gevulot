package pg

// BackendKeyDataMessageType identifies BackendKeyDataMessage message.
const BackendKeyDataMessageType = 'K'

// BackendKeyDataMessage is sent by a backend to tell frontend about secret-key data.
// The frontend must save these values if it wishes to be able to issue CancelRequestMessage later.
type BackendKeyDataMessage struct {
	// The process ID of this backend.
	ProcessID int32

	// The secret key of this backend.
	Key int32
}

// Compile time check to make sure that BackendKeyDataMessage implements the Message interface.
var _ Message = &BackendKeyDataMessage{}

// ParseBackendKeyDataMessage parses BackendKeyDataMessage from a network frame.
func ParseBackendKeyDataMessage(frame Frame) (*BackendKeyDataMessage, error) {
	// Assert the message type
	if frame.MessageType() != BackendKeyDataMessageType {
		return nil, ErrMalformedMessage
	}

	messageData := ReadBuffer(frame.MessageBody())

	processID, err := messageData.ReadInt32()

	if err != nil {
		return nil, err
	}

	key, err := messageData.ReadInt32()

	if err != nil {
		return nil, err
	}

	return &BackendKeyDataMessage{ProcessID: processID, Key: key}, nil
}

// Frame serializes the message into a network frame.
func (m *BackendKeyDataMessage) Frame() Frame {
	var messageBuffer WriteBuffer

	messageBuffer.WriteInt32(m.ProcessID)
	messageBuffer.WriteInt32(m.Key)

	return NewStandardFrame(BackendKeyDataMessageType, messageBuffer)
}
