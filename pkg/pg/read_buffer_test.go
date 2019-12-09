package pg

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadBufferLen(t *testing.T) {
	buf := make(ReadBuffer, 42)
	assert.Equal(t, buf.Len(), 42)
}

func TestReadBufferReadBytes(t *testing.T) {
	buf := ReadBuffer{0xDE, 0xAD, 0xBE, 0xEF, 0xFE, 0xE1, 0xDE, 0xAD}

	// Regular read
	bytesRead, err := buf.ReadBytes(4)

	assert.NoError(t, err)
	assert.Equal(t, bytesRead, []byte{0xDE, 0xAD, 0xBE, 0xEF})
	assert.Equal(t, buf.Len(), 4) // buffer advances

	// Not enough bytes
	_, err = buf.ReadBytes(100)

	assert.Equal(t, err, io.EOF)
	assert.Equal(t, buf.Len(), 4) // buffer doesn't advance
}

func TestReadBufferReadByte(t *testing.T) {
	buf := ReadBuffer{0xFF}

	// Regular read
	byteRead, err := buf.ReadByte()

	assert.NoError(t, err)
	assert.Equal(t, byteRead, byte(0xFF))
	assert.Equal(t, buf.Len(), 0) // buffer advances

	// Buffer is empty
	_, err = buf.ReadByte()

	assert.Equal(t, err, io.EOF)
}

func TestReadBufferReadInt32(t *testing.T) {
	buf := ReadBuffer{
		// 31337 in Big-endian
		0x0, 0x0, 0x7A, 0x69,

		// 2 bytes — not enough for 32 bit integer
		0xCA, 0xFE,
	}

	// Regular read
	num, err := buf.ReadInt32()

	assert.NoError(t, err)
	assert.Equal(t, num, int32(31337))
	assert.Equal(t, buf.Len(), 2) // buffer advances

	// Not enough bytes
	_, err = buf.ReadInt32()

	assert.Equal(t, err, io.EOF)
	assert.Equal(t, buf.Len(), 2) // buffer doesn't advance
}

func TestReadBufferReadInt16(t *testing.T) {
	buf := ReadBuffer{
		// 42 in Big-endian
		0x0, 0x2A,

		// 1 byte — not enough for 16 bit integer
		0xFF,
	}

	// Regular read
	num, err := buf.ReadInt16()

	assert.NoError(t, err)
	assert.Equal(t, num, int16(42))
	assert.Equal(t, buf.Len(), 1) // buffer advances

	// Not enough bytes
	_, err = buf.ReadInt16()

	assert.Equal(t, err, io.EOF)
	assert.Equal(t, buf.Len(), 1) // buffer doesn't advance
}

func TestReadBufferReadString(t *testing.T) {
	buf := ReadBuffer{
		'h', 'e', 'l', 'l', 'o', ' ', 'w', 'o', 'r', 'l', 'd', 0x00, // 12 bytes
		'n', 'o', ' ', 't', 'e', 'r', 'm', 'i', 'n', 'a', 't', 'o', 'r', // 13 bytes
	}

	// Regular read
	str, err := buf.ReadString()

	assert.NoError(t, err)
	assert.Equal(t, str, "hello world")
	assert.Equal(t, buf.Len(), 13) // buffer advances

	// Cannot find null byte
	_, err = buf.ReadString()

	assert.EqualError(t, err, "invalid message format: expected string terminator")
	assert.Equal(t, buf.Len(), 13) // buffer doesn't advance
}

func TestReadBufferReadInt32Array(t *testing.T) {
	buf := ReadBuffer{
		// 1, 2, 3, 4, 5, 6 in Big-endian
		0x0, 0x0, 0x0, 0x1,
		0x0, 0x0, 0x0, 0x2,
		0x0, 0x0, 0x0, 0x3,
		0x0, 0x0, 0x0, 0x4,
		0x0, 0x0, 0x0, 0x5,
		0x0, 0x0, 0x0, 0x6,
	}

	// Regular read
	arr, err := buf.ReadInt32Array(5)

	assert.NoError(t, err)
	assert.Equal(t, arr, []int32{1, 2, 3, 4, 5})
	assert.Equal(t, buf.Len(), 4) // buffer advances

	// Not enough bytes
	_, err = buf.ReadInt32Array(100)

	assert.Equal(t, err, io.EOF)
	assert.Equal(t, buf.Len(), 4) // buffer doesn't advance
}

func TestReadBufferReadInt16Array(t *testing.T) {
	buf := ReadBuffer{
		// 6, 5, 4, 3, 2, 1 in Big-endian
		0x0, 0x6,
		0x0, 0x5,
		0x0, 0x4,
		0x0, 0x3,
		0x0, 0x2,
		0x0, 0x1,
	}

	// Regular read
	arr, err := buf.ReadInt16Array(5)

	assert.NoError(t, err)
	assert.Equal(t, arr, []int16{6, 5, 4, 3, 2})
	assert.Equal(t, buf.Len(), 2) // buffer advances

	// Not enough bytes
	_, err = buf.ReadInt16Array(100)

	assert.Equal(t, err, io.EOF)
	assert.Equal(t, buf.Len(), 2) // buffer doesn't advance
}

func TestReadBufferInspect(t *testing.T) {
	buf := ReadBuffer{0xDE, 0xAD, 0xBE, 0xEF, 0xFE, 0xE1, 0xDE, 0xAD}
	output := buf.Inspect()

	assert.Contains(t, output, "00000000  de ad be ef fe e1 de ad")
}
