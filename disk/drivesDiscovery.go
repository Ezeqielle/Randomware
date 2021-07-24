package disk

import (
	"os"
)

// Discover : returns array of char with every char representing a network drive
func Discover() (r []string) {
	osFileTest := os.Getenv("SystemDrive") + "\\test.txt"
	testFile := ""
	for _, drive := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		testFile = string(drive) + ":\\test.txt"
		if osFileTest != testFile {
			f, err := os.Open(testFile)
			if err == nil {
				r = append(r, string(drive)+":\\")
				f.Close()
			}
		}

	}

	return
}
