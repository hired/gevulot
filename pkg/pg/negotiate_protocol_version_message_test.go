package pg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Real negotiate protocol message packet from pg captured with Wireshark
//   Supported protocol version: 196608
//   Unrecognized options: _pq_.test
const GoldenNegotiateProtocolVersionMessagePacket = "\x76\x00\x00\x00\x16\x00\x03\x00\x00\x00\x00\x00\x01\x5f\x70\x71" +
	"\x5f\x2e\x74\x65\x73\x74\x00"

func TestParseNegotiateProtocolVersionMessage(t *testing.T) {
	{
		msg, err := ParseNegotiateProtocolVersionMessage(StandardFrame(GoldenNegotiateProtocolVersionMessagePacket))

		assert.NoError(t, err)
		assert.Equal(t, int32(196608), msg.SupportedProtocolVersion)
		assert.Equal(t, []string{"_pq_.test"}, msg.UnrecognizedOptions)
	}

	// Test invalid type
	{
		_, err := ParseNegotiateProtocolVersionMessage(append(StandardFrame{'X'}, GoldenNegotiateProtocolVersionMessagePacket[1:]...))

		assert.Equal(t, ErrMalformedMessage, err)
	}
}

func TestNegotiateProtocolVersionMessageFrame(t *testing.T) {
	msg := &NegotiateProtocolVersionMessage{
		SupportedProtocolVersion: 196608,
		UnrecognizedOptions:      []string{"_pq_.test"},
	}

	assert.Equal(t, []byte(GoldenNegotiateProtocolVersionMessagePacket), msg.Frame().Bytes())
}
