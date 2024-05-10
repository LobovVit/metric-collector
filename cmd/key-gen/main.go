package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"os"
)

func main() {
	reader := rand.Reader
	bitSize := 2048
	key, err := rsa.GenerateKey(reader, bitSize)
	if err != nil {
		panic(err)
	}
	publicKey := key.PublicKey
	savePEMKey("private.pem", key)
	savePublicPEMKey("public.pem", publicKey)
}

func savePEMKey(fileName string, key *rsa.PrivateKey) {
	outFile, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	var privateKey = &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	err = pem.Encode(outFile, privateKey)
	if err != nil {
		panic(err)
	}
}

func savePublicPEMKey(fileName string, pubKey rsa.PublicKey) {
	asn1Bytes, err := asn1.Marshal(pubKey)
	if err != nil {
		panic(err)
	}

	var pemKey = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	pemFile, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer pemFile.Close()

	err = pem.Encode(pemFile, pemKey)
	if err != nil {
		panic(err)
	}
}
