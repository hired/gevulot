package pg

// ReadyForQueryMessageType identifies ReadyForQueryMessage message.
const ReadyForQueryMessageType = 'Z'

// ReadyForQueryMessage is sent whenever the backend is ready for a new query cycle.
type ReadyForQueryMessage struct {
	// Current backend transaction status
	TxStatus TxStatus
}

// TxStatus represents transaction status code.
type TxStatus byte

const (
	TxStatusIdle   TxStatus = 'I' // Not in a transaction block
	TxStatusActive TxStatus = 'T' // In a transaction block
	TxStatusFailed TxStatus = 'E' // Failed transaction block
)

// Compile time check to make sure that QueryMessage implements the Message interface.
var _ Message = &ReadyForQueryMessage{}

// ParseReadyForQueryMessage parses ReadyForQueryMessage from a network frame.
func ParseReadyForQueryMessage(frame Frame) (*ReadyForQueryMessage, error) {
	// Assert the message type
	if frame.MessageType() != ReadyForQueryMessageType {
		return nil, ErrMalformedMessage
	}

	messageData := ReadBuffer(frame.MessageBody())

	status, err := messageData.ReadByte()

	if err != nil {
		return nil, err
	}

	return &ReadyForQueryMessage{TxStatus: TxStatus(status)}, nil
}

// Frame serializes the message into a network frame.
func (m *ReadyForQueryMessage) Frame() Frame {
	return NewStandardFrame(ReadyForQueryMessageType, []byte{byte(m.TxStatus)})
}
