package pg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Real empty query response message packet from pg captured with Wireshark
const GoldenEmptyQueryResponseMesagePacket = "\x49\x00\x00\x00\x04"

func TestParseEmptyQueryResponseMessage(t *testing.T) {
	{
		_, err := ParseEmptyQueryResponseMessage(StandardFrame(GoldenEmptyQueryResponseMesagePacket))
		assert.NoError(t, err)
	}

	// Test invalid type
	{
		_, err := ParseEmptyQueryResponseMessage(append(StandardFrame{'X'}, GoldenEmptyQueryResponseMesagePacket[1:]...))
		assert.Equal(t, ErrMalformedMessage, err)
	}
}

func TestEmptyQueryResponseMessageFrame(t *testing.T) {
	msg := &EmptyQueryResponseMessage{}
	assert.Equal(t, []byte(GoldenEmptyQueryResponseMesagePacket), msg.Frame().Bytes())
}
