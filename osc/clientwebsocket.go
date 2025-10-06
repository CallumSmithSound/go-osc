package osc

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"time"

	"github.com/coder/websocket"
)

// Client enables you to send OSC packets. It sends OSC messages and bundles to
// the given IP address and port.
type WebClient struct {
	host  string
	path  string
	laddr *net.UDPAddr
	conn  *websocket.Conn
}

// NewClient creates a new OSC client. The Client is used to send OSC
// messages and OSC bundles over an UDP network connection. The `ip` argument
// specifies the IP address and `port` defines the target port where the
// messages and bundles will be send to.
func NewWebClient(host string, path string) *WebClient {
	return &WebClient{host: host, path: path, laddr: nil}
}

func (c *WebClient) Connect() error {
	u := url.URL{Scheme: "ws", Host: c.host, Path: c.path}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	var err error
	c.conn, _, err = websocket.Dial(ctx, u.String(), nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *WebClient) Close() error {
	return c.conn.Close(websocket.StatusNormalClosure, "")
}

// Send sends an OSC Bundle or an OSC Message.
func (c *WebClient) Send(packet Packet) error {
	data, err := packet.MarshalBinary()
	if err != nil {
		return err
	}

	if c.conn == nil {
		return fmt.Errorf("no osc connection has been made")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	if err = c.conn.Write(ctx, websocket.MessageBinary, data); err != nil {
		return err
	}
	return nil
}
