package main

import (
	"strings"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

// PubKeyFile : is the public key file name
const PubKeyFile string = "rsa_public_key.pub"

// PrivKeyFile : is the private key file name
const PrivKeyFile string = "rsa_private_key.priv"

func main() {
	/* if false {
		//File Encryption
		privateKey, publicKey := keys.GenerateKeyPair(4096)
		file.BytesToNewFile(PubKeyFile, keys.PublicKeyToBytes(publicKey))
		file.BytesToNewFile(PrivKeyFile, keys.PrivateKeyToBytes(privateKey))
		//Encrypt file
		var key *[]byte
		key = encryption.GenKey()
		file.BytesToNewFile("safe_key", keys.EncryptWithPublicKey(key, publicKey))
		nbrFiles, err := file.EncryptAll("C:\\Users\\peter\\Downloads", key)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Encrypted files number: ", nbrFiles)
	} else {
		//File Decryption
		encryptedKey, err := file.BytesFromFile("safe_key")
		if err != nil {
			log.Fatal(err)
		}
		privateKeyBytes, err := file.BytesFromFile(PrivKeyFile)
		if err != nil {
			log.Fatal(err)
		}
		privateKey := keys.BytesToPrivateKey(privateKeyBytes)
		var key *[]byte
		key = keys.DecryptWithPrivateKey(encryptedKey, privateKey)
		nbrFiles, err := file.DecryptAll("C:\\Users\\peter\\Downloads", key)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Decrypted files number: ", nbrFiles)
	} */

	var inTE, outTE *walk.TextEdit

	MainWindow{
		Title:   "SCREAMO",
		MinSize: Size{600, 400},
		Layout:  VBox{},
		Children: []Widget{
			HSplitter{
				Children: []Widget{
					TextEdit{AssignTo: &inTE},
					TextEdit{AssignTo: &outTE, ReadOnly: true},
				},
			},
			PushButton{
				Text: "SCREAM",
				OnClicked: func() {
					outTE.SetText(strings.ToUpper(inTE.Text()))
				},
			},
		},
	}.Run()

}
