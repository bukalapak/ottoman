// Package jose provides an implementation of the Javascript Object Signing and Encryption.
package jose

import (
	"crypto/rsa"
	"strings"

	jose "gopkg.in/square/go-jose.v2"
)

// Stamper is the interface that handles commonly used JOSE operations.
type Stamper interface {
	Encode(payload []byte) (string, error)
	Encrypt(payload []byte) (string, error)
	Decode(data string) ([]byte, error)
	Decrypt(data string) ([]byte, error)
}

// Standard implements Stamper interface.
type Standard struct {
	signer     jose.Signer
	encrypter  jose.Encrypter
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

// New returns a Standard. It's using predefined key, compression and algorithm.
func New(publicKey, privateKey string) (*Standard, error) {
	prv, err := RSAPrivateKey(readCertificate(privateKey))
	if err != nil {
		return nil, err
	}

	pub, err := RSAPublicKey(readCertificate(publicKey))
	if err != nil {
		return nil, err
	}

	signer, _ := jose.NewSigner(jose.SigningKey{Key: prv, Algorithm: jose.RS256}, nil)
	encrypter, _ := jose.NewEncrypter(jose.A256GCM,
		jose.Recipient{Key: pub, Algorithm: jose.RSA_OAEP_256},
		&jose.EncrypterOptions{Compression: jose.NONE},
	)

	return &Standard{
		signer:     signer,
		encrypter:  encrypter,
		privateKey: prv,
		publicKey:  pub,
	}, nil
}

// Encode signs a payload and produces a signed JWS object.
func (s *Standard) Encode(payload []byte) (string, error) {
	token, _ := s.signer.Sign(payload)
	return token.CompactSerialize()
}

// Decode parses a signed message in compact or full serialization format.
func (s *Standard) Decode(data string) ([]byte, error) {
	token, err := jose.ParseSigned(data)
	if err != nil {
		return nil, err
	}

	return token.Verify(s.publicKey)
}

// Encrypt encrypts a payload and produces an encrypted JWE object.
func (s *Standard) Encrypt(payload []byte) (string, error) {
	token, _ := s.encrypter.Encrypt(payload)
	return token.CompactSerialize()
}

// Decrypt parses an encrypted message in compact or full serialization format.
func (s *Standard) Decrypt(data string) ([]byte, error) {
	token, err := jose.ParseEncrypted(data)
	if err != nil {
		return nil, err
	}

	return token.Decrypt(s.privateKey)
}

func readCertificate(s string) []byte {
	return []byte(strings.TrimSpace(s))
}
