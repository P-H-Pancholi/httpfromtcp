package response

import (
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/P-H-Pancholi/httpfromtcp/internal/headers"
)

type StatusCode int64

const (
	httpOk                  StatusCode = 200
	httpBadReq              StatusCode = 400
	httpInternalServerError StatusCode = 500
)

type writerState int

const (
	initState writerState = iota
	statusLineState
	headerState
	bodyState
)

type Writer struct {
	data  *io.Writer
	state writerState
}

func NewWriter(wr io.Writer) Writer {
	w := Writer{
		data:  &wr,
		state: initState,
	}
	return w
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.state != initState {
		return errors.New("improper sequence")
	}
	s := "HTTP/1.1" + " " + strconv.Itoa(int(statusCode)) + " "
	switch statusCode {
	case httpOk:
		s += "OK"
	case httpBadReq:
		s += "Bad Request"
	case httpInternalServerError:
		s += "Internal Server Error"
	}
	s += "\r\n"
	writer := *w.data
	writer.Write([]byte(s))
	w.state = statusLineState
	return nil
}

func (w *Writer) WriteHeaders(h headers.Headers) error {
	if w.state != statusLineState {
		return errors.New("improper sequence")
	}

	// Write ALL headers, not just specific ones
	for key, value := range h {
		s := fmt.Sprintf("%s: %s\r\n", key, value)
		writer := *w.data
		writer.Write([]byte(s))
	}

	// End headers section
	writer := *w.data
	writer.Write([]byte("\r\n"))
	w.state = headerState
	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.state != headerState {
		return 0, errors.New("improper sequence")
	}
	writer := *w.data
	writer.Write(p)
	return len(p), nil
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	if w.state != headerState {
		return 0, errors.New("improper sequence")
	}
	wr := *w.data
	s := fmt.Sprintf("%x\r\n", len(p))
	wr.Write([]byte(s))
	wr.Write(p)
	wr.Write([]byte("\r\n"))
	return len(p), nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	if w.state != headerState {
		return 0, errors.New("improper sequence")
	}
	wr := *w.data
	wr.Write([]byte("0\r\n"))
	return 0, nil
}

func (w *Writer) WriteTrailers(h headers.Headers) error {

	s := ""
	xSHA, _ := h.Get("X-Content-Sha256")
	s += fmt.Sprintf("X-Content-Sha256: %s\r\n", xSHA)
	xContLen, _ := h.Get("X-Content-Length")
	s += fmt.Sprintf("X-Content-Length: %s\r\n", xContLen)

	s += "\r\n"
	wr := *w.data
	wr.Write([]byte(s))
	return nil
}

// func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
// 	s := "HTTP/1.1" + " " + strconv.Itoa(int(statusCode)) + " "
// 	switch statusCode {
// 	case httpOk:
// 		s += "OK"
// 	case httpBadReq:
// 		s += "Bad Request"
// 	case httpInternalServerError:
// 		s += "Internal Server Error"
// 	}
// 	s += "\r\n"
// 	w.Write([]byte(s))
// 	return nil
// }

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.Headers{}
	h["content-length"] = strconv.Itoa(contentLen)
	h["connection"] = "close"
	h["content-type"] = "text/plain"
	return h
}

// func WriteHeaders(w io.Writer, headers headers.Headers) error {
// 	s := ""

// 	contentLen, ok := headers.Get("Content-length")
// 	if !ok {
// 		return errors.New("content-length does not exists in headers")
// 	}
// 	s += fmt.Sprintf("Content-Length: %s\r\n", contentLen)

// 	con, ok := headers.Get("Connection")
// 	if !ok {
// 		return errors.New("connection does not exists in headers")
// 	}
// 	s += fmt.Sprintf("Connection: %s\r\n", con)

// 	conType, ok := headers.Get("Content-Type")
// 	if !ok {
// 		return errors.New("content-Type does not exists in headers")
// 	}
// 	s += fmt.Sprintf("Content-Type: %s\r\n", conType)

// 	s += "\r\n"

// 	w.Write([]byte(s))
// 	return nil
// }
