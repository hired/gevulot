package pg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Real data row message packet from pg captured with Wireshark
// SQL: SELECT id, old_id, name, email FROM users WHERE email = 'denis.diachkov@hired.com';
const GoldenDataRowMesagePacket = "\x44\x00\x00\x00\x43\x00\x04\x00\x00\x00\x07\x31\x34\x35\x32\x32" +
	"\x36\x38\xff\xff\xff\xff\x00\x00\x00\x0e\x44\x65\x6e\x69\x73\x20" +
	"\x44\x69\x61\x63\x68\x6b\x6f\x76\x00\x00\x00\x18\x64\x65\x6e\x69" +
	"\x73\x2e\x64\x69\x61\x63\x68\x6b\x6f\x76\x40\x68\x69\x72\x65\x64" +
	"\x2e\x63\x6f\x6d"

func TestParseDataRowMessage(t *testing.T) {
	msg, err := ParseDataRowMessage(StandardFrame(GoldenDataRowMesagePacket))

	assert.NoError(t, err)
	assert.Len(t, msg.Values, 4)
	assert.Nil(t, msg.Values[1])
	assert.Equal(t, []byte("denis.diachkov@hired.com"), msg.Values[3])

	// Test invalid type
	_, err = ParseDataRowMessage(append(StandardFrame{'X'}, GoldenDataRowMesagePacket[1:]...))

	assert.Equal(t, ErrMalformedMessage, err)
}

func TestDataRowMessageFrame(t *testing.T) {
	msg := &DataRowMessage{
		Values: [][]byte{
			{ 0x31, 0x34, 0x35, 0x32, 0x32, 0x36, 0x38 }, // id = 1452268
			nil, // old_id = null
			[]byte("Denis Diachkov"),
			[]byte("denis.diachkov@hired.com"),
		},
	}

	frame := msg.Frame()

	assert.Equal(t, []byte(GoldenDataRowMesagePacket), frame.Bytes())
}