package pg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Real terminate message packet from psql captured with Wireshark
const GoldenTerminateMesagePacket = "\x58\x00\x00\x00\x04"

func TestParseTerminateMessage(t *testing.T) {
	{
		_, err := ParseTerminateMessage(StandardFrame(GoldenTerminateMesagePacket))
		assert.NoError(t, err)
	}

	// Test invalid type
	{
		_, err := ParseTerminateMessage(append(StandardFrame{'!'}, GoldenTerminateMesagePacket[1:]...))
		assert.Equal(t, ErrMalformedMessage, err)
	}
}

func TestTerminateMessageFrame(t *testing.T) {
	msg := &TerminateMessage{}
	frame := msg.Frame()

	assert.Equal(t, []byte(GoldenTerminateMesagePacket), frame.Bytes())
}
