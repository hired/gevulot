package pg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Real parameter status message packet from pg captured with Wireshark
//   Name: server_encoding
//   Value: UTF8
const GoldenParameterStatusMessagePacket = "\x53\x00\x00\x00\x19\x73\x65\x72\x76\x65\x72\x5f\x65\x6e\x63\x6f" +
	"\x64\x69\x6e\x67\x00\x55\x54\x46\x38\x00"

func TestParseParameterStatusMessage(t *testing.T) {
	{
		msg, err := ParseParameterStatusMessage(StandardFrame(GoldenParameterStatusMessagePacket))

		assert.NoError(t, err)
		assert.Equal(t, "server_encoding", msg.Name)
		assert.Equal(t, "UTF8", msg.Value)
	}

	// Test invalid type
	{
		_, err := ParseParameterStatusMessage(append(StandardFrame{'X'}, GoldenParameterStatusMessagePacket[1:]...))

		assert.Equal(t, ErrMalformedMessage, err)
	}
}

func TestParameterStatusMessageFrame(t *testing.T) {
	msg := &ParameterStatusMessage{"server_encoding", "UTF8"}
	assert.Equal(t, []byte(GoldenParameterStatusMessagePacket), msg.Frame().Bytes())
}
