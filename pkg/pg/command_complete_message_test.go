package pg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Real command completion packet from pg captured with Wireshark
// Tag: SELECT 1
const GoldenCommandCompleteMesagePacket = "\x43\x00\x00\x00\x0d\x53\x45\x4c\x45\x43\x54\x20\x31\x00"

func TestParseCommandCompleteMessage(t *testing.T) {
	msg, err := ParseCommandCompleteMessage(StandardFrame(GoldenCommandCompleteMesagePacket))

	assert.NoError(t, err)
	assert.Equal(t, "SELECT 1", msg.Tag)

	// Test invalid type
	_, err = ParseCommandCompleteMessage(append(StandardFrame{'X'}, GoldenCommandCompleteMesagePacket[1:]...))

	assert.Equal(t, ErrMalformedMessage, err)
}

func TestCommandCompleteMessageFrame(t *testing.T) {
	msg := &CommandCompleteMessage{"SELECT 1"}
	frame := msg.Frame()

	assert.Equal(t, []byte(GoldenCommandCompleteMesagePacket), frame.Bytes())
}
