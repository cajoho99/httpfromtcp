package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

const port = ":42069"

func main() {

	tcpListener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}

	defer tcpListener.Close()

	fmt.Println("Listening for TCP traffic on", port)

	for {

		conn, err := tcpListener.Accept()
		if err != nil {
			log.Fatal(err)
			break
		}

		fmt.Printf("Connection established with %s\n", conn.RemoteAddr())

		for line := range getLinesChannel(conn) {
			fmt.Printf("%s\n", line)
		}
		fmt.Println("Connection to ", conn.RemoteAddr(), "closed")
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	out := make(chan string)

	go func() {
		defer f.Close()
		defer close(out)

		currentLine := ""

		for {
			buff := make([]byte, 8)
			n, err := f.Read(buff)
			if err != nil {

				if errors.Is(err, io.EOF) {
					return
				}
				break
			}
			str := string(buff[:n])
			splitStrings := strings.Split(str, "\n")
			currentLine += splitStrings[0]
			if len(splitStrings) > 1 {
				out <- currentLine
				currentLine = ""
				currentLine += splitStrings[1]
			}
		}
	}()

	return out

}
