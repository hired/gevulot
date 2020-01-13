package pg

import (
	"io"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnRecvStartupMessage(t *testing.T) {
	client, server := net.Pipe()

	defer client.Close()
	defer server.Close()

	go func() {
		_, err := server.Write([]byte(GoldenStartupMessagePacket))

		if err != nil {
			panic(err)
		}
	}()

	pgConn := NewConn(client)
	msg, err := pgConn.RecvStartupMessage()

	assert.NoError(t, err)
	assert.Equal(t, int32(DefaultProtocolVersion), msg.ProtocolVersion)
}

func TestConnRecvMessage(t *testing.T) {
	client, server := net.Pipe()

	defer client.Close()
	defer server.Close()

	go func() {
		msg := &GenericMessage{
			Type: '$',
			Body: []byte("test test"),
		}

		_, err := server.Write(msg.Frame().Bytes())

		if err != nil {
			panic(err)
		}
	}()

	pgConn := NewConn(client)
	msg, err := pgConn.RecvMessage()

	assert.NoError(t, err)

	if gm, ok := msg.(*GenericMessage); ok {
		assert.Equal(t, byte('$'), gm.Type)
		assert.Equal(t, []byte("test test"), gm.Body)
	} else {
		assert.Fail(t, "received message is not a GenericMessage")
	}
}

func TestConnSendMessage(t *testing.T) {
	client, server := net.Pipe()

	defer client.Close()
	defer server.Close()

	msg := &GenericMessage{
		Type: '$',
		Body: []byte("test test"),
	}

	go func() {
		pgConn := NewConn(client)
		err := pgConn.SendMessage(msg)

		if err != nil {
			panic(err)
		}
	}()

	expectedFrame := msg.Frame()
	receivedFrame, err := ReadStandardFrame(server)

	assert.NoError(t, err)
	assert.Equal(t, expectedFrame, receivedFrame)
}

func TestConnSendByte(t *testing.T) {
	client, server := net.Pipe()

	defer client.Close()
	defer server.Close()

	go func() {
		pgConn := NewConn(client)
		err := pgConn.SendByte('$')

		if err != nil {
			panic(err)
		}
	}()

	buf := make([]byte, 1)
	_, err := server.Read(buf)

	assert.NoError(t, err)
	assert.Equal(t, []byte{'$'}, buf)
}

func TestConnUnwrap(t *testing.T) {
	client, server := net.Pipe()

	defer client.Close()
	defer server.Close()

	pgConn := NewConn(client)

	assert.Same(t, client, pgConn.Unwrap())
}

func TestConnClose(t *testing.T) {
	client, server := net.Pipe()

	defer client.Close()
	defer server.Close()

	pgConn := NewConn(client)
	pgConn.Close()

	_, err := client.Read(make([]byte, 1))

	assert.Equal(t, io.ErrClosedPipe, err)
}

func BenchmarkConnThroughput(b *testing.B) {
	client, server := net.Pipe()

	defer client.Close()
	defer server.Close()

	b.ResetTimer()

	go func() {
		pgConn := NewConn(server)

		for n := 0; n < b.N; n++ {
			// FIXME: use real message
			msg := &GenericMessage{
				Type: 'Q',
				Body: []byte("SELECT * FROM users"),
			}

			err := pgConn.SendMessage(msg)

			if err != nil {
				panic(err)
			}
		}
	}()

	pgConn := NewConn(client)

	for n := 0; n < b.N; n++ {
		_, err := pgConn.RecvMessage()

		if err != nil {
			b.Fatal(err)
		}
	}
}
