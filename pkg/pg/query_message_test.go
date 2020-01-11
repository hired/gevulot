package pg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Real query message packet from psql captured with Wireshark
// SQL: SELECT count(*) FROM users;
const GoldenQueryMesage = "\x51\x00\x00\x00\x20\x53\x45\x4c\x45\x43\x54\x20\x63\x6f\x75\x6e" +
	"\x74\x28\x2a\x29\x20\x46\x52\x4f\x4d\x20\x75\x73\x65\x72\x73\x3b\x00"

func TestParseQueryMessage(t *testing.T) {
	msg, err := ParseQueryMessage(StandardFrame(GoldenQueryMesage))

	assert.NoError(t, err)
	assert.Equal(t, "SELECT count(*) FROM users;", msg.Query)

	// Test invalid type
	_, err = ParseQueryMessage(StandardFrame(append([]byte{'X'}, GoldenQueryMesage[1:]...)))

	assert.Equal(t, ErrMalformedMessage, err)
}

func TestQueryMessageFrame(t *testing.T) {
	msg := &QueryMessage{
		Query: "SELECT count(*) FROM users;",
	}

	frame := msg.Frame()

	assert.Equal(t, []byte(GoldenQueryMesage), frame.Bytes())
}
