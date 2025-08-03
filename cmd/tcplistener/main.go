package main

import (
	"fmt"
	"log"
	"net"

	"github.com/P-H-Pancholi/httpfromtcp/internal/request"
)

func main() {
	list, err := net.Listen("tcp", ":42069")

	if err != nil {
		log.Fatal(err)
	}

	defer list.Close()

	conn, err := list.Accept()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("connection has been accepted")
	// outChan := getLinesChannel(conn)
	// for {

	// 	s, ok := <-outChan

	// 	if !ok {
	// 		fmt.Println("connection has been closed")
	// 		return
	// 	}

	// 	fmt.Println(s)
	// }
	r, err := request.RequestFromReader(conn)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Request line:\n")
	fmt.Printf("- Method: %s\n- Target: %s\n- Version: %s\n", r.RequestLine.Method, r.RequestLine.RequestTarget, r.RequestLine.HttpVersion)
	fmt.Println("Headers:")
	for key, value := range r.Headers {
		fmt.Printf("- %s: %s\n", key, value)
	}
	fmt.Println("Body:")
	fmt.Printf("%s\n", string(r.Body))
}

// func getLinesChannel(f io.ReadCloser) <-chan string {
// 	linesChan := make(chan string)
// 	go func() {
// 		chunk := make([]byte, 8)
// 		currLine := ""

// 		for {
// 			n, err := f.Read(chunk)
// 			if err == io.EOF {
// 				linesChan <- currLine
// 				close(linesChan)
// 				return
// 			}
// 			parts := strings.Split(string(chunk[:n]), "\n")
// 			if len(parts) > 1 {
// 				currLine += parts[0]
// 				linesChan <- currLine
// 				currLine = parts[1]
// 			} else {
// 				currLine += parts[0]
// 			}
// 		}
// 	}()
// 	return linesChan
// }
