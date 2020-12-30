// Generate Public and Private key for Encryption and Decryption function.

package main

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"log"
	"os"
)

func main() {
	keys, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal(err)
	}
	keyPubFile, err := os.Create("rsa_key.pub")
	if err != nil {
		log.Fatal(err)
	}

	_, err = keyPubFile.WriteString(fmt.Sprintf("%d\n", keys.PublicKey.N))
	_, err = keyPubFile.WriteString(fmt.Sprintf("%d\n", keys.PublicKey.E))

	keyPubFile.Close()

	keyPrivFile, err := os.Create("rsa_key.priv")
	if err != nil {
		log.Fatal(err)
	}
	_, err = keyPrivFile.WriteString(fmt.Sprintf("%d\n", keys.D))
	_, err = keyPrivFile.WriteString(fmt.Sprintf("%d\n", keys.Primes))
	keyPrivFile.Close()
}
