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

	buf.WriteBytes([]byte{0xDE, 0xAD, 0xBE, 0xEF})
	assert.DeepEqual(t, buf, WriteBuffer{0xDE, 0xAD, 0xBE, 0xEF})

	buf.WriteBytes([]byte{0xFE, 0xE1, 0xDE, 0xAD})
	assert.DeepEqual(t, buf, WriteBuffer{0xDE, 0xAD, 0xBE, 0xEF, 0xFE, 0xE1, 0xDE, 0xAD})
}

func TestWriteBufferWriteByte(t *testing.T) {
	var buf WriteBuffer

	assert.Equal(t, buf.Len(), 0)

	buf.WriteByte('A')
	assert.DeepEqual(t, buf, WriteBuffer{'A'})

	buf.WriteByte('B')
	assert.DeepEqual(t, buf, WriteBuffer{'A', 'B'})
}

func TestWriteBufferWriteInt32(t *testing.T) {
	var buf WriteBuffer

	assert.Equal(t, buf.Len(), 0)

	buf.WriteInt32(31337)
	assert.DeepEqual(t, buf, WriteBuffer{0x00, 0x00, 0x7A, 0x69})

	buf.WriteInt32(42)
	assert.DeepEqual(t, buf, WriteBuffer{
		0x00, 0x00, 0x7A, 0x69, // 31337 in Big-endian
		0x00, 0x00, 0x00, 0x2A, // 42 in Big-endian
	})
}

func TestWriteBufferWriteInt16(t *testing.T) {
	var buf WriteBuffer

	assert.Equal(t, buf.Len(), 0)

	buf.WriteInt16(31337)
	assert.DeepEqual(t, buf, WriteBuffer{0x7A, 0x69})

	buf.WriteInt16(42)
	assert.DeepEqual(t, buf, WriteBuffer{
		0x7A, 0x69, // 31337 in Big-endian
		0x00, 0x2A, // 42 in Big-endian
	})
}

func TestWriteBufferWriteString(t *testing.T) {
	var buf WriteBuffer

	assert.Equal(t, buf.Len(), 0)

	buf.WriteString("hello")
	assert.DeepEqual(t, buf, WriteBuffer{'h', 'e', 'l', 'l', 'o', 0x00})

	buf.WriteString("world")
	assert.DeepEqual(t, buf, WriteBuffer{
		'h', 'e', 'l', 'l', 'o', 0x00,
		'w', 'o', 'r', 'l', 'd', 0x00,
	})
}

func TestWriteBufferWriteInt32Array(t *testing.T) {
	var buf WriteBuffer

	assert.Equal(t, buf.Len(), 0)

	buf.WriteInt32Array([]int32{31337, 42})
	assert.DeepEqual(t, buf, WriteBuffer{
		0x00, 0x00, 0x7A, 0x69, // 31337 in Big-endian
		0x00, 0x00, 0x00, 0x2A, // 42 in Big-endian
	})

	buf.WriteInt32Array([]int32{42, 31337})
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

	buf.WriteInt16Array([]int16{31337, 42})
	assert.DeepEqual(t, buf, WriteBuffer{
		0x7A, 0x69, // 31337 in Big-endian
		0x00, 0x2A, // 42 in Big-endian
	})

	buf.WriteInt16Array([]int16{42, 31337})
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
