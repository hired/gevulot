package pg

import (
	"testing"

	"github.com/lib/pq/oid"
	"github.com/stretchr/testify/assert"
)

const GoldenRowDescriptionMessagePacket = "\x54\x00\x00\x00\x4a\x00\x03\x69\x64\x00\x00\x0c\x06\x01\x00\x11" +
	"\x00\x00\x00\x17\x00\x04\xff\xff\xff\xff\x00\x00\x6e\x61\x6d\x65" +
	"\x00\x00\x0c\x06\x01\x00\x10\x00\x00\x04\x13\xff\xff\x00\x00\x01" +
	"\x03\x00\x00\x65\x6d\x61\x69\x6c\x00\x00\x0c\x06\x01\x00\x02\x00" +
	"\x00\x04\x13\xff\xff\x00\x00\x01\x03\x00\x00"

func TestParseRowDescriptionMessage(t *testing.T) {
	msg, err := ParseRowDescriptionMessage(StandardFrame(GoldenRowDescriptionMessagePacket))

	assert.NoError(t, err)
	assert.Equal(t, 3, len(msg.Fields))

	assert.Equal(t, "id", msg.Fields[0].Name)
	assert.Equal(t, oid.T_int4, msg.Fields[0].DataTypeOID)

	assert.Equal(t, "name", msg.Fields[1].Name)
	assert.Equal(t, oid.T_varchar, msg.Fields[1].DataTypeOID)

	assert.Equal(t, "email", msg.Fields[2].Name)
	assert.Equal(t, oid.T_varchar, msg.Fields[2].DataTypeOID)

	// Test invalid message type
	_, err = ParseRowDescriptionMessage(append(StandardFrame{'X'}, GoldenRowDescriptionMessagePacket[1:]...))
	assert.Equal(t, ErrMalformedMessage, err)
}

func TestRowDescriptionMessageFrame(t *testing.T) {
	fields := []*FieldDescriptor{
		{
			Name:             "id",
			TableOID:         787969,
			ColumnIndex:      17,
			DataTypeOID:      oid.T_int4,
			DataTypeSize:     4,
			DataTypeModifier: -1,
			Format:           DataFormatText,
		},
		{
			Name:             "name",
			TableOID:         787969,
			ColumnIndex:      16,
			DataTypeOID:      oid.T_varchar,
			DataTypeSize:     -1,
			DataTypeModifier: 259,
			Format:           DataFormatText,
		},
		{
			Name:             "email",
			TableOID:         787969,
			ColumnIndex:      2,
			DataTypeOID:      oid.T_varchar,
			DataTypeSize:     -1,
			DataTypeModifier: 259,
			Format:           DataFormatText,
		},
	}

	msg := &RowDescriptionMessage{
		Fields: fields,
	}

	frame := msg.Frame()

	assert.Equal(t, []byte(GoldenRowDescriptionMessagePacket), frame.Bytes())
}
