// Package jose provides an implementation of the Javascript Object Signing and Encryption.
package jose

import (
	"crypto/rsa"
	"strings"

	"gopkg.in/square/go-jose.v2"
)

// Signature is the interface that handles commonly used JWS operations.
type Signature interface {
	Encode(payload []byte) (string, error)
	Decode(data string) ([]byte, error)
}

// Encryption is the interface that handles commonly used JWE operations.
type Encryption interface {
	Encrypt(payload []byte) (string, error)
	Decrypt(data string) ([]byte, error)
}

// Stamper is the interface that handles commonly used JOSE operations.
type Stamper interface {
	Signature
	Encryption
}

type signatureStandard struct {
	signer    jose.Signer
	publicKey *rsa.PublicKey
}

type encryptionStandard struct {
	encrypter  jose.Encrypter
	privateKey *rsa.PrivateKey
}

type standard struct {
	signature  Signature
	encryption Encryption
}

// NewSignature returns a Signature. It's using predefined key and algorithm.
func NewSignature(publicKey, privateKey string) (Signature, error) {
	pub, prv, err := readKeys(publicKey, privateKey)
	if err != nil {
		return nil, err
	}

	signer, _ := jose.NewSigner(jose.SigningKey{Key: prv, Algorithm: jose.RS256}, nil)

	return &signatureStandard{
		signer:    signer,
		publicKey: pub,
	}, nil
}

// Encode signs a payload and produces a signed JWS object.
func (s *signatureStandard) Encode(payload []byte) (string, error) {
	token, _ := s.signer.Sign(payload)
	return token.CompactSerialize()
}

// Decode parses a signed message in compact or full serialization format.
func (s *signatureStandard) Decode(data string) ([]byte, error) {
	return decode(s.publicKey, data)
}

// New returns a Encryption. It's using predefined key, compression and algorithm.
func NewEncryption(publicKey, privateKey string) (Encryption, error) {
	pub, prv, err := readKeys(publicKey, privateKey)
	if err != nil {
		return nil, err
	}

	encrypter, _ := jose.NewEncrypter(jose.A256GCM,
		jose.Recipient{Key: pub, Algorithm: jose.RSA_OAEP_256},
		&jose.EncrypterOptions{Compression: jose.NONE},
	)

	return &encryptionStandard{
		encrypter:  encrypter,
		privateKey: prv,
	}, nil
}

// Encrypt encrypts a payload and produces an encrypted JWE object.
func (s *encryptionStandard) Encrypt(payload []byte) (string, error) {
	token, _ := s.encrypter.Encrypt(payload)
	return token.CompactSerialize()
}

// Decrypt parses an encrypted message in compact or full serialization format.
func (s *encryptionStandard) Decrypt(data string) ([]byte, error) {
	return decrypt(s.privateKey, data)
}

// New returns a Stamper. It's using predefined key, compression and algorithm.
func New(publicKey, privateKey string) (Stamper, error) {
	_, _, err := readKeys(publicKey, privateKey)
	if err != nil {
		return nil, err
	}

	signature, _ := NewSignature(publicKey, privateKey)
	encryption, _ := NewEncryption(publicKey, privateKey)

	return &standard{
		signature:  signature,
		encryption: encryption,
	}, nil
}

// Encode signs a payload and produces a signed JWS object.
func (s *standard) Encode(payload []byte) (string, error) {
	return s.signature.Encode(payload)
}

// Decode parses a signed message in compact or full serialization format.
func (s *standard) Decode(data string) ([]byte, error) {
	return s.signature.Decode(data)
}

// Encrypt encrypts a payload and produces an encrypted JWE object.
func (s *standard) Encrypt(payload []byte) (string, error) {
	return s.encryption.Encrypt(payload)
}

// Decrypt parses an encrypted message in compact or full serialization format.
func (s *standard) Decrypt(data string) ([]byte, error) {
	return s.encryption.Decrypt(data)
}

// Decode parses a signed message in compact or full serialization format.
func Decode(publicKey string, data string) ([]byte, error) {
	pub, err := RSAPublicKey(readCertificate(publicKey))
	if err != nil {
		return nil, err
	}

	return decode(pub, data)
}

// Decrypt parses an encrypted message in compact or full serialization format.
func Decrypt(privateKey string, data string) ([]byte, error) {
	prv, err := RSAPrivateKey(readCertificate(privateKey))
	if err != nil {
		return nil, err
	}

	return decrypt(prv, data)
}

func decode(pub *rsa.PublicKey, data string) ([]byte, error) {
	token, err := jose.ParseSigned(data)
	if err != nil {
		return nil, err
	}

	return token.Verify(pub)
}

func decrypt(prv *rsa.PrivateKey, data string) ([]byte, error) {
	token, err := jose.ParseEncrypted(data)
	if err != nil {
		return nil, err
	}

	return token.Decrypt(prv)
}

func readCertificate(s string) []byte {
	return []byte(strings.TrimSpace(s))
}

func readKeys(publicKey, privateKey string) (pub *rsa.PublicKey, prv *rsa.PrivateKey, err error) {
	pub, err = RSAPublicKey(readCertificate(publicKey))
	if err != nil {
		return
	}

	prv, err = RSAPrivateKey(readCertificate(privateKey))
	if err != nil {
		return
	}

	return
}
