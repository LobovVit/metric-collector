package signature

import (
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"fmt"
)

func CreateSignature(data []byte, key string) ([]byte, error) {
	h := hmac.New(sha256.New, []byte(key))
	_, err := h.Write(data)
	if err != nil {
		return nil, fmt.Errorf("sign: %w", err)
	}
	dst := h.Sum(nil)
	return dst, nil
}

func CheckSignature(data []byte, hash string, key string) error {
	h := hmac.New(sha256.New, []byte(key))
	_, err := h.Write(data)
	if err != nil {
		return fmt.Errorf("sign: %w", err)
	}
	sign := h.Sum(nil)
	if !hmac.Equal([]byte(fmt.Sprintf("%x", sign)), []byte(hash)) {
		return errors.New("signature is not correct")
	}
	return nil
}
