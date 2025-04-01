package telnet

import (
	"io"
	"net"
	"time"
)

// TelnetClient defines the interface for a Telnet client.
type TelnetClient interface { //nolint
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

// telnetClient implements the TelnetClient interface.
type telnetClient struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    io.ReadWriteCloser
}

// Connect establishes a TCP connection to the specified address with the given timeout.
func (c *telnetClient) Connect() error {
	conn, err := net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

// Close closes the connection if it is not nil.
func (c *telnetClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// Send writes data from the input reader to the connection.
func (c *telnetClient) Send() error {
	if c.conn == nil {
		return io.EOF
	}
	_, err := io.Copy(c.conn, c.in)
	return err
}

// Receive reads data from the connection and writes it to the output writer.
func (c *telnetClient) Receive() error {
	if c.conn == nil {
		return io.EOF
	}
	_, err := io.Copy(c.out, c.conn)
	return err
}

// NewTelnetClient creates a new TelnetClient
// instance with the provided address, timeout, input reader, and output writer.
func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}
