package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

const bufferSize = 4096

func main() {
	fmt.Println("vim-go")
	// open file
	file, err := os.Open("1MB.db")
	defer file.Close()
	if err != nil {
		log.Printf("failed to open file %v", err)
		return
	}

	fileInfo, err := file.Stat()
	if err != nil {
		log.Printf("failed to fetch file stat %v", err)
		return
	}

	log.Printf("file size %d", fileInfo.Size())
	var numberOfWorker = 1
	if fileInfo.Size() > bufferSize {
		numberOfWorker = int(fileInfo.Size()) / bufferSize
	}
	log.Printf("number of worker %d", numberOfWorker)
	// create reader
	r := bufio.NewReader(file)
	// create scanner
	out := make([]byte, bufferSize)
	// read data concurrently

	wg := new(sync.WaitGroup)
	wg.Add(numberOfWorker)
	for i := 0; i < numberOfWorker; i++ {
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			n, err := r.Read(out)
			if err != nil {
				if err == io.EOF {
					log.Printf("end of file")
					return
				}
				log.Printf("failed to read file %v", err)
				return
			}
			log.Printf("size read %d data %s", n, string(out))
		}(wg)
	}
	wg.Wait()
}
