package main

import (
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

	defer file.Close()

	const chunkSize = 8

	buff := make([]byte, chunkSize)

	currentLine := ""

	for {
		bytesRead, err := file.Read(buff)

		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
			break
		}
		str := string(buff[:bytesRead])
		splitStrings := strings.Split(str, "\n")
		currentLine += splitStrings[0]
		if len(splitStrings) > 1 {
			fmt.Printf("read: %s\n", currentLine)
			currentLine = ""
			currentLine += splitStrings[1]
		}
	}
}
