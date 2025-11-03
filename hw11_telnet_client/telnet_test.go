package main

import (
	"bytes"
	"io"
	"net"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTelnetClient(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			in := &bytes.Buffer{}
			out := &bytes.Buffer{}

			timeout, err := time.ParseDuration("10s")
			require.NoError(t, err)

			client := NewTelnetClient(l.Addr().String(), timeout, io.NopCloser(in), out)
			require.NoError(t, client.Connect())
			defer func() { require.NoError(t, client.Close()) }()

			in.WriteString("hello\n")
			err = client.Send()
			require.NoError(t, err)

			err = client.Receive()
			require.NoError(t, err)
			require.Equal(t, "world\n", out.String())
		}()

		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			defer func() { require.NoError(t, conn.Close()) }()

			request := make([]byte, 1024)
			n, err := conn.Read(request)
			require.NoError(t, err)
			require.Equal(t, "hello\n", string(request)[:n])

			n, err = conn.Write([]byte("world\n"))
			require.NoError(t, err)
			require.NotEqual(t, 0, n)
		}()

		wg.Wait()
	})

	t.Run("timeout", func(t *testing.T) {
		timeOut := 500 * time.Millisecond

		input := io.NopCloser(strings.NewReader("hello\n"))
		output := &bytes.Buffer{}

		tc := NewTelnetClient("192.0.2.1:1234", timeOut, input, output)
		err := tc.Connect()

		var netErr net.Error
		require.ErrorAs(t, err, &netErr)
		require.True(t, netErr.Timeout())
	})

	t.Run("closed connection", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer l.Close()

		client := NewTelnetClient(l.Addr().String(), 1*time.Second, io.NopCloser(strings.NewReader("")), &bytes.Buffer{})
		require.NoError(t, client.Connect())
		require.NoError(t, client.Close())

		err = client.Send()
		require.Error(t, err)

		err = client.Receive()
		require.Error(t, err)
	})

	t.Run("send-eof", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer l.Close()

		client := NewTelnetClient(l.Addr().String(), 1*time.Second, io.NopCloser(strings.NewReader("")), &bytes.Buffer{})
		require.NoError(t, client.Connect())
		defer client.Close()

		err = client.Send()
		require.NoError(t, err)
	})

	// Если сокет закрылся со стороны сервера,
	// то при следующей попытке отправить сообщение программа
	// должна завершаться
	t.Run("server closed connection", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer l.Close()

		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			time.Sleep(50 * time.Millisecond)
			conn.Close()
		}()

		go func() {
			defer wg.Done()

			clientInput := strings.Repeat("data line\n", 1024*1024)
			client := NewTelnetClient(
				l.Addr().String(),
				time.Second,
				io.NopCloser(strings.NewReader(clientInput)),
				&bytes.Buffer{},
			)
			require.NoError(t, client.Connect())
			defer client.Close()

			err = client.Send()

			// сбой в процессе передачи большого объема данных
			// из-за того что сокет "неожиданно" закрылся со стороны сервера
			require.Error(t, err)
			var errno syscall.Errno
			require.ErrorAs(t, err, &errno)
			require.True(t, errno == syscall.EPIPE || errno == syscall.ECONNRESET,
				"ожидался EPIPE или ECONNRESET, но получили %v", errno)
		}()

		wg.Wait()
	})
}
