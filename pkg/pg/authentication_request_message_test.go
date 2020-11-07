package pg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Real clear-text password request message packet from pg captured with Wireshark
const GoldenAuthenticationCleartextPasswordMessagePacket = "\x52\x00\x00\x00\x08\x00\x00\x00\x03"

// Real MD5 password request message packet from pg captured with Wireshark
// Salt: 8f012d4c
//nolint:gosec
const GoldenAuthenticationMD5PasswordMessagePacket = "\x52\x00\x00\x00\x0c\x00\x00\x00\x05\x8f\x01\x2d\x4c"

// Real auth ok message packet from pg captured with Wireshark
const GoldenAuthenticationOkMessagePacket = "\x52\x00\x00\x00\x08\x00\x00\x00\x00"

func TestParseAuthenticationRequestMessage(t *testing.T) {
	{
		msg, err := ParseAuthenticationRequestMessage(StandardFrame(GoldenAuthenticationCleartextPasswordMessagePacket))

		assert.NoError(t, err)
		assert.IsType(t, &AuthenticationCleartextPasswordMessage{}, msg)
	}

	{
		msg, err := ParseAuthenticationRequestMessage(StandardFrame(GoldenAuthenticationMD5PasswordMessagePacket))

		assert.NoError(t, err)

		if assert.IsType(t, &AuthenticationMD5PasswordMessage{}, msg) {
			assert.Equal(t, [4]byte{0x8f, 0x01, 0x2d, 0x4c}, msg.(*AuthenticationMD5PasswordMessage).Salt)
		}
	}

	{
		msg, err := ParseAuthenticationRequestMessage(StandardFrame(GoldenAuthenticationOkMessagePacket))

		assert.NoError(t, err)
		assert.IsType(t, &AuthenticationOkMessage{}, msg)
	}

	// Test invalid type
	{
		_, err := ParseAuthenticationRequestMessage(append(StandardFrame{'X'}, GoldenAuthenticationOkMessagePacket[1:]...))

		assert.Equal(t, ErrMalformedMessage, err)
	}
}

func TestAuthenticationRequestMessageFrame(t *testing.T) {
	{
		msg := &AuthenticationCleartextPasswordMessage{}
		assert.Equal(t, []byte(GoldenAuthenticationCleartextPasswordMessagePacket), msg.Frame().Bytes())
	}

	{
		msg := &AuthenticationMD5PasswordMessage{[4]byte{0x8f, 0x01, 0x2d, 0x4c}}
		assert.Equal(t, []byte(GoldenAuthenticationMD5PasswordMessagePacket), msg.Frame().Bytes())
	}

	{
		msg := &AuthenticationOkMessage{}
		assert.Equal(t, []byte(GoldenAuthenticationOkMessagePacket), msg.Frame().Bytes())
	}
}
