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

// Marshall serializes the message to send it over the network.
func (m *StartupMessage) Marshall() ([]byte, error) {
	var paramsBuffer WriteBuffer

	// Dictionary is encoded as string key-value pair + ending /0x00
	for _, param := range m.Parameters {
		paramsBuffer.WriteString(param.Name)
		paramsBuffer.WriteString(param.Value)
	}

	if paramsBuffer.Len() > 0 {
		paramsBuffer.WriteByte(0)
	}

	// Total message length including length itself
	messageLength := paramsBuffer.Len() + 8 // 8 = length (int32) + protocol version (int32)
	messageBuffer := make(WriteBuffer, 0, messageLength)

	// NB: StartupMessage does not have a leading type byte!
	messageBuffer.WriteInt32(int32(messageLength))
	messageBuffer.WriteInt32(m.ProtocolVersion)
	messageBuffer.WriteBytes(paramsBuffer)

	return messageBuffer, nil
}
