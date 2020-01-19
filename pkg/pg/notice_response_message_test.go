package pg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Real notice message packet from pg captured with Wireshark
//   Severity (localized): WARNING
//   Severity: WARNING
//   Code: 01000
//   Message: GLOBAL is deprecated in temporary table creation
//   Position: 8
//   File: gram.y
//   Line: 3279
//   Routine: base_yyparse
const GoldenNoticeMessagePacket = "\x4e\x00\x00\x00\x6f\x53\x57\x41\x52\x4e\x49\x4e\x47\x00\x56\x57" +
	"\x41\x52\x4e\x49\x4e\x47\x00\x43\x30\x31\x30\x30\x30\x00\x4d\x47" +
	"\x4c\x4f\x42\x41\x4c\x20\x69\x73\x20\x64\x65\x70\x72\x65\x63\x61" +
	"\x74\x65\x64\x20\x69\x6e\x20\x74\x65\x6d\x70\x6f\x72\x61\x72\x79" +
	"\x20\x74\x61\x62\x6c\x65\x20\x63\x72\x65\x61\x74\x69\x6f\x6e\x00" +
	"\x50\x38\x00\x46\x67\x72\x61\x6d\x2e\x79\x00\x4c\x33\x32\x37\x39" +
	"\x00\x52\x62\x61\x73\x65\x5f\x79\x79\x70\x61\x72\x73\x65\x00\x00"

func TestParseNoticeResponseMessage(t *testing.T) {
	{
		msg, err := ParseNoticeResponseMessage(StandardFrame(GoldenNoticeMessagePacket))

		assert.NoError(t, err)

		expectedFields := []*MessageField{
			{MessageFieldSeverityLocalized, "WARNING"},
			{MessageFieldSeverity, "WARNING"},
			{MessageFieldCode, "01000"},
			{MessageFieldMessage, "GLOBAL is deprecated in temporary table creation"},
			{MessageFieldPosition, "8"},
			{MessageFieldFile, "gram.y"},
			{MessageFieldLine, "3279"},
			{MessageFieldRoutine, "base_yyparse"},
		}

		assert.Equal(t, expectedFields, msg.Fields)
	}

	// Test invalid type
	{
		_, err := ParseNoticeResponseMessage(append(StandardFrame{'X'}, GoldenNoticeMessagePacket[1:]...))
		assert.Equal(t, ErrMalformedMessage, err)
	}
}

func TestParseNoticeResponseMessageFrame(t *testing.T) {
	msg := &NoticeResponseMessage{
		Fields: []*MessageField{
			{MessageFieldSeverityLocalized, "WARNING"},
			{MessageFieldSeverity, "WARNING"},
			{MessageFieldCode, "01000"},
			{MessageFieldMessage, "GLOBAL is deprecated in temporary table creation"},
			{MessageFieldPosition, "8"},
			{MessageFieldFile, "gram.y"},
			{MessageFieldLine, "3279"},
			{MessageFieldRoutine, "base_yyparse"},
		},
	}

	assert.Equal(t, []byte(GoldenNoticeMessagePacket), msg.Frame().Bytes())
}
