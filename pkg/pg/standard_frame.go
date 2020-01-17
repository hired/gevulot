package pg

import (
	"bytes"
	"io"
)

// StandardFrame conveys arbitrary PostgreSQL message.
// The first byte of a frame identifies the message type, and the next four bytes specify the length
// of the the attached message (this length count includes itself, but not the message-type byte).
type StandardFrame []byte

// Compile time check to make sure that StandardFrame implements the Frame interface.
var _ Frame = StandardFrame{}

// MessageType returns type of the message.
func (f StandardFrame) MessageType() byte {
	return f[0]
}

// MessageBody returns bytes of the message.
func (f StandardFrame) MessageBody() []byte {
	return f[5:]
}

// Bytes returns this frame as a slice of bytes.
func (f StandardFrame) Bytes() []byte {
	return f
}

// NewStandardFrame initializes a new StandardFrame with the given type and message.
func NewStandardFrame(messageType byte, messageBody []byte) StandardFrame {
	var buffer WriteBuffer

	buffer.WriteByte(messageType)
	buffer.WriteInt32(int32(len(messageBody) + 4))
	buffer.WriteBytes(messageBody)

	return StandardFrame(buffer)
}

// ReadStandardFrame reads a standard frame from the given reader.
func ReadStandardFrame(r io.Reader) (StandardFrame, error) {
	buffer := &bytes.Buffer{}

	// Grow the buffer to reduce memory allocations
	buffer.Grow(512)

	// Read frame header
	_, err := io.CopyN(buffer, r, 5) // 1 byte of message type + 4 bytes of message length

	if err != nil {
		return nil, err
	}

	// Decode frame header
	frameHeader := ReadBuffer(buffer.Bytes())

	// Skip the first byte
	_, err = frameHeader.ReadByte()

	if err != nil {
		return nil, err
	}

	// Get the message length
	frameLength, err := frameHeader.ReadInt32()

	if err != nil {
		return nil, err
	}

	// Read the rest of the frame
	_, err = io.CopyN(buffer, r, int64(frameLength-4))

	if err != nil {
		return nil, err
	}

	// Cast all read bytes to StandardFrame
	return StandardFrame(buffer.Bytes()), nil
}
