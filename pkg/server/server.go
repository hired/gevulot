package server

import (
	"errors"
	"net"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	// ErrServerClosed is returned by the Server's Start, Serve and ServeConn methods after a call to Close.
	ErrServerClosed = errors.New("server: Server closed")

	// ErrServerAlreadyStarted is returned by the Server's Start if its already have been called.
	ErrServerAlreadyStarted = errors.New("server: Server already started")
)

// Server is a masking PostgreSQL proxy server.
type Server struct {
	// Configuration provider
	config ConfigStore

	// We use this waitgroup to track goroutines that server is creating.
	// Close() waits for all of them to finish.
	wg sync.WaitGroup

	// Guards listener and sessions
	mu sync.Mutex

	// Fired when Start is called
	start *Event

	// Fired when Close is called
	shutdown *Event

	// Current active listener (use changeListener to change)
	listener net.Listener

	// List of currently active database sessions
	sessions map[*Session]struct{}

	// When set, called after Serve successfully set a new listener
	// but before is started to accept client connections
	testHookServe func(net.Listener)

	// When set, called after ServeConn is ready to handle client connection.
	testHookServeConn func(*Session)
}

// NewServer initializes a new Server instance.
func NewServer(config ConfigStore) *Server {
	return &Server{
		config:   config,
		start:    NewEvent(),
		shutdown: NewEvent(),
	}
}

// Start listens on the TCP network address specified in the Server's
// config and then calls Serve to handle incoming connections from the clients
// to the proxied database.
//
// When Server's configuration changed (e.g., when listen port has been modified
// in the configuration file), Start automatically changes the listener.
//
// Start can only be called once per Server instance.
//
// Start always returns a non-nil error.
// After Close, the returned error is ErrServerClosed;
// after consequent Start, the returned error is ErrServerAlreadyStarted.
func (srv *Server) Start() error {
	// Return error if server is closed
	if srv.shutdown.HasFired() {
		return ErrServerClosed
	}

	// Make sure that Start is only called once
	if !srv.start.Fire() {
		return ErrServerAlreadyStarted
	}

	// Watch for config changes
	serverConfigurationUpdates := make(chan *Config, 1)
	defer close(serverConfigurationUpdates)

	err := srv.config.Subscribe(serverConfigurationUpdates, func(oldConfig, newConfig *Config) bool {
		return oldConfig == nil || oldConfig.Listen != newConfig.Listen
	})

	if err != nil {
		return err
	}

	// For debugging purposes
	defer log.Debug("server: Start() loop finished")

	for {
		select {
		// Wait for the new config
		case config := <-serverConfigurationUpdates:
			log.Infof("server: serving on %s", config.Listen)

			// Initialize a new listener
			ln, err := net.Listen("tcp", config.Listen)

			if err != nil {
				log.Errorf("server: can't listen on %s: %v", config.Listen, err)
				continue
			}

			// Serve the new listener
			srv.do(func() { _ = srv.Serve(ln) })

		// Wait for the server shutdown
		case <-srv.shutdown.Done():
			return ErrServerClosed
		}
	}
}

// Serve sets the Server listener (closing existing one if set) and then
// calls ServeConn for every accepted client connection.
//
// Serve blocks until the listener returns a non-nil error. The caller typically
// invokes Serve in a go statement.
func (srv *Server) Serve(ln net.Listener) error {
	// Update connection; the err could be ErrServerClosed
	err := srv.changeListener(ln)

	if err != nil {
		log.Errorf("server: error has been occurred while changing listeners: %v", err)
		return err
	}

	log.Info("server: ready to accept client connections")

	// Notify tests that server is listening
	if srv.testHookServe != nil {
		srv.testHookServe(ln)
	}

	// For debugging purposes
	defer log.Debug("server: Serve() loop finished")

	for {
		// NB: we use l not srv.listener here because later can be nil
		conn, err := ln.Accept()

		if err != nil {
			// Log error unless it is a closed listener
			if !strings.Contains(err.Error(), "use of closed network connection") {
				log.Errorf("server: error while accepting client connection: %v", err)
			}

			// Continue after temporary errors
			var ne net.Error
			if errors.As(err, &ne) && ne.Temporary() {
				continue
			}

			// ...otherwise stop the loop (most likely listener is closed)
			return err
		}

		// Serve client in a new goroutine
		srv.do(func() { _ = srv.ServeConn(conn) })
	}
}

