package pg

// NegotiateProtocolVersionMessageType identifies NegotiateProtocolVersionMessage message.
const NegotiateProtocolVersionMessageType = 'v'

// NegotiateProtocolVersionMessage is sent by a backend when it does not support the minor protocol version requested by the frontend.
type NegotiateProtocolVersionMessage struct {
	// Newest minor protocol version supported by the backend for the major protocol version requested by the frontend.
	SupportedProtocolVersion int32

	// List of protocol options not recognized by the backend.
	UnrecognizedOptions []string
}

// Compile time check to make sure that NegotiateProtocolVersionMessage implements the Message interface.
var _ Message = &NegotiateProtocolVersionMessage{}

// ParseNegotiateProtocolVersionMessage parses NegotiateProtocolVersionMessage from a network frame.
func ParseNegotiateProtocolVersionMessage(frame Frame) (*NegotiateProtocolVersionMessage, error) {
	// Assert the message type
	if frame.MessageType() != NegotiateProtocolVersionMessageType {
		return nil, ErrMalformedMessage
	}

	messageData := ReadBuffer(frame.MessageBody())

	// Read the proto version
	supportedProtocol, err := messageData.ReadInt32()

	if err != nil {
		return nil, err
	}

	// Read unrecognized startup params (could be 0)
	unrecognizedOptionsCount, err := messageData.ReadInt32()

	if err != nil {
		return nil, err
	}

	unrecognizedOptions := make([]string, unrecognizedOptionsCount)

	for i := 0; i < int(unrecognizedOptionsCount); i++ {
		unrecognizedOptions[i], err = messageData.ReadString()

		if err != nil {
			return nil, err
		}
	}

	return &NegotiateProtocolVersionMessage{
		SupportedProtocolVersion: supportedProtocol,
		UnrecognizedOptions:      unrecognizedOptions,
	}, nil
}

// Frame serializes the message into a network frame.
func (m *NegotiateProtocolVersionMessage) Frame() Frame {
	var messageBuffer WriteBuffer

	messageBuffer.WriteInt32(m.SupportedProtocolVersion)
	messageBuffer.WriteInt32(int32(len(m.UnrecognizedOptions)))

	for _, option := range m.UnrecognizedOptions {
		messageBuffer.WriteString(option)
	}

	return NewStandardFrame(NegotiateProtocolVersionMessageType, messageBuffer)
}
