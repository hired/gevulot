package pg

import (
	"testing"

	"gotest.tools/assert"
)

// Real SSL request packet from psql captured with Wireshark
const GoldenSSLRequestMessagePacket = "\x00\x00\x00\x08\x04\xd2\x16\x2f"

// Real startup message packet from psql captured with Wireshark
// Parameters:
//   user = hired
//   database = hired_dev
//   application_name = psql
//   client_encoding = UTF8
const GoldenStartupMessagePacket = "\x00\x00\x00\x52\x00\x03\x00\x00\x75\x73\x65\x72\x00\x68\x69\x72" +
	"\x65\x64\x00\x64\x61\x74\x61\x62\x61\x73\x65\x00\x68\x69\x72\x65" +
	"\x64\x5f\x64\x65\x76\x00\x61\x70\x70\x6c\x69\x63\x61\x74\x69\x6f" +
	"\x6e\x5f\x6e\x61\x6d\x65\x00\x70\x73\x71\x6c\x00\x63\x6c\x69\x65" +
	"\x6e\x74\x5f\x65\x6e\x63\x6f\x64\x69\x6e\x67\x00\x55\x54\x46\x38" +
	"\x00\x00"

func TestStartupMessageMarshall(t *testing.T) {
	// Test SSL request
	msg := &StartupMessage{
		ProtocolVersion: SSLRequestMagic,
	}

	packet, err := msg.Marshall()

	assert.NilError(t, err)
	assert.DeepEqual(t, packet, []byte(GoldenSSLRequestMessagePacket))

	// Test regular start up
	msg = &StartupMessage{
		ProtocolVersion: DefaultProtocolVersion,
		Parameters: []StartupMessageParameter{
			{Name: "user", Value: "hired"},
			{Name: "database", Value: "hired_dev"},
			{Name: "application_name", Value: "psql"},
			{Name: "client_encoding", Value: "UTF8"},
		},
	}

	packet, err = msg.Marshall()

	assert.NilError(t, err)
	assert.DeepEqual(t, packet, []byte(GoldenStartupMessagePacket))
}
