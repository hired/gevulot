package pg

import (
	"testing"

	"gotest.tools/assert"
	is "gotest.tools/assert/cmp"
)

func TestWriteBufferLen(t *testing.T) {
	buf := make(WriteBuffer, 42)
	assert.Equal(t, buf.Len(), 42)
}

func TestWriteBufferWriteBytes(t *testing.T) {
	var buf WriteBuffer

	assert.Equal(t, buf.Len(), 0)

	err := buf.WriteBytes([]byte{0xDE, 0xAD, 0xBE, 0xEF})

	assert.NilError(t, err)
	assert.DeepEqual(t, buf, WriteBuffer{0xDE, 0xAD, 0xBE, 0xEF})

	err = buf.WriteBytes([]byte{0xFE, 0xE1, 0xDE, 0xAD})

	assert.NilError(t, err)
	assert.DeepEqual(t, buf, WriteBuffer{0xDE, 0xAD, 0xBE, 0xEF, 0xFE, 0xE1, 0xDE, 0xAD})
}

func TestWriteBufferWriteByte(t *testing.T) {
	var buf WriteBuffer

	assert.Equal(t, buf.Len(), 0)

	err := buf.WriteByte('A')

	assert.NilError(t, err)
	assert.DeepEqual(t, buf, WriteBuffer{'A'})

	err = buf.WriteByte('B')

	assert.NilError(t, err)
	assert.DeepEqual(t, buf, WriteBuffer{'A', 'B'})
}

func TestWriteBufferWriteInt32(t *testing.T) {
	var buf WriteBuffer

	assert.Equal(t, buf.Len(), 0)

	err := buf.WriteInt32(31337)

	assert.NilError(t, err)
	assert.DeepEqual(t, buf, WriteBuffer{0x00, 0x00, 0x7A, 0x69})

	err = buf.WriteInt32(42)

	assert.NilError(t, err)
	assert.DeepEqual(t, buf, WriteBuffer{
		0x00, 0x00, 0x7A, 0x69, // 31337 in Big-endian
		0x00, 0x00, 0x00, 0x2A, // 42 in Big-endian
	})
}

func TestWriteBufferWriteInt16(t *testing.T) {
	var buf WriteBuffer

	assert.Equal(t, buf.Len(), 0)

	err := buf.WriteInt16(31337)

	assert.NilError(t, err)
	assert.DeepEqual(t, buf, WriteBuffer{0x7A, 0x69})

	err = buf.WriteInt16(42)

	assert.NilError(t, err)
	assert.DeepEqual(t, buf, WriteBuffer{
		0x7A, 0x69, // 31337 in Big-endian
		0x00, 0x2A, // 42 in Big-endian
	})
}

func TestWriteBufferWriteString(t *testing.T) {
	var buf WriteBuffer

	assert.Equal(t, buf.Len(), 0)

	err := buf.WriteString("hello")

	assert.NilError(t, err)
	assert.DeepEqual(t, buf, WriteBuffer{'h', 'e', 'l', 'l', 'o', 0x00})

	err = buf.WriteString("world")

	assert.NilError(t, err)
	assert.DeepEqual(t, buf, WriteBuffer{
		'h', 'e', 'l', 'l', 'o', 0x00,
		'w', 'o', 'r', 'l', 'd', 0x00,
	})
}

func TestWriteBufferWriteInt32Array(t *testing.T) {
	var buf WriteBuffer

	assert.Equal(t, buf.Len(), 0)

	err := buf.WriteInt32Array([]int32{31337, 42})

	assert.NilError(t, err)
	assert.DeepEqual(t, buf, WriteBuffer{
		0x00, 0x00, 0x7A, 0x69, // 31337 in Big-endian
		0x00, 0x00, 0x00, 0x2A, // 42 in Big-endian
	})

	err = buf.WriteInt32Array([]int32{42, 31337})

	assert.NilError(t, err)
	assert.DeepEqual(t, buf, WriteBuffer{
		0x00, 0x00, 0x7A, 0x69,
		0x00, 0x00, 0x00, 0x2A,
		0x00, 0x00, 0x00, 0x2A,
		0x00, 0x00, 0x7A, 0x69,
	})
}

func TestWriteBufferWriteInt16Array(t *testing.T) {
	var buf WriteBuffer

	assert.Equal(t, buf.Len(), 0)

	err := buf.WriteInt16Array([]int16{31337, 42})

	assert.NilError(t, err)
	assert.DeepEqual(t, buf, WriteBuffer{
		0x7A, 0x69, // 31337 in Big-endian
		0x00, 0x2A, // 42 in Big-endian
	})

	err = buf.WriteInt16Array([]int16{42, 31337})

	assert.NilError(t, err)
	assert.DeepEqual(t, buf, WriteBuffer{
		0x7A, 0x69,
		0x00, 0x2A,
		0x00, 0x2A,
		0x7A, 0x69,
	})
}

func TestWriteBufferInspect(t *testing.T) {
	buf := WriteBuffer{0xDE, 0xAD, 0xBE, 0xEF, 0xFE, 0xE1, 0xDE, 0xAD}
	output := buf.Inspect()

	assert.Assert(t, is.Contains(output, "00000000  de ad be ef fe e1 de ad"))
}
