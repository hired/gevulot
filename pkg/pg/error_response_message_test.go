package pg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Real error message packet from pg captured with Wireshark
//   Severity (localized): FATAL
//   Severity: FATAL
//   Code: 28P01
//   Message: password authentication failed for user "chrno"
//   File: auth.c
//   Line: 333
//   Routine: auth_failed
const GoldenErrorMessagePacket = "\x45\x00\x00\x00\x65\x53\x46\x41\x54\x41\x4c\x00\x56\x46\x41\x54" +
	"\x41\x4c\x00\x43\x32\x38\x50\x30\x31\x00\x4d\x70\x61\x73\x73\x77" +
	"\x6f\x72\x64\x20\x61\x75\x74\x68\x65\x6e\x74\x69\x63\x61\x74\x69" +
	"\x6f\x6e\x20\x66\x61\x69\x6c\x65\x64\x20\x66\x6f\x72\x20\x75\x73" +
	"\x65\x72\x20\x22\x63\x68\x72\x6e\x6f\x22\x00\x46\x61\x75\x74\x68" +
	"\x2e\x63\x00\x4c\x33\x33\x33\x00\x52\x61\x75\x74\x68\x5f\x66\x61" +
	"\x69\x6c\x65\x64\x00\x00"

func TestParseErrorResponseMessage(t *testing.T) {
	{
		msg, err := ParseErrorResponseMessage(StandardFrame(GoldenErrorMessagePacket))

		assert.NoError(t, err)

		expectedFields := []*MessageField{
			{MessageFieldSeverityLocalized, "FATAL"},
			{MessageFieldSeverity, "FATAL"},
			{MessageFieldCode, "28P01"},
			{MessageFieldMessage, `password authentication failed for user "chrno"`},
			{MessageFieldFile, "auth.c"},
			{MessageFieldLine, "333"},
			{MessageFieldRoutine, "auth_failed"},
		}

		assert.Equal(t, expectedFields, msg.Fields)
	}

	// Test invalid type
	{
		_, err := ParseErrorResponseMessage(append(StandardFrame{'X'}, GoldenErrorMessagePacket[1:]...))
		assert.Equal(t, ErrMalformedMessage, err)
	}
}

func TestParseErrorResponseMessageFrame(t *testing.T) {
	msg := &ErrorResponseMessage{
		Fields: []*MessageField{
			{MessageFieldSeverityLocalized, "FATAL"},
			{MessageFieldSeverity, "FATAL"},
			{MessageFieldCode, "28P01"},
			{MessageFieldMessage, `password authentication failed for user "chrno"`},
			{MessageFieldFile, "auth.c"},
			{MessageFieldLine, "333"},
			{MessageFieldRoutine, "auth_failed"},
		},
	}

	assert.Equal(t, []byte(GoldenErrorMessagePacket), msg.Frame().Bytes())
}
