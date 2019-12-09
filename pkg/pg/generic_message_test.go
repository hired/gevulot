package pg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseGenericMessage(t *testing.T) {
	frame := NewStandardFrame('X', []byte("test test"))
	msg, err := ParseGenericMessage(frame)

	assert.NoError(t, err)
	assert.Equal(t, frame.MessageType(), msg.Type)
	assert.Equal(t, frame.MessageBody(), msg.Body)
}

func TestGenericMessageFrame(t *testing.T) {
	msg := &GenericMessage{
		Type: 'X',
		Body: []byte("test"),
	}

	frame := msg.Frame()

	assert.Equal(t, []byte{'X', 0x00, 0x00, 0x00, 0x08, 't', 'e', 's', 't'}, frame.Bytes())
}
