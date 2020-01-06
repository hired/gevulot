package server

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServerStart(t *testing.T) {
	configChan := make(chan *Config, 1)
	defer close(configChan)

	configPublisher := NewConfigDistributor(configChan)
	defer configPublisher.Close()

	configChan <- &Config{Listen: "0.0.0.0:4242"}

	serverIsListening := make(chan bool)
	defer close(serverIsListening)

	serverErr := make(chan error)
	defer close(serverErr)

	srv := NewServer(configPublisher)
	srv.testHookServe = func(net.Listener) {
		serverIsListening <- true
	}

	defer srv.Close()

	go func() {
		serverErr <- srv.Start()
	}()

	<-serverIsListening

	// Server is listening with initial config
	assert.NotNil(t, srv.listener)
	assert.Equal(t, "[::]:4242", srv.listener.Addr().String())

	configChan <- &Config{Listen: "0.0.0.0:31337"}

	<-serverIsListening

	// Server is listening with updated config
	assert.NotNil(t, srv.listener)
	assert.Equal(t, "[::]:31337", srv.listener.Addr().String())

	// Start returns ErrServerAlreadyStarted
	err := srv.Start()
	assert.Equal(t, ErrServerAlreadyStarted, err)

	// Close stops the Start loop
	closeErr := srv.Close()
	assert.NoError(t, closeErr)
	assert.Equal(t, ErrServerClosed, <-serverErr)

	// Start returns ErrServerClosed after Close
	err = srv.Start()
	assert.Equal(t, ErrServerClosed, err)
}

func TestServerServe(t *testing.T) {
	serverIsListening := make(chan bool)
	defer close(serverIsListening)

	serverSessions := make(chan *Session)
	defer close(serverSessions)

	serveError := make(chan error)
	defer close(serveError)

	srv := NewServer(nil)
	srv.testHookServe = func(net.Listener) {
		serverIsListening <- true
	}
	srv.testHookServeConn = func(sess *Session) {
		serverSessions <- sess
	}

	defer srv.Close()

	l, err := net.Listen("tcp", "127.0.0.1:0")
	assert.NoError(t, err)

	defer l.Close()

	go func() {
		serveError <- srv.Serve(l)
	}()

	<-serverIsListening

	// Serve sets the Server listener
	assert.Same(t, l, srv.listener)

	client, err := net.Dial("tcp", l.Addr().String())
	assert.NoError(t, err)

	defer client.Close()

	// Server accepting client connections
	assert.Equal(t, client.LocalAddr(), (<-serverSessions).clientConn.Unwrap().RemoteAddr())

	// Close stops the Serve loop
	err = srv.Close()
	assert.NoError(t, err)
	assert.Contains(t, (<-serveError).Error(), "use of closed network connection")

	// Serve returns ErrServerClosed after Close
	err = srv.Serve(l)
	assert.Equal(t, ErrServerClosed, err)
}

func TestServerServeConn(t *testing.T) {
	serverSessions := make(chan *Session)
	defer close(serverSessions)

	srv := NewServer(nil)
	srv.testHookServeConn = func(sess *Session) {
		serverSessions <- sess
	}

	defer srv.Close()

	in, out := net.Pipe()

	defer in.Close()
	defer out.Close()

	go func() {
		_ = srv.ServeConn(out)
	}()

	session := <-serverSessions

	// Server registers session
	assert.Len(t, srv.sessions, 1)
	assert.Contains(t, srv.sessions, session)

	err := srv.Close()

	// Server unregister session after Close
	assert.NoError(t, err)
	assert.Len(t, srv.sessions, 0)

	// ServeConn returns ErrServerClosed after Close
	err = srv.ServeConn(out)
	assert.Equal(t, ErrServerClosed, err)
}