// ServeConn proxies given connection to the PostgreSQL database specified in the Server's config.
//
// ServeConn blocks, serving the connection until the client or the database hangs up.
// The caller typically invokes ServeConn in a go statement.
func (srv *Server) ServeConn(conn net.Conn) error {
	// Return error if server is closed
	if srv.shutdown.HasFired() {
		return ErrServerClosed
	}

	// For debugging purposes
	defer log.Debugf("serve: ServeConn() for %s exited", conn.RemoteAddr().String())

	log.Infof("server: new client connection from %s", conn.RemoteAddr().String())

	// Initialize a new session
	session := NewSession(conn, srv.config)

	// Register session in the list of active server sessions; the err could be ErrServerClosed
	err := srv.registerSession(session)

	if err != nil {
		return err
	}

	// Automatically remove session from the list
	defer srv.removeSession(session)

	// Notify tests that server is about to handle client connection
	if srv.testHookServeConn != nil {
		srv.testHookServeConn(session)
	}

	// Start the proxied DB session
	return session.Start()
}

// Close immediately closes the Server's listener and all active sessions.
// Close returns any error returned from closing the listener.
//
// Once Close has been called on a server, it may not be reused;
// future calls to methods such as Serve or Start will return ErrServerClosed.
func (srv *Server) Close() error {
	// Do nothing if Close has been called already
	if !srv.shutdown.Fire() {
		return nil
	}

	log.Info("server: closing")

	// For debugging purposes
	defer log.Debug("server: Close() exited")

	// Error to be returned from Close
	var resultErr error

	srv.mu.Lock()
	{
		// Close the listener
		if srv.listener != nil {
			resultErr = srv.listener.Close()
			srv.listener = nil
		}

		// Close every active session
		for s := range srv.sessions {
			err := s.Close()

			if err != nil {
				log.Errorf("server: error closing session: %v", err)
			}
		}

		srv.sessions = nil
	}
	srv.mu.Unlock()

	// Run a watchdog in background that will panic if there are any running goroutines
	// left after we closed everything.
	watchdog := time.AfterFunc(time.Second*5, func() {
		panic("server: some goroutines are still running after server is closed")
	})

	// Wait for Server's goroutines to finish
	srv.wg.Wait()

	// All good — no need to panic
	watchdog.Stop()

	return resultErr
}

// changeListener sets the new Server's listener closing existing one. It can be called concurrently.
func (srv *Server) changeListener(ln net.Listener) error {
	// Refuse to change a listener if Server is closed
	if srv.shutdown.HasFired() {
		return ErrServerClosed
	}

	srv.mu.Lock()
	defer srv.mu.Unlock()

	if srv.listener != nil {
		log.Debug("server: closing the old listener")

		err := srv.listener.Close()

		if err != nil {
			log.Errorf("server: error while closing old listener: %v", err)
		}
	}

	srv.listener = ln

	return nil
}

// do creates a goroutine, but maintains a record of it to ensure that execution completes
// before the server is shutdown.
func (srv *Server) do(f func()) {
	srv.wg.Add(1)

	go func() {
		defer srv.wg.Done()
		f()
	}()
}

// registerSession adds the given session to the sessions map. It can be called concurrently.
func (srv *Server) registerSession(s *Session) error {
	// Refuse to register a session if Server is closed
	if srv.shutdown.HasFired() {
		return ErrServerClosed
	}

	srv.mu.Lock()
	defer srv.mu.Unlock()

	if srv.sessions == nil {
		srv.sessions = make(map[*Session]struct{})
	}

	srv.sessions[s] = struct{}{}

	return nil
}

// removeSession remove session from sessions map. It can be called concurrently.
func (srv *Server) removeSession(s *Session) {
	srv.mu.Lock()
	defer srv.mu.Unlock()

	delete(srv.sessions, s)
}
