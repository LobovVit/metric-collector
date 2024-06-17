package main

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
)

func Test_savePEMKey(t *testing.T) {

	tests := []struct {
		name     string
		fileName string
	}{
		{name: "test", fileName: "./test.test"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := rand.Reader
			bitSize := 2048
			key, err := rsa.GenerateKey(reader, bitSize)
			if err != nil {
				panic(err)
			}
			savePEMKey(tt.fileName, key)
		})
	}
}

func Test_savePublicPEMKey(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
	}{
		{name: "test", fileName: "./test2.test"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := rand.Reader
			bitSize := 2048
			key, err := rsa.GenerateKey(reader, bitSize)
			if err != nil {
				panic(err)
			}
			publicKey := key.PublicKey
			savePublicPEMKey(tt.fileName, publicKey)
		})
	}
}
