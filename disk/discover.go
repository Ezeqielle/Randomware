package disk

import (
	"os"
)

// Discover : returns array of char with every char representing a network drive
func Discover() (r []string) {
	for _, drive := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		f, err := os.Open(string(drive) + ":\\")
		if err == nil {
			r = append(r, string(drive))
			f.Close()
		}
	}

	return
}
