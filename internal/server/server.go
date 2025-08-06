package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/P-H-Pancholi/httpfromtcp/internal/request"
	"github.com/P-H-Pancholi/httpfromtcp/internal/response"
)

type State int64

type Server struct {
	Listner     net.Listener
	closed      atomic.Bool
	handlerfunc Handler
}

type HandlerError struct {
	StatusCode int
	Message    string
}

type Handler func(w response.Writer, req *request.Request)

// func writeError(w io.Writer, handlerErr *HandlerError) {
// 	h := response.GetDefaultHeaders(len(handlerErr.Message))
// 	if err := response.WriteStatusLine(w, response.StatusCode(handlerErr.StatusCode)); err != nil {
// 		log.Printf("Error while writing status line: %v", err)
// 		return
// 	}
// 	if err := response.WriteHeaders(w, h); err != nil {
// 		log.Printf("Error while writing headers line: %v", err)
// 		return
// 	}
// 	w.Write([]byte(handlerErr.Message))

// }

func Serve(port int, handlerfunc Handler) (*Server, error) {
	list, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	// go func() {
	// 	list.
	// }()
	serv := Server{
		Listner:     list,
		handlerfunc: handlerfunc,
	}
	go serv.listen()
	return &serv, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	if s.Listner != nil {
		return s.Listner.Close()
	}
	return nil
}

func (s *Server) listen() {
	for {
		conn, err := s.Listner.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			log.Printf("Error accepting connection %v\n", err)
			continue
		}
		go s.handler(conn)
	}
}

func (s *Server) handler(conn net.Conn) {
	defer conn.Close()
	// data := "HTTP/1.1 200 OK\r\n" + "Content-Type: text/plain\r\n" + "\r\n" + "Hello World!\n"
	// conn.Write([]byte(data))

	r, err := request.RequestFromReader(conn)
	if err != nil {
		log.Printf("Error while parsing request from connection: %v", err)
	}
	w := response.NewWriter(conn)
	s.handlerfunc(w, r)

}
