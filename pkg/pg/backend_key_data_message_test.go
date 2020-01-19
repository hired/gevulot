package pg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Real key data message packet from pg captured with Wireshark
//   PID: 28822
//   Key: -636998870
const GoldenBakendKeyDataMesagePacket = "\x4b\x00\x00\x00\x0c\x00\x00\x70\x96\xda\x08\x2b\x2a"

func TestParseBackendKeyDataMessage(t *testing.T) {
	{
		msg, err := ParseBackendKeyDataMessage(StandardFrame(GoldenBakendKeyDataMesagePacket))

		assert.NoError(t, err)
		assert.Equal(t, int32(28822), msg.ProcessID)
		assert.Equal(t, int32(-636998870), msg.Key)
	}

	// Test invalid type
	{
		_, err := ParseQueryMessage(append(StandardFrame{'X'}, GoldenBakendKeyDataMesagePacket[1:]...))
		assert.Equal(t, ErrMalformedMessage, err)
	}
}

func TestBackendKeyDataMessageFrame(t *testing.T) {
	msg := &BackendKeyDataMessage{ProcessID: 28822, Key: -636998870}
	assert.Equal(t, []byte(GoldenBakendKeyDataMesagePacket), msg.Frame().Bytes())
}
