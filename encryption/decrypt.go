package encryption

import (
	"crypto/aes"
	"crypto/cipher"
)

// UnpadBytes : pads missing bytes at the end
func UnpadBytes(data *[]byte) {
	dataLen := uint(len(*data))
	if (*data)[dataLen-2] == 0 || (*data)[dataLen-1] == 0 {
		*data = (*data)[0 : dataLen-uint((*data)[dataLen-1])-1]
	}
}

// Decrypt : decrypts chunk of data with key and counter
func Decrypt(block *[]byte, key *[]byte, counter *[]byte) ([]byte, error) {
	cipherKey, err := aes.NewCipher([]byte(*key))
	if err != nil {
		return *block, err
	}

	plaintext := make([]byte, len(*block))
	stream := cipher.NewCTR(cipherKey, *counter)
	stream.XORKeyStream(plaintext, *block)

	*block = plaintext
	UnpadBytes(block)
	return plaintext, nil
}
