package main

import (
	"Ransomware/file"
	"fmt"
	"os"
)

func main() {
	key := []byte("PeterBalivet2347")
	file, err := file.Decrypt("test.mp4.ec", key)

	if err != nil {
		os.Exit(1)
	}

	fmt.Println(file)

}
