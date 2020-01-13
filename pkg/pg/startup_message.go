package pg

// SSLRequestMagic is magic protocol version that client is using to request a SSL session.
const SSLRequestMagic = 80877103

// DefaultProtocolVersion is a default protocol version of PostgreSQL since version 7
const DefaultProtocolVersion = 196608

// StartupMessage represents the first message a client sends to the DB after establishing a connection.
// This message includes the names of the user and of the database the user wants to connect to;
// it also identifies the particular protocol version to be used as well as any additional run-time parameters.
type StartupMessage struct {
	ProtocolVersion int32
	Parameters      []StartupMessageParameter // NB: we want to preserve order so we can't use map here
}

// StartupMessageParameter represents run-time parameter in a StartupMessage.
type StartupMessageParameter struct {
	Name  string
	Value string
}

// Compile time check to make sure that StartupMessage implements the Message interface.
var _ Message = &StartupMessage{}

// ParseStartupMessage parses StartupMessage from a network frame.
func ParseStartupMessage(frame Frame) (*StartupMessage, error) {
	messageData := ReadBuffer(frame.MessageBody())

	// Start decoding the message
	protocolVersion, err := messageData.ReadInt32()

	if err != nil {
		return nil, err
	}

	var parameters []StartupMessageParameter

	// NB: startup message won't contain any parameters if this is SSL negotiation message
	if messageData.Len() > 0 {
		// Dictionary encoded as tuples of strings + ending \0
		for {
			key, err := messageData.ReadString()

			if err != nil {
				return nil, err
			}

			if key == "" {
				break
			}

			value, err := messageData.ReadString()

			if err != nil {
				return nil, err
			}

			parameters = append(parameters, StartupMessageParameter{key, value})
		}
	}

	message := &StartupMessage{
		ProtocolVersion: protocolVersion,
		Parameters:      parameters,
	}

	return message, nil
}

// Frame serializes the message into a network frame.
func (m *StartupMessage) Frame() Frame {
	var messageBuffer WriteBuffer

	messageBuffer.WriteInt32(m.ProtocolVersion)

	hasParams := false

	// Dictionary is encoded as string key-value pair + ending /0x00
	for _, param := range m.Parameters {
		messageBuffer.WriteString(param.Name)
		messageBuffer.WriteString(param.Value)

		hasParams = true
	}

	if hasParams {
		messageBuffer.WriteByte(0)
	}

	return NewStartupFrame(messageBuffer)
}

// GetParameter returns parameter value by its name. If parameter is not set it returns an empty string.
func (m *StartupMessage) GetParameter(name string) string {
	for _, param := range m.Parameters {
		if param.Name == name {
			return param.Value
		}
	}

	return ""
}
