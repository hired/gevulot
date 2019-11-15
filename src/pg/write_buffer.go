package pg

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

// WriteBuffer provides bufio.Scanner-like API to encode PG data types described here:
// https://www.postgresql.org/docs/9.6/protocol-message-types.html
type WriteBuffer []byte

//
// Main use case:
//
//   var buf PgReadBuffer     // No need to allocate if you don't know the size
//   buf.WriteInt32()         // Write values
//   buf.WiteString("hello")
//   ...
//

// Len returns the number of bytes in the buffer.
func (b *WriteBuffer) Len() int {
	return len(*b)
}

// WriteBytes appends given bytes to the buffer.
func (b *WriteBuffer) WriteBytes(v []byte) {
	*b = append(*b, v...)
}

// WriteByte appends given byte to the buffer.
func (b *WriteBuffer) WriteByte(c byte) {
	*b = append(*b, c)
}

// WriteInt32 encodes a 32-bit signed integer and appends it to the buffer.
func (b *WriteBuffer) WriteInt32(num int32) {
	numBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(numBytes, uint32(num))
	b.WriteBytes(numBytes)
}

// WriteInt16 encodes a 16-bit signed integer and appends it to the buffer.
func (b *WriteBuffer) WriteInt16(num int16) {
	numBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(numBytes, uint16(num))
	b.WriteBytes(numBytes)
}

// WriteString encodes and appends given string to the buffer.
func (b *WriteBuffer) WriteString(str string) {
	*b = append(append(*b, str...), '\000')
}

// WriteInt32Array encodes an array of 32-bit signed integers and append it into the buffer.
func (b *WriteBuffer) WriteInt32Array(arr []int32) {
	arrBytes := make([]byte, len(arr)*4)

	for i := range arr {
		binary.BigEndian.PutUint32(arrBytes[i*4:i*4+4], uint32(arr[i]))
	}

	b.WriteBytes(arrBytes)
}

// WriteInt16Array encodes an array of 16-bit signed integers and append it into the buffer.
func (b *WriteBuffer) WriteInt16Array(arr []int16) {
	arrBytes := make([]byte, len(arr)*2)

	for i := range arr {
		binary.BigEndian.PutUint16(arrBytes[i*2:i*2+2], uint16(arr[i]))
	}

	b.WriteBytes(arrBytes)
}

// Inspect returns content of the buffer as a hex dump for debug purposes.
func (b *WriteBuffer) Inspect() string {
	return fmt.Sprintf("buffer length: %d\n%s", b.Len(), hex.Dump(*b))
}
