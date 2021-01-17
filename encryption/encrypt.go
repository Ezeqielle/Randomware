package encryption

import (
	"crypto/aes"
	"crypto/cipher"
)

// PadBytes : pads missing bytes at the end
func PadBytes(data *[]byte) {
	dataLen := uint(len(*data))
	maxBytes := (((dataLen / uint(SIZE)) + 1) * uint(SIZE)) - 1
	for i := dataLen; i < maxBytes; i++ {
		*data = append(*data, 0x00)
	}
	*data = append(*data, byte(maxBytes-dataLen))
}

// Encrypt : encrypts chunk of data with key and counter
func Encrypt(block *[]byte, key *[]byte, counter *[]byte) ([]byte, error) {
	cipherKey, err := aes.NewCipher([]byte(*key))
	if err != nil {
		return *block, err
	}
	if len(*block)%int(SIZE) != 0 {
		PadBytes(block)
	}
	ciphertext := make([]byte, len(*block))
	stream := cipher.NewCTR(cipherKey, *counter)
	stream.XORKeyStream(ciphertext, *block)
	// The IV needs to be unique, but not secure.

	*block = ciphertext
	return ciphertext, nil
}
