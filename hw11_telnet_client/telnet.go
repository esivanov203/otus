package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
	Run()
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

func (c *telnetClient) Run() {
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		err := c.Send()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Отправка завершена: %v\n", err)
		} else {
			_, _ = fmt.Fprintf(os.Stderr, "Отправлено EOF\n")
		}
		_ = c.Close()
	}()

	go func() {
		defer wg.Done()
		err := c.Receive()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Соединение закрыто сервером: %v\n", err)
		} else {
			_, _ = fmt.Fprintf(os.Stderr, "Соединение закрыто сервером\n")
		}
		_ = c.Close()
	}()

	wg.Wait()
}
