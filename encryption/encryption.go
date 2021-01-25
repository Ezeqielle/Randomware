package encryption

import "math/rand"

// SIZE : size of bytes we would like to encrypt at one time
const SIZE uint8 = 16

// GenKey : generates a random Bytes array of size SIZE
func GenKey() *[]byte {
	var i uint8
	key := make([]byte, SIZE)
	for i = 0; i < SIZE; i++ {
		key[i] = uint8(rand.Intn(255))
	}
	return &key
}
