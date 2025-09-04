package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

const addr = "localhost:42069"

func main() {

	a, err := net.ResolveUDPAddr("udp", addr)

	if err != nil {
		log.Fatal("Could not resolve address ", err)
	}

	log.Println("Found address:", a)

	conn, err := net.DialUDP("udp", nil, a)
	if err != nil {
		log.Fatal("Connection refused to address", a, err)
	}

	defer conn.Close()

	log.Println("Connected to :", conn.LocalAddr())

	r := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")

		l, _, err := r.ReadLine()

		if err != nil {
			fmt.Println("Could not parse message", err)
		}

		n, err := conn.Write(l)
		if err != nil {
			fmt.Println("Could not send message", err)
		}

		fmt.Println("Sending message with number of bytes =", n)

	}

}
