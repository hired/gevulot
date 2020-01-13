package pg

// DataRowMessageType identifies DataRowMessage message.
const DataRowMessageType = 'D'

// DataRowMessage represents a single data row sent by a backend as a query result.
// RowDescriptionMessage always precedes the DataRowMessage.
type DataRowMessage struct {
	// Row raw values. To decode them you need RowDescriptionMessage.
	Values [][]byte
}

// Compile time check to make sure that DataRowMessage implements the Message interface.
var _ Message = &DataRowMessage{}

// ParseRowDescriptionMessage parses RowDescriptionMessage from a network frame.
func ParseDataRowMessage(frame Frame) (*DataRowMessage, error) {
	// Assert the message type
	if frame.MessageType() != DataRowMessageType {
		return nil, ErrMalformedMessage
	}

	messageData := ReadBuffer(frame.MessageBody())

	// Number of columns in the row (could be 0)
	valuesCount, err := messageData.ReadInt16()

	if err != nil {
		return nil, err
	}

	values := make([][]byte, valuesCount)

	for i := 0; i < int(valuesCount); i++ {
		size, err := messageData.ReadInt32()

		if err != nil {
			return nil, err
		}

		if size == -1 {
			continue
		}

		values[i], err = messageData.ReadBytes(int(size))

		if err != nil {
			return nil, err
		}
	}

	return &DataRowMessage{Values: values}, nil
}

// Frame serializes the message into a network frame.
func (m *DataRowMessage) Frame() Frame {
	var messageBuffer WriteBuffer

	messageBuffer.WriteInt16(int16(len(m.Values)))

	for i := 0; i < len(m.Values); i++ {
		value := m.Values[i]

		if value == nil {
			// -1 represents NULL
			messageBuffer.WriteInt32(-1)
			continue
		}

		messageBuffer.WriteInt32(int32(len(value)))
		messageBuffer.WriteBytes(value)
	}

	return NewStandardFrame(DataRowMessageType, messageBuffer)
}
