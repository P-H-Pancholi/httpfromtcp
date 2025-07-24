package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	udpAdd, err := net.ResolveUDPAddr("udp", "localhost:42069")

	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialUDP("udp", nil, udpAdd)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println(">")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}

		n, err := conn.Write([]byte(input))
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("%d bytes send to connection", n)
	}
}
