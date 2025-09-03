package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("./messages.txt")
	if err != nil {
		fmt.Println(err)
	}

	lines := getLinesChannel(file)

	for line := range lines {

		fmt.Printf("read: %s\n", line)
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
