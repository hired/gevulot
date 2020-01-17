package pg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Real password message packet from psql captured with Wireshark
// Password: password
const GoldenPasswordMessagePacket = "\x70\x00\x00\x00\x0d\x70\x61\x73\x73\x77\x6f\x72\x64\x00"

func TestParsePasswordMessage(t *testing.T) {
	{
		msg, err := ParsePasswordMessage(StandardFrame(GoldenPasswordMessagePacket))

		assert.NoError(t, err)
		assert.Equal(t, "password", msg.Password)
	}

	// Test invalid type
	{
		_, err := ParsePasswordMessage(append(StandardFrame{'X'}, GoldenPasswordMessagePacket[1:]...))

		assert.Equal(t, ErrMalformedMessage, err)
	}
}

func TestPasswordMessageFrame(t *testing.T) {
	msg := &PasswordMessage{"password"}
	assert.Equal(t, []byte(GoldenPasswordMessagePacket), msg.Frame().Bytes())
}
