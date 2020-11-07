package server

import (
	"errors"
	"fmt"
	"net"
	"sync"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	"github.com/hired/gevulot/pkg/dispatcher"
	"github.com/hired/gevulot/pkg/pg"
)

// Session represents a proxied PostgreSQL database session.
type Session struct {
	// Guards dbConnectionParams
	mu sync.Mutex

	// Global configuration
	cfg ConfigStore

	// Connection from a client to the Gevulot
	clientConn *pg.Conn

	// Connection from the Gevulot to the database
	dbConn *pg.Conn

	// Cached database connection parameters from the config
	dbConnectionParams pg.ConnectionParams

	// ┌──────────┐                  ┌─────────────────┐                  ┌──────────┐
	// │          │◀───── dbOut ─────│                 │◀─── clientIn ────│          │
	// │    DB    │                  │     Gevulot     │                  │  Client  │
	// │          │────── dbIn ─────▶│                 │──── clientOut ──▶│          │
	// └──────────┘                  └─────────────────┘                  └──────────┘

	//
	// Set of channels to send/receive PG messages.
	//

	clientIn  chan pg.Message // client -> Gevulot
	clientOut chan pg.Message // Gevulot -> client
	dbIn      chan pg.Message // DB -> Gevulot
	dbOut     chan pg.Message // Gevulot -> DB

	// Fired when session is closed
	closed *Event
}

var (
	// ErrSessionClosed is returned by the Server's Start after a call to Close.
	ErrSessionClosed = errors.New("session: Session closed")
)

// NewSession initializes a new Session.
func NewSession(client net.Conn, config ConfigStore) *Session {
	return &Session{
		cfg:        config,
		clientConn: pg.NewConn(client),

		clientIn:  make(chan pg.Message, 64),
		clientOut: make(chan pg.Message, 64),
		dbIn:      make(chan pg.Message, 64),
		dbOut:     make(chan pg.Message, 64),

		closed: NewEvent(),
	}
}

// Start starts the session between a client and the database.
// Start blocks, serving the connection until the client or the database hangs up.
// The caller typically invokes Start in a go statement.
func (s *Session) Start() error {
	log.Info("session: initializing a new session")
	defer log.Info("session: closed")

	if s.closed.HasFired() {
		return ErrSessionClosed
	}

	// Automatically close the session
	defer s.Close()

	// Initialize connection params (SSL, encoding, etc.)
	err := s.negotiateSessionParams()

	if err != nil {
		return err
	}

	g := errgroup.Group{}

	// Run session goroutines capturing errors
	g.Go(s.startClientInPump)
	g.Go(s.startClientOutPump)
	g.Go(s.startDBInPump)
	g.Go(s.startDBOutPump)
	g.Go(s.startProcessing)

	// Wait for the first error (or successful exit)
	err = g.Wait()

	if err != nil {
		log.Errorf("session: error: %v", err)
	} else {
		log.Info("session: completed")
	}

	return err
}

// Close immediately closes Session's underlying network connections.
// Close returns any error returned from closing db/client connections.
//
// Once Close has been called on a Session, it may not be reused.
func (s *Session) Close() (err error) {
	// Ensure that we close session only once
	if !s.closed.Fire() {
		return nil
	}

	log.Info("session: closing")

	// Close network connections — this will stop in pumps
	if s.clientConn != nil {
		err = s.clientConn.Close()

		log.Debugf("session: client connection is closed; err = %v", err)
	}

	if s.dbConn != nil {
		err = s.dbConn.Close()

		log.Debugf("session: db connection is closed; err = %v", err)
	}

	// Close all message channels — this will stop out pumps and processing goroutine
	close(s.clientIn)
	close(s.clientOut)
	close(s.dbIn)
	close(s.dbOut)

	log.Debug("session: channels closed")

	return
}

