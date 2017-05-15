package jose

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

// RSAPrivateKey parses data as *rsa.PrivateKey
func RSAPrivateKey(data []byte) (*rsa.PrivateKey, error) {
	input := pemDecode(data)

	var err error
	var key interface{}

	if key, err = x509.ParsePKCS1PrivateKey(input); err != nil {
		if key, err = x509.ParsePKCS8PrivateKey(input); err != nil {
			return nil, err
		}
	}

	return key.(*rsa.PrivateKey), nil
}

// RSAPublicKey parses data as *rsa.PublicKey
func RSAPublicKey(data []byte) (*rsa.PublicKey, error) {
	input := pemDecode(data)

	var err error
	var key interface{}

	if key, err = x509.ParsePKIXPublicKey(input); err != nil {
		if cert, err := x509.ParseCertificate(input); err == nil {
			key = cert.PublicKey
		} else {
			return nil, err
		}
	}

	return key.(*rsa.PublicKey), nil
}

func pemDecode(data []byte) []byte {
	if block, _ := pem.Decode(data); block != nil {
		return block.Bytes
	}

	return data
}
