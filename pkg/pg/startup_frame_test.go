package pg

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStartupFrame(t *testing.T) {
	frameData := []byte{
		0x00, 0x00, 0x00, 0x11, // frame size including this 4 bytes
		'm', 'e', 's', 's', 'a', 'g', 'e',
	}

	frame := StartupFrame(frameData)

	assert.Equal(t, byte(0x00), frame.MessageType(), "StartupMessage doesn't have a message type")
	assert.Equal(t, []byte("message"), frame.MessageBody())
	assert.Equal(t, frameData, frame.Bytes())
}

func TestNewStartupFrame(t *testing.T) {
	frame := NewStartupFrame([]byte("hello world"))

	assert.Equal(
		t,
		[]byte{0x00, 0x00, 0x00, 0x0f, 'h', 'e', 'l', 'l', 'o', ' ', 'w', 'o', 'r', 'l', 'd'},
		frame.Bytes(),
	)
}

func TestReadStartupFrame(t *testing.T) {
	testFrameData := []byte{
		0x00, 0x00, 0x00, 0x09, // frame size including this 4 bytes
		'h', 'e', 'l', 'l', 'o',
	}

	r := bytes.NewBuffer(testFrameData)
	frame, err := ReadStartupFrame(r)

	assert.NoError(t, err)
	assert.Equal(t, testFrameData, frame.Bytes())
	assert.Equal(t, []byte("hello"), frame.MessageBody())

	// Test IO error
	frame, err = ReadStartupFrame(r) // r is empty
	assert.Equal(t, io.EOF, err)
}
