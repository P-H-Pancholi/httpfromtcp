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

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
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
	w.Write([]byte(s))
	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.Headers{}
	h["content-length"] = strconv.Itoa(contentLen)
	h["connection"] = "close"
	h["content-type"] = "text/plain"
	return h
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	s := ""

	contentLen, ok := headers.Get("Content-length")
	if !ok {
		return errors.New("content-length does not exists in headers")
	}
	s += fmt.Sprintf("Content-Length: %s\r\n", contentLen)

	con, ok := headers.Get("Connection")
	if !ok {
		return errors.New("connection does not exists in headers")
	}
	s += fmt.Sprintf("Connection: %s\r\n", con)

	conType, ok := headers.Get("Content-Type")
	if !ok {
		return errors.New("content-Type does not exists in headers")
	}
	s += fmt.Sprintf("Content-Type: %s\r\n", conType)

	s += "\r\n"

	w.Write([]byte(s))
	return nil
}
