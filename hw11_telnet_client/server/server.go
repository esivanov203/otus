package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	// #nosec G102
	l, err := net.Listen("tcp", ":4242")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Listening on port 4242")
	defer func() { _ = l.Close() }()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go func(c net.Conn) {
			defer func() { _ = c.Close() }()
			_, err := io.WriteString(c, "Welcome to server\n")
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Printf("New connection from %s\n", c.RemoteAddr().String())

			scanner := bufio.NewScanner(c)
			for scanner.Scan() {
				line := scanner.Text()
				fmt.Printf("Received from %s: %s\n", c.RemoteAddr().String(), line)
				if line == "quit" || line == "q" {
					_, _ = io.WriteString(c, "Bye\n")
					return
				}
				if line != "" {
					_, _ = io.WriteString(c, fmt.Sprintf("Processed by server: %s\n", line))
				}
			}
		}(conn)
	}
}
