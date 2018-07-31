package random

import (
	"encoding/hex"
	"io"
)

func Bytes(r io.Reader, n int) ([]byte, error) {
	bytes := make([]byte, n)

	if _, err := r.Read(bytes); err != nil {
		return nil, err
	}

	return bytes, nil
}

func Hex(r io.Reader, n int) (string, error) {
	bytes, err := Bytes(r, n)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}
