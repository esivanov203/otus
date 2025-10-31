package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type telnetClient struct {
	address string
	timeout time.Duration
	conn    net.Conn
	in      io.ReadCloser
	out     io.Writer
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (c *telnetClient) Connect() error {
	conn, err := net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return err
	}
	c.conn = conn

	return nil
}

func (c *telnetClient) Send() error {
	if c.conn == nil {
		return fmt.Errorf("connection closed")
	}
	_, err := io.Copy(c.conn, c.in)
	if errors.Is(err, io.EOF) {
		return nil
	}
	return err
}

func (c *telnetClient) Receive() error {
	if c.conn == nil {
		return fmt.Errorf("connection closed")
	}
	_, err := io.Copy(c.out, c.conn)
	if errors.Is(err, io.EOF) {
		return nil
	}
	return err
}

func (c *telnetClient) Close() error {
	if c.conn != nil {
		err := c.conn.Close()
		c.conn = nil // чтобы Send/Receive знали, что соединение закрыто
		return err
	}
	return nil
}
