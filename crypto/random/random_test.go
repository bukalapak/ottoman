package random_test

import (
	"crypto/rand"
	"errors"
	"testing"

	"github.com/bukalapak/ottoman/crypto/random"
	"github.com/stretchr/testify/assert"
)

func TestBytes(t *testing.T) {
	b, err := random.Bytes(rand.Reader, 10)
	assert.Nil(t, err)
	assert.Len(t, b, 10)
}

func TestBytes_error(t *testing.T) {
	b, err := random.Bytes(&errorReader{}, 10)
	assert.NotNil(t, err)
	assert.Empty(t, b)
}

func TestHex(t *testing.T) {
	s, err := random.Hex(rand.Reader, 16)
	assert.Nil(t, err)
	assert.Len(t, s, 32)
}

func TestHex_error(t *testing.T) {
	s, err := random.Hex(&errorReader{}, 10)
	assert.NotNil(t, err)
	assert.Empty(t, s)
}

type errorReader struct{}

func (r errorReader) Read(b []byte) (n int, err error) {
	return 0, errors.New("invalid reader")
}
