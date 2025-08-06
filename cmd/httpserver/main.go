package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/P-H-Pancholi/httpfromtcp/internal/headers"
	"github.com/P-H-Pancholi/httpfromtcp/internal/request"
	"github.com/P-H-Pancholi/httpfromtcp/internal/response"
	"github.com/P-H-Pancholi/httpfromtcp/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func handler(w response.Writer, req *request.Request) {
	h := headers.Headers{}
	h["connection"] = "close"
	h["content-type"] = "text/html"
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		w.WriteStatusLine(response.StatusCode(400))
		respStr := `
		<html>
  			<head>
    			<title>400 Bad Request</title>
  			</head>
  			<body>
    			<h1>Bad Request</h1>
    			<p>Your request honestly kinda sucked.</p>
  			</body>
		</html>
		`
		contLen := len([]byte(respStr))
		h["content-length"] = strconv.Itoa(contLen)
		w.WriteHeaders(h)
		w.WriteBody([]byte(respStr))
	case "/myproblem":
		w.WriteStatusLine(response.StatusCode(500))
		respStr := `
		<html>
  			<head>
    			<title>500 Internal Server Error</title>
  			</head>
  			<body>
    			<h1>Internal Server Error</h1>
    			<p>Okay, you know what? This one is on me.</p>
  			</body>
		</html>
		`
		contLen := len([]byte(respStr))
		h["content-length"] = strconv.Itoa(contLen)
		w.WriteHeaders(h)
		w.WriteBody([]byte(respStr))
	default:
		w.WriteStatusLine(response.StatusCode(200))
		respStr := `
		<html>
  			<head>
    			<title>200 OK</title>
  			</head>
  			<body>
    			<h1>Success!</h1>
    			<p>Your request was an absolute banger.</p>
  			</body>
		</html>
		`
		contLen := len([]byte(respStr))
		h["content-length"] = strconv.Itoa(contLen)
		w.WriteHeaders(h)
		w.WriteBody([]byte(respStr))
	}

}
