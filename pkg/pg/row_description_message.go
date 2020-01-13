package pg

import (
	"github.com/lib/pq/oid"
)

// RowDescriptionMessageType identifies RowDescriptionMessage message.
const RowDescriptionMessageType = 'T'

// RowDescriptionMessage represent a message sent by a backend to describe query result fields.
type RowDescriptionMessage struct {
	// List of fields to be returned
	Fields []*FieldDescriptor
}

// FieldDescriptor describes a field of a DataRow.
type FieldDescriptor struct {
	// The field name.
	Name string

	// If the field can be identified as a column of a specific table, the object ID of the table; otherwise zero.
	TableOID oid.Oid

	// If the field can be identified as a column of a specific table, the attribute number of the column; otherwise zero.
	ColumnIndex int16

	// The object ID of the field's data type.
	DataTypeOID oid.Oid

	// The data type size (see pg_type.typlen). Note that negative values denote variable-width types.
	DataTypeSize int16

	// The type modifier (see pg_attribute.atttypmod). The meaning of the modifier is type-specific.
	DataTypeModifier int32

	// The format code being used for the field. Currently will be zero (text) or one (binary). In a RowDescription
	// returned from the statement variant of Describe, the format code is not yet known and will always be zero.
	Format DataFormat
}

// DataFormat is a code for PostgreSQL data format.
type DataFormat int16

const (
	DataFormatText   DataFormat = 0 // Plain text
	DataFormatBinary DataFormat = 1 // Binary representation
)

// Compile time check to make sure that RowDescriptionMessage implements the Message interface.
var _ Message = &RowDescriptionMessage{}

// ParseRowDescriptionMessage parses RowDescriptionMessage from the network frame.
func ParseRowDescriptionMessage(frame Frame) (*RowDescriptionMessage, error) {
	// Assert the message type
	if frame.MessageType() != RowDescriptionMessageType {
		return nil, ErrMalformedMessage
	}

	messageData := ReadBuffer(frame.MessageBody())

	// Number of fields (could be 0)
	fieldsCount, err := messageData.ReadInt16()

	if err != nil {
		return nil, err
	}

	fields := make([]*FieldDescriptor, fieldsCount)

	for i := 0; i < int(fieldsCount); i++ {
		name, err := messageData.ReadString()

		if err != nil {
			return nil, err
		}

		tableOID, err := messageData.ReadInt32()

		if err != nil {
			return nil, err
		}

		columnNum, err := messageData.ReadInt16()

		if err != nil {
			return nil, err
		}

		dataTypeOID, err := messageData.ReadInt32()

		if err != nil {
			return nil, err
		}

		dataTypeSize, err := messageData.ReadInt16()

		if err != nil {
			return nil, err
		}

		dataTypeModifier, err := messageData.ReadInt32()

		if err != nil {
			return nil, err
		}

		format, err := messageData.ReadInt16()

		if err != nil {
			return nil, err
		}

		fields[i] = &FieldDescriptor{
			Name:             name,
			TableOID:         oid.Oid(tableOID),
			ColumnIndex:      columnNum,
			DataTypeOID:      oid.Oid(dataTypeOID),
			DataTypeSize:     dataTypeSize,
			DataTypeModifier: dataTypeModifier,
			Format:           DataFormat(format),
		}
	}

	return &RowDescriptionMessage{Fields: fields}, nil
}

// Frame serializes the message into a network frame.
func (m *RowDescriptionMessage) Frame() Frame {
	var messageBuffer WriteBuffer

	messageBuffer.WriteInt16(int16(len(m.Fields)))

	for i := 0; i < len(m.Fields); i++ {
		f := m.Fields[i]

		messageBuffer.WriteString(f.Name)
		messageBuffer.WriteInt32(int32(f.TableOID))
		messageBuffer.WriteInt16(f.ColumnIndex)
		messageBuffer.WriteInt32(int32(f.DataTypeOID))
		messageBuffer.WriteInt16(f.DataTypeSize)
		messageBuffer.WriteInt32(f.DataTypeModifier)
		messageBuffer.WriteInt16(int16(f.Format))
	}

	return NewStandardFrame(RowDescriptionMessageType, messageBuffer)
}
