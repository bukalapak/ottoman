package jose_test

import (
	"testing"

	"github.com/bukalapak/ottoman/crypto/jose"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	_, err := jose.New("", rsaPrivateKey)
	assert.NotNil(t, err)

	_, err = jose.New(rsaPublicKey, "")
	assert.NotNil(t, err)
}

func TestEncodeDecode(t *testing.T) {
	b := []byte(`{"foo":"bar"}`)
	n, err := jose.New(rsaPublicKey, rsaPrivateKey)
	assert.Nil(t, err)

	token, err := n.Encode(b)
	assert.Nil(t, err)
	assert.NotNil(t, token)

	data, err := n.Decode(token)
	assert.Nil(t, err)
	assert.Equal(t, b, data)

	out, err := n.Decode("x")
	assert.NotNil(t, err)
	assert.Nil(t, out)
}

func TestEncryptDecrypt(t *testing.T) {
	b := []byte(`{"foo":"bar"}`)
	n, err := jose.New(rsaPublicKey, rsaPrivateKey)
	assert.Nil(t, err)

	token, err := n.Encrypt(b)
	assert.Nil(t, err)
	assert.NotNil(t, token)

	data, err := n.Decrypt(token)
	assert.Nil(t, err)
	assert.Equal(t, b, data)

	out, err := n.Decrypt("x")
	assert.NotNil(t, err)
	assert.Nil(t, out)
}
