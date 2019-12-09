package pg

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
)

// ReadBuffer provides bufio.Scanner-like API to decode PG data types described here:
// https://www.postgresql.org/docs/9.6/protocol-message-types.html
//
// NB: buffer is read-only! You cannot reset the state!
type ReadBuffer []byte

//
// Main use case:
//
//   buf := make(PgReadBuffer, 31337) // Allocate buffer
//   io.ReadFull(r, buf)              // Write data to the buffer
//   value, err := buf.ReadInt32()    // Start decoding
//   value, err := buf.ReadString()
//   ...
//

// Len returns the number of bytes left in the buffer.
func (b *ReadBuffer) Len() int {
	return len(*b)
}

// ReadBytes reads and returns exactly n bytes from the buffer. The bytes are removed from the buffer afterwards.
// If there are not enough bytes in the buffer available, it returns error io.EOF and does not advance.
func (b *ReadBuffer) ReadBytes(n int) ([]byte, error) {
	if b.Len() < n {
		return nil, io.EOF
	}

	bytesRead := (*b)[:n]

	// Advance the slice: n = 3, b = { 1, 2, 3, 4, 5 } -> b = { 4, 5 }
	*b = (*b)[n:]

	return bytesRead, nil
}

// ReadByte reads and returns the next byte from the buffer. The byte is removed from the buffer afterwards.
// If no byte is available, it returns error io.EOF.
func (b *ReadBuffer) ReadByte() (byte, error) {
	value, err := b.ReadBytes(1)

	if err != nil {
		return 0, err
	}

	return value[0], nil
}

// ReadInt32 decodes a 32-bit signed integer and advances the buffer over it.
// If there are not enough bytes in the buffer available, it returns error io.EOF and does not advance.
func (b *ReadBuffer) ReadInt32() (int32, error) {
	intBytes, err := b.ReadBytes(4)

	if err != nil {
		return 0, err
	}

	return int32(binary.BigEndian.Uint32(intBytes)), nil
}

// ReadInt16 decodes a 16-bit signed integer and advances the buffer over it.
// If there are not enough bytes in the buffer available, it returns error io.EOF and does not advance.
func (b *ReadBuffer) ReadInt16() (int16, error) {
	intBytes, err := b.ReadBytes(2)

	if err != nil {
		return 0, err
	}

	return int16(binary.BigEndian.Uint16(intBytes)), nil
}

// ReadString reads a null-terminated string and advances the buffer over it.
// If there is no terminator present in the buffer, it returns error and does not advance.
func (b *ReadBuffer) ReadString() (string, error) {
	// Find the first occurrence of \0x00
	nullIndex := bytes.IndexByte(*b, 0)

	if nullIndex < 0 {
		return "", fmt.Errorf("invalid message format: expected string terminator")
	}

	// Read string bytes including \0x00
	stringBytes, err := b.ReadBytes(nullIndex + 1)

	if err != nil {
		return "", err
	}

	// Convert to string removing \0x00
	return string(stringBytes[:nullIndex]), nil
}

// ReadInt32Array decodes an array of k 32-bit signed integers and advances the buffer over it.
// If there are not enough bytes in the buffer available, it returns error io.EOF and does not advance.
func (b *ReadBuffer) ReadInt32Array(k int) ([]int32, error) {
	arrBytes, err := b.ReadBytes(k * 4)

	if err != nil {
		return nil, err
	}

	arr := make([]int32, k)

	for i := range arr {
		arr[i] = int32(binary.BigEndian.Uint32(arrBytes[4*i:]))
	}

	return arr, nil
}

// ReadInt16Array decodes an array of k 16-bit signed integers and advances the buffer over it.
// If there are not enough bytes in the buffer available, it returns error io.EOF and does not advance.
func (b *ReadBuffer) ReadInt16Array(k int) ([]int16, error) {
	arrBytes, err := b.ReadBytes(k * 2)

	if err != nil {
		return nil, err
	}

	arr := make([]int16, k)

	for i := range arr {
		arr[i] = int16(binary.BigEndian.Uint16(arrBytes[2*i:]))
	}

	return arr, nil
}

// Inspect returns content of the buffer as a hex dump for debug purposes.
func (b *ReadBuffer) Inspect() string {
	return fmt.Sprintf("buffer length: %d\n%s", b.Len(), hex.Dump(*b))
}
