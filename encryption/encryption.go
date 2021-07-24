package encryption

import (
	"math/rand"
	"time"
)

// SIZE : size of bytes we would like to encrypt at one time
const SIZE uint8 = 16

// GenKey : generates a random Bytes array of size SIZE
func GenKey() *[]byte {
	rand.Seed(time.Now().UnixNano())
	key := make([]byte, SIZE)
	rand.Read(key)
	return &key
}
