package pg

import (
	"bytes"
	"io"
)

// StartupFrame conveys StartupMessage — the first message a client sends to begin a session with a DB.
// The only difference between a StandardFrame and a StartupFrame is that the latter doesn't have a message type
// as the first byte of the frame header. It is exist in PG protocol for historical reason.
type StartupFrame []byte

// Compile time check to make sure that StartupFrame implements the Frame interface.
var _ Frame = StartupFrame{}

// MessageType returns 0x00 since StartupFrame doesn't have message type.
func (f StartupFrame) MessageType() byte {
	return 0x00
}

// MessageBody returns bytes of a message containing within this frame.
func (f StartupFrame) MessageBody() []byte {
	return f[4:]
}

// Bytes returns this frame as a slice of bytes.
func (f StartupFrame) Bytes() []byte {
	return f
}

// NewStartupFrame initializes a new StartupFrame with the given message.
func NewStartupFrame(messageBody []byte) StartupFrame {
	var buffer WriteBuffer

	// NB: no message type — just the frame length and the message
	buffer.WriteInt32(int32(len(messageBody) + 4))
	buffer.WriteBytes(messageBody)

	return StartupFrame(buffer)
}

// ReadStartupFrame reads start-up frame from the given reader.
func ReadStartupFrame(r io.Reader) (StartupFrame, error) {
	buffer := &bytes.Buffer{}

	// Grow the buffer to reduce memory allocations
	buffer.Grow(256)

	// Read frame header
	_, err := io.CopyN(buffer, r, 4) // 4 bytes of message length

	if err != nil {
		return nil, err
	}

	// Decode frame header
	frameHeader := ReadBuffer(buffer.Bytes())

	// First 4 bytes is the length of the frame.
	// Note that StartupFrame doesn't have message type!
	frameLength, err := frameHeader.ReadInt32()

	if err != nil {
		return nil, err
	}

	// Read the rest of the frame
	_, err = io.CopyN(buffer, r, int64(frameLength-4))

	if err != nil {
		return nil, err
	}

	// Cast all read bytes to StartupFrame
	return StartupFrame(buffer.Bytes()), nil
}
