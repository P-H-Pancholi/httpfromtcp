package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	headers := make(Headers)
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	assert.NoError(t, err)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.False(t, done)
	assert.Equal(t, len(data), n+2)

	headers = make(Headers)
	data = []byte("      Host : localhost:42069   \r\n\r\n")
	n, done, err = headers.Parse(data)
	assert.Error(t, err)
	assert.False(t, done)
	assert.Equal(t, 0, n)

	headers = make(Headers)
	data = []byte("      Host: localhost:42069          \r\n\r\n")
	n, done, err = headers.Parse(data)
	assert.NoError(t, err)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.False(t, done)
	assert.Equal(t, len(data), n+2)

	headers = make(Headers)
	data = []byte("\r\n\r\n")
	n, done, err = headers.Parse(data)
	assert.True(t, done)
	assert.Equal(t, 2, n)
	assert.NoError(t, err)

	headers = make(Headers)
	data = []byte("      H©st: localhost:42069          \r\n\r\n")
	n, done, err = headers.Parse(data)
	assert.Error(t, err)
	assert.False(t, done)
	assert.Equal(t, 0, n)

	headers = make(Headers)
	headers["language"] = "golang"
	data = []byte("      Language: python      \r\n\r\n")
	n, done, err = headers.Parse(data)
	assert.NoError(t, err)
	assert.False(t, done)
	assert.Equal(t, n+2, len(data))
	assert.Equal(t, "golang, python", headers["language"])

	// Test: Valid 2 headers with existing headers
	headers = map[string]string{"host": "localhost:42069"}
	data = []byte("User-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, "curl/7.81.0", headers["user-agent"])
	assert.Equal(t, 25, n)
	assert.False(t, done)
}
