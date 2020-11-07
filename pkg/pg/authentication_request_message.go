package pg

import (
	"errors"
	"fmt"
	"reflect"
)

// AuthenticationRequestMessageType identifies authentication request message.
const AuthenticationRequestMessageType = 'R'

// AuthenticationRequestMessage is one of backend authentication requests.
type AuthenticationRequestMessage interface {
	Message
}

// AuthenticationOkMessage is sent by a backend when the authentication was successful.
type AuthenticationOkMessage struct{}

// AuthenticationKerberosV5Message is sent by a backend when Kerberos V5 authentication is required.
type AuthenticationKerberosV5Message struct{}

// AuthenticationCleartextPasswordMessage is sent by a backend when a clear-text password is required.
type AuthenticationCleartextPasswordMessage struct{}

// AuthenticationMD5PasswordMessage is sent by a backend when an MD5-encrypted password is required.
type AuthenticationMD5PasswordMessage struct {
	// The salt to use when encrypting the password.
	Salt [4]byte
}

// AuthenticationSCMCredentialMessage is sent by a backend when an SCM credentials message is required.
type AuthenticationSCMCredentialMessage struct{}

// AuthenticationGSSMessage is sent by a backend when GSSAPI authentication is required.
type AuthenticationGSSMessage struct{}

// AuthenticationGSSContinueMessage is sent by a frontend to authenticate with GSSAPI or SSPI.
type AuthenticationGSSContinueMessage struct {
	// GSSAPI or SSPI authentication data.
	Data []byte
}

// AuthenticationSSPIMessage is sent by a backend when SSPI authentication is required.
type AuthenticationSSPIMessage struct{}

var (
	// ErrunsupportedAuthenticationRequest is returned when ParseAuthenticationRequestMessage
	// cannot handle unknown auth status.
	ErrunsupportedAuthenticationRequest = errors.New("pg: unsupported auth request from a backend")
)

// All valid auth status codes.
const (
	authStatusCodeOk                = 0
	authStatusCodeKerberosV5        = 2
	authStatusCodeCleartextPassword = 3
	authStatusCodeMD5Password       = 5
	authStatusCodeSCMCredential     = 6
	authStatusCodeGSS               = 7
	authStatusCodeGSSContinue       = 8
	authStatusCodeSSPI              = 9
)

// Mapping between status code and implementation of AuthenticationRequestMessage.
var authStatusMap = map[int32]reflect.Type{ //nolint:gocheckboglobals
	authStatusCodeOk:                reflect.TypeOf(AuthenticationOkMessage{}),
	authStatusCodeKerberosV5:        reflect.TypeOf(AuthenticationKerberosV5Message{}),
	authStatusCodeCleartextPassword: reflect.TypeOf(AuthenticationCleartextPasswordMessage{}),
	authStatusCodeMD5Password:       reflect.TypeOf(AuthenticationMD5PasswordMessage{}),
	authStatusCodeSCMCredential:     reflect.TypeOf(AuthenticationSCMCredentialMessage{}),
	authStatusCodeGSS:               reflect.TypeOf(AuthenticationGSSMessage{}),
	authStatusCodeGSSContinue:       reflect.TypeOf(AuthenticationGSSContinueMessage{}),
	authStatusCodeSSPI:              reflect.TypeOf(AuthenticationSSPIMessage{}),
}

// ParseAuthenticationRequestMessage parses authentication request from a network frame.
func ParseAuthenticationRequestMessage(frame Frame) (AuthenticationRequestMessage, error) {
	// Assert the message type
	if frame.MessageType() != AuthenticationRequestMessageType {
		return nil, ErrMalformedMessage
	}

	messageData := ReadBuffer(frame.MessageBody())

	// Read the status cide
	status, err := messageData.ReadInt32()

	if err != nil {
		return nil, err
	}

	// Resolve concrete message struct type
	messageType := authStatusMap[status]

	// Unknown status code. Probably newer protocol?
	if messageType == nil {
		return nil, ErrunsupportedAuthenticationRequest
	}

	// Initialize message struct from the type
	message, ok := reflect.New(messageType).Interface().(AuthenticationRequestMessage)

	if !ok {
		// We failed to typecast message to the AuthenticationRequestMessage; this cannot happen under any circumstances
		panic(fmt.Sprintf("pg: ParseAuthenticationRequestMessage: error constructing auth message %#v", messageType))
	}

	// Some message types contain additional fields — parse them
	switch v := message.(type) {
	case *AuthenticationMD5PasswordMessage:
		salt, err := messageData.ReadBytes(4)

		if err != nil {
			return nil, err
		}

		copy(v.Salt[:], salt)

	case *AuthenticationGSSContinueMessage:
		v.Data, err = messageData.ReadBytes(messageData.Len())

		if err != nil {
			return nil, err
		}
	}

	return message, nil
}

// Frame serializes the message into a network frame.
func (m *AuthenticationOkMessage) Frame() Frame {
	var messageBuffer WriteBuffer
	messageBuffer.WriteInt32(authStatusCodeOk)

	return NewStandardFrame(AuthenticationRequestMessageType, messageBuffer)
}

// Frame serializes the message into a network frame.
func (m *AuthenticationKerberosV5Message) Frame() Frame {
	var messageBuffer WriteBuffer
	messageBuffer.WriteInt32(authStatusCodeKerberosV5)

	return NewStandardFrame(AuthenticationRequestMessageType, messageBuffer)
}

// Frame serializes the message into a network frame.
func (m *AuthenticationCleartextPasswordMessage) Frame() Frame {
	var messageBuffer WriteBuffer
	messageBuffer.WriteInt32(authStatusCodeCleartextPassword)

	return NewStandardFrame(AuthenticationRequestMessageType, messageBuffer)
}

// Frame serializes the message into a network frame.
func (m *AuthenticationMD5PasswordMessage) Frame() Frame {
	var messageBuffer WriteBuffer

	messageBuffer.WriteInt32(authStatusCodeMD5Password)
	messageBuffer.WriteBytes(m.Salt[:])

	return NewStandardFrame(AuthenticationRequestMessageType, messageBuffer)
}

// Frame serializes the message into a network frame.
func (m *AuthenticationSCMCredentialMessage) Frame() Frame {
	var messageBuffer WriteBuffer
	messageBuffer.WriteInt32(authStatusCodeSCMCredential)

	return NewStandardFrame(AuthenticationRequestMessageType, messageBuffer)
}

// Frame serializes the message into a network frame.
func (m *AuthenticationGSSMessage) Frame() Frame {
	var messageBuffer WriteBuffer
	messageBuffer.WriteInt32(authStatusCodeGSS)

	return NewStandardFrame(AuthenticationRequestMessageType, messageBuffer)
}

// Frame serializes the message into a network frame.
func (m *AuthenticationGSSContinueMessage) Frame() Frame {
	var messageBuffer WriteBuffer

	messageBuffer.WriteInt32(authStatusCodeGSSContinue)
	messageBuffer.WriteBytes(m.Data)

	return NewStandardFrame(AuthenticationRequestMessageType, messageBuffer)
}

// Frame serializes the message into a network frame.
func (m *AuthenticationSSPIMessage) Frame() Frame {
	var messageBuffer WriteBuffer
	messageBuffer.WriteInt32(authStatusCodeSSPI)

	return NewStandardFrame(AuthenticationRequestMessageType, messageBuffer)
}
