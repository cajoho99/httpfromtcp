package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	"github.com/cajoho99/httpfromtcp/internal/request"
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
		}

		fmt.Printf("Connection established with %s\n", conn.RemoteAddr())

		req, err := request.RequestFromReader(conn)

		req.RequestLine.Print()

		fmt.Println("Connection to ", conn.RemoteAddr(), "closed")

		conn.Close()
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	out := make(chan string)

	go func() {
		defer f.Close()
		defer close(out)

		currentLine := ""

		for {
			buff := make([]byte, 8, 8)
			n, err := f.Read(buff)
			if err != nil {
				if currentLine != "" {
					out <- currentLine
				}

				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Printf("error: %s\n", err.Error())
				return
			}
			str := string(buff[:n])
			splitStrings := strings.Split(str, "\n")
			for i := 0; i < len(splitStrings)-1; i++ {
				out <- fmt.Sprintf("%s%s", currentLine, splitStrings[i])
				currentLine = ""

			}
			currentLine += splitStrings[len(splitStrings)-1]
		}
	}()

	return out

}
