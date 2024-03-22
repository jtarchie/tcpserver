package tcpserver

import (
	"bufio"
	"fmt"
	"net"
)

type Client struct {
	conn   *net.TCPConn
	port   int
	reader *bufio.Reader
}

func NewClient(port int) (*Client, error) {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return nil, fmt.Errorf("could resolve address (localhost:%d): %w", port, err)
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return nil, fmt.Errorf("could not dial (localhost:%d): %w", port, err)
	}

	return &Client{
		conn:   conn,
		port:   port,
		reader: bufio.NewReader(conn),
	}, nil
}

func (c *Client) WriteString(message string) error {
	_, err := c.conn.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("could not write message (localhost:%d): %w", c.port, err)
	}

	return nil
}

func (c *Client) Close() error {
	err := c.conn.Close()
	if err != nil {
		return fmt.Errorf("could not close connection (localhost:%d): %w", c.port, err)
	}

	return nil
}

func (c *Client) ReadlineString() (string, error) {
	contents, _, err := c.reader.ReadLine()
	if err != nil {
		return "", fmt.Errorf("could read from (localhost:%d): %w", c.port, err)
	}

	return string(contents), nil
}

func Write(port int, message string) (string, error) {
	client, err := NewClient(port)
	if err != nil {
		return "", fmt.Errorf("could not create client: %w", err)
	}
	defer client.Close()

	err = client.WriteString(message)
	if err != nil {
		return "", fmt.Errorf("could not write message: %w", err)
	}

	message, err = client.ReadlineString()
	if err != nil {
		return "", fmt.Errorf("could not read line: %w", err)
	}

	return message, nil
}
