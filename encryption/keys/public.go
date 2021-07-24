package keys

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"log"
)

// PublicKeyToBytes : public key to bytes
func PublicKeyToBytes(pub *rsa.PublicKey) []byte {
	pubASN1 := x509.MarshalPKCS1PublicKey(pub)
	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubASN1,
	})

	return pubBytes
}

// BytesToPublicKey : bytes to public key
func BytesToPublicKey(pub []byte) *rsa.PublicKey {
	block, _ := pem.Decode(pub)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			log.Fatal(err)
		}
	}
	key, err := x509.ParsePKCS1PublicKey(b)
	if err != nil {
		log.Fatal(err)
	}
	return key
}

// EncryptWithPublicKey : encrypts data with public key
func EncryptWithPublicKey(msg *[]byte, pub *rsa.PublicKey) []byte {
	hash := sha256.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, pub, *msg, nil)
	if err != nil {
		log.Fatal(err)
	}
	return ciphertext
}
