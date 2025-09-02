package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	file, err := os.Open("./messages.txt")
	if err != nil {
		fmt.Println(err)
	}

	defer file.Close()

	const chunkSize = 8

	buff := make([]byte, chunkSize)

	for {
		bytesRead, err := file.Read(buff)

		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
			break
		}
		fmt.Printf("read: %s\n", string(buff[:bytesRead]))
	}
}
