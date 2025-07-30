package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/P-H-Pancholi/httpfromtcp/internal/headers"
)

type State int64

const bufferSize = 8

const (
	Initialized State = iota
	parseHeaders
	Done
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	State       State
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize, bufferSize)

	readToIndex := 0
	r := &Request{
		State:   Initialized,
		Headers: make(headers.Headers),
	}

	for r.State != Done {
		if readToIndex >= cap(buf) {
			newbuf := make([]byte, len(buf)*2)
			copy(newbuf, buf)
			buf = newbuf
		}

		// read into the buffer
		n, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				r.State = Done
				break
			}
			return nil, err
		}
		readToIndex += n

		// parse from the buffer
		bytesRead, err := r.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[bytesRead:])
		readToIndex -= bytesRead
	}
	return r, nil
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	idx := bytes.Index(data, []byte("\r\n"))
	if idx == -1 {
		return &RequestLine{}, 0, nil
	}

	reqLine := string(data[:idx])
	splitedReqLine := strings.Split(reqLine, " ")

	if len(splitedReqLine) != 3 {
		return &RequestLine{}, 0, fmt.Errorf("request line does not have all sections")
	}

	for _, s := range splitedReqLine[0] {
		if unicode.IsLetter(s) && !unicode.IsUpper(s) {
			return &RequestLine{}, 0, fmt.Errorf("request method is invalid")
		}
	}
	var r RequestLine

	r.Method = splitedReqLine[0]

	r.RequestTarget = splitedReqLine[1]

	httpVersion := strings.Split(splitedReqLine[2], "/")

	if httpVersion[1] != "1.1" && httpVersion[0] != "HTTP" {
		return &RequestLine{}, 0, fmt.Errorf("http version is invalid")
	}
	r.HttpVersion = httpVersion[1]

	return &r, idx + 2, nil
}

func (r *Request) parse(data []byte) (int, error) {
	numOfBytesParsed := 0
	for r.State != Done {
		n, err := r.parseSingle(data[numOfBytesParsed:])
		if err != nil {
			return 0, err
		}
		numOfBytesParsed += n
		if n == 0 {
			break
		}

	}
	return numOfBytesParsed, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.State {
	case Initialized:
		reqLine, n, err := parseRequestLine(data)

		if err != nil {
			return 0, err
		}

		if n == 0 {
			return 0, nil
		}

		r.RequestLine = *reqLine
		r.State = parseHeaders
		return n, nil
	case parseHeaders:
		n, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if done {
			r.State = Done
		}
		return n, nil
	case Done:
		return 0, fmt.Errorf("error: trying to read data in done state")
	default:
		return 0, fmt.Errorf("error: unknown state")
	}
}
