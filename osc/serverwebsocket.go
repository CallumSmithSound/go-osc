package osc

import (
	"bufio"
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/coder/websocket"
)

// Server represents an OSC server. The server listens on Address and Port for
// incoming OSC packets and bundles.
type WebServer struct {
	Addr        string
	Pattern     string
	Dispatcher  Dispatcher
	ReadTimeout time.Duration
	close       func() error
}

// ListenAndServe retrieves incoming OSC packets and dispatches the retrieved
// OSC packets.
func (s *WebServer) ListenAndServe() error {
	http.HandleFunc(s.Pattern, s.Handler)

	defer s.CloseConnection()

	if s.Dispatcher == nil {
		s.Dispatcher = NewStandardDispatcher()
	}

	return http.ListenAndServe(s.Addr, nil)
}

// Serve retrieves incoming OSC packets from the given connection and dispatches
// retrieved OSC packets. If something goes wrong an error is returned.
func (s *WebServer) Handler(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		// ...
		return
	}
	defer c.CloseNow()

	// Set the context as needed. Use of r.Context() is not recommended
	// to avoid surprising behavior (see http.Hijacker).
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, reader, err := c.Reader(ctx)
	if err != nil {
		// ...
		return
	}

	msg, err := readPacket(bufio.NewReader(reader))
	if err != nil {
		// ...
		return
	}

	s.Dispatcher.Dispatch(msg, r.RemoteAddr)
}

// CloseConnection forcibly closes a server's connection.
//
// This causes a "use of closed network connection" error the next time the
// server attempts to read from the connection.
func (s *WebServer) CloseConnection() error {
	if s.close == nil {
		return nil
	}

	err := s.close()
	// If we get "use of closed network connection", it's not a problem because
	// closing the network connection is exactly what we wanted to do!
	if err != nil && !strings.Contains(
		err.Error(), "use of closed network connection",
	) {
		return err
	}

	return nil
}
