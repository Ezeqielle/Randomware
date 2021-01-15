package main

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"os"
)

//SIZE is the size of bytes we would like to encrypt at the time
const SIZE uint8 = 16

//MaxSize is the maximun allowed Bytes size to encrypt
const MaxSize int = 104857600

func padBytes(data *[]byte) {
	dataLen := uint(len(*data))
	maxBytes := (((dataLen / uint(SIZE)) + 1) * uint(SIZE)) - 1
	for i := dataLen; i < maxBytes; i++ {
		*data = append(*data, 0x00)
	}
	*data = append(*data, byte(maxBytes-dataLen))
}

func unpadBytes(data []byte) []byte {
	dataLen := uint(len(data))
	if data[dataLen-2] == 0 || data[dataLen-1] == 0 {
		return data[0 : dataLen-uint(data[dataLen-1])-1]
	}
	return data
}

func fileEncrypt(fileSrc string, key []byte) (string, error) {

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
		encrypt(&data, &key, &counter)
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
		encrypt(&data, &key, &counter)
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

func encrypt(block *[]byte, key *[]byte, counter *[]byte) ([]byte, error) {
	cipherKey, err := aes.NewCipher([]byte(*key))
	if err != nil {
		return *block, err
	}
	if len(*block)%int(SIZE) != 0 {
		padBytes(block)
	}
	ciphertext := make([]byte, len(*block))
	stream := cipher.NewCTR(cipherKey, *counter)
	stream.XORKeyStream(ciphertext, *block)
	// The IV needs to be unique, but not secure.

	*block = ciphertext
	return ciphertext, nil
}

func fileDecrypt(fileSrc string, key []byte) (string, error) {

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
		decrypt(&data, &key, &counter)
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
		decrypt(&data, &key, &counter)
		data = unpadBytes(data)
		_, err = decryptedFile.Write(data)
		if err != nil {
			return "", err
		}

		decryptedFile.Sync()
	}

	decryptedFile.Close()

	os.Remove(fileSrc)

	return decryptedFileName, nil
}

func decrypt(block *[]byte, key *[]byte, counter *[]byte) ([]byte, error) {
	cipherKey, err := aes.NewCipher([]byte(*key))
	if err != nil {
		return *block, err
	}
	if len(*block)%int(SIZE) != 0 {
		padBytes(block)
	}
	plaintext := make([]byte, len(*block))
	stream := cipher.NewCTR(cipherKey, *counter)
	stream.XORKeyStream(plaintext, *block)
	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.

	*block = plaintext
	return plaintext, nil
}

func main() {
	key := []byte("PeterBalivet2347")
	file, err := fileEncrypt("test.MOV", key)

	if err != nil {
		os.Exit(1)
	}

	fmt.Println(file)

}
