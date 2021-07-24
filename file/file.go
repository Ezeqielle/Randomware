package file

import (
	"os"
	"strings"
)

// MaxThreads : is the maximum threads for files encryption/decryption
const MaxThreads int = 10

//MaxSize : maximun allowed Bytes size to be loaded in memory
const MaxSize int = 26214400

// EncryptedExt : extension for all encrypted files
const EncryptedExt string = ".ec"

// HostName : Computer hostname
var HostName string

// SafeFiles : files not to encrypt
var SafeFiles []string

// BytesToNewFile : creates file from fileName and writes data
func BytesToNewFile(fileName string, data []byte) (*os.File, error) {
	//Creates File
	newFile, err := os.Create(fileName)
	if err != nil {
		return newFile, err
	}

	defer newFile.Close()
	//Write to file
	newFile.Write(data)
	newFile.Sync()

	return newFile, nil
}

// BytesFromFile : reads all bytes from file
func BytesFromFile(fileName string) ([]byte, error) {
	//Opens file
	openedFile, err := os.Open(fileName)
	if err != nil {
		return make([]byte, 0), err
	}

	fi, err := openedFile.Stat()

	if err != nil {
		return make([]byte, 0), err
	}

	//Reads from file
	bytesData := make([]byte, fi.Size())
	openedFile.Read(bytesData)
	openedFile.Close()

	return bytesData, nil
}

// IsSafeFile : checks if file is to be encrypted or not
func IsSafeFile(a *string) bool {
	for _, b := range SafeFiles {
		if strings.Contains(*a, b) {
			return true
		}
	}
	return false
}