// negotiateSessionParams establish session parameters with the client, then connects
// to the specified in config DB on behalf of the client.
func (s *Session) negotiateSessionParams() error {
	log.Debug("session: waiting for the client startup message")

	// Receiving initial startup message from the client. It contains username, database name etc.
	startupMessage, err := s.clientConn.RecvStartupMessage()

	if err != nil {
		return err
	}

	// Check if startup message is a SSL request
	if startupMessage.ProtocolVersion == pg.SSLRequestMagic {
		log.Info("session: client requested SSL; denying")

		err = s.clientConn.SendByte('N')

		if err != nil {
			return err
		}

		// DANGER! There is a possibility of an infinite loop here.
		return s.negotiateSessionParams()
	}

	// Check the protocol version just in case
	if startupMessage.ProtocolVersion != pg.DefaultProtocolVersion {
		return fmt.Errorf("session: unsupported PG protocol version %v", startupMessage.ProtocolVersion)
	}

	// Check that client is trying to connect to the database that we are proxying
	allowedDB, err := s.getDBConnnectionParam("database")

	if err != nil {
		return err
	}

	if dbName := startupMessage.GetParameter("database"); dbName != allowedDB {
		return fmt.Errorf("session: database mismatch: %v != %v", dbName, allowedDB)
	}

	// Establish DB connection on behalf of the client
	return s.establishDBConnection(startupMessage)
}

// establishDBConnection connects to the database using connection parameters from the config.
func (s *Session) establishDBConnection(startupMessage pg.Message) error {
	// Get database connection params from the config
	host, err := s.getDBConnnectionParam("host")

	if err != nil {
		return err
	}

	port, err := s.getDBConnnectionParam("port")

	if err != nil {
		return err
	}

	// Connect to the database
	conn, err := net.Dial("tcp", net.JoinHostPort(host, port))

	if err != nil {
		return err
	}

	// Convert to pg.NewConn
	s.dbConn = pg.NewConn(conn)

	// Send initial startup message that we received from the client
	err = s.dbConn.SendMessage(startupMessage)

	if err != nil {
		return err
	}

	return nil
}

// startClientInPump pumps messages from the client into the clientIn channel.
func (s *Session) startClientInPump() error {
	for {
		message, err := s.clientConn.RecvMessage()

		if err != nil {
			return err
		}

		s.clientIn <- message
	}
}

// startClientOutPump pumps messages from the clientOut channel to the client.
func (s *Session) startClientOutPump() error {
	for {
		// Get next message in queue
		message, ok := <-s.clientOut

		// Channel is closed — exit
		if !ok {
			return nil
		}

		// Send the message over the network
		err := s.clientConn.SendMessage(message)

		if err != nil {
			return err
		}
	}
}

// startClientOutPump pumps messages from the database into the dbIn channel.
func (s *Session) startDBInPump() error {
	for {
		message, err := s.dbConn.RecvMessage()

		if err != nil {
			return err
		}

		s.dbIn <- message
	}
}

// startDBOutPump pumps messages from the dbOut channel to the database.
func (s *Session) startDBOutPump() error {
	for {
		// Get next message in queue
		message, ok := <-s.dbOut

		// Channel is closed — exit
		if !ok {
			return nil
		}

		// Send the message over the network
		err := s.dbConn.SendMessage(message)

		if err != nil {
			return err
		}
	}
}

// startProcessing dispatches messages between the database and the client.
func (s *Session) startProcessing() error {
	// FIXME: temporary just proxy everything from the client to the database and vice versa
	for {
		select {
		case clientMsg, ok := <-s.clientIn:
			if !ok {
				return nil
			}

			log.Infof("-> %s", clientMsg.Frame().Bytes())

			s.dbOut <- clientMsg

		case dbMsg, ok := <-s.dbIn:
			if !ok {
				return nil
			}

			log.Infof("<- %s", dbMsg.Frame().Bytes())

			s.clientOut <- dbMsg
		}
	}
}

// getDBConnnectionParam returns connection parameter with given name from the config.
func (s *Session) getDBConnnectionParam(name string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// First call to getDBConnnectionParam — parse params from the config
	if s.dbConnectionParams == nil {
		// Get() will block until we have a config
		config, err := s.cfg.Get()

		if err != nil {
			return "", err
		}

		// Parse database URI
		s.dbConnectionParams, err = pg.ParseDatabaseURI(config.DatabaseURL)

		if err != nil {
			return "", err
		}
	}

	return s.dbConnectionParams[name], nil
}
