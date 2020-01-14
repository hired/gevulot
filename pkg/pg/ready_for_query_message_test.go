package pg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Real ready for query message packet from pg captured with Wireshark
// Tx status = idle
const GoldenReadyForQueryMesagePacket = "\x5a\x00\x00\x00\x05\x49"

func TestParseReadyForQueryMessage(t *testing.T) {
	msg, err := ParseReadyForQueryMessage(StandardFrame(GoldenReadyForQueryMesagePacket))

	assert.NoError(t, err)
	assert.Equal(t, TxStatusIdle, msg.TxStatus)

	// Test invalid type
	_, err = ParseQueryMessage(append(StandardFrame{'X'}, GoldenReadyForQueryMesagePacket[1:]...))

	assert.Equal(t, ErrMalformedMessage, err)
}

func TestReadyForQueryMessageFrame(t *testing.T) {
	msg := &ReadyForQueryMessage{TxStatusIdle}
	frame := msg.Frame()

	assert.Equal(t, []byte(GoldenReadyForQueryMesagePacket), frame.Bytes())
}
