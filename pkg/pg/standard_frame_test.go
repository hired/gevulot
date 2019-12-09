package pg

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStandardFrame(t *testing.T) {
	frameData := []byte{
		'X',                    // message type
		0x00, 0x00, 0x00, 0x11, // frame size including this 4 bytes
		'm', 'e', 's', 's', 'a', 'g', 'e',
	}

	frame := StandardFrame(frameData)

	assert.Equal(t, byte('X'), frame.MessageType())
	assert.Equal(t, []byte("message"), frame.MessageBody())
	assert.Equal(t, frameData, frame.Bytes())
}

func TestNewStandardFrame(t *testing.T) {
	frame := NewStandardFrame('X', []byte("hello world"))

	assert.Equal(
		t,
		[]byte{'X', 0x00, 0x00, 0x00, 0x0f, 'h', 'e', 'l', 'l', 'o', ' ', 'w', 'o', 'r', 'l', 'd'},
		frame.Bytes(),
	)
}

func TestReadStandardFrame(t *testing.T) {
	testFrameData := []byte{
		'X',
		0x00, 0x00, 0x00, 0x09, // frame size including this 4 bytes
		'h', 'e', 'l', 'l', 'o',
	}

	r := bytes.NewBuffer(testFrameData)
	frame, err := ReadStandardFrame(r)

	assert.NoError(t, err)
	assert.Equal(t, testFrameData, frame.Bytes())
	assert.Equal(t, byte('X'), frame.MessageType())
	assert.Equal(t, []byte("hello"), frame.MessageBody())

	// Test IO error
	frame, err = ReadStandardFrame(r) // r is empty
	assert.Equal(t, io.EOF, err)
}
