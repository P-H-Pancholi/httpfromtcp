package headers

import (
	"bytes"
	"errors"
	"strings"
	"unicode"
)

type Headers map[string]string

const clrf = "\r\n"

func (h Headers) Parse(data []byte) (int, bool, error) {

	idx := bytes.Index(data, []byte(clrf))
	if idx == -1 {
		return 0, false, nil
	}
	if idx == 0 {
		return 2, true, nil
	}

	currHeaders := string(data[:idx])
	currHeaders = strings.TrimSpace(currHeaders)

	colonIdx := strings.Index(currHeaders, ":")
	if colonIdx == -1 {
		return 0, false, errors.New("malformed header")
	}
	key := currHeaders[:colonIdx]
	value := currHeaders[colonIdx+1:]

	value = strings.TrimSpace(value)
	if unicode.IsSpace(rune(key[len(key)-1])) {
		return 0, false, errors.New("key has trailing whitespace")
	}

	key = strings.ToLower(key)
	for _, c := range key {
		if !(c >= 'a' && c <= 'z') &&
			!(c >= '0' && c <= '9') &&
			c != '&' && c != '#' && c != '!' && c != '$' && c != '%' && c != '\'' &&
			c != '*' && c != '+' && c != '-' && c != '.' && c != '^' && c != '_' && c != '`' &&
			c != '|' && c != '~' {
			// Invalid character in header key
			return 0, false, errors.New("invalid character in header key")
		}
	}
	_, exists := h[key]
	if exists {
		h[key] = h[key] + ", " + value
	} else {
		h[key] = value
	}
	return idx + 2, false, nil
}

func (h Headers) Get(key string) (string, bool) {
	s, exists := h[strings.ToLower(key)]
	return s, exists
}
