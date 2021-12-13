package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"sync"
)

const (
	mb = 1024 * 1024
	gb = mb * 1024
)

func main() {
	// data size
	dataCh := make(chan int) // number of bytes read from the file

	// output completed signal
	doneCh := make(chan struct{})

	// buffer size
	var bufferSize int64 = 10 * mb

	// offset
	var current int64 = 0

	// open file
	file, err := os.Open("100MB.bin")
	defer file.Close()
	if err != nil {
		log.Printf("failed to open file %v", err)
		return
	}
	// file size
	fileInfo, err := file.Stat()
	log.Printf("file size in byte : %d", fileInfo.Size())

	// create reader
	fileReader := bufio.NewReader(file)




	// output reader
	go func(dataCh chan int, doneCh chan struct{}){
		//count := 0
		//for data := range dataCh {
		//	count += data
		//}

		count := 0
		count1 := 0
		for {

			select {
				case data,open:= <- dataCh :
					count1 ++
					count += data
					if !open {
						log.Printf("number of bytes read from the file : %d %d", count, count1)
						doneCh <- struct{}{}
						return
					}

			default:

			}
		}
		close(doneCh)

	}(dataCh, doneCh)

	// wait group to complete all the task
	wg := new(sync.WaitGroup)
	wg.Add(10)
	// file reader
	for i := 1; i <= 10; i++ {
		go func(wg *sync.WaitGroup, dataCh chan int, file *os.File, fileReader *bufio.Reader, bufferSize, current int64) {
			defer wg.Done()
			reader(dataCh, file, fileReader, bufferSize, current)
		}(wg,dataCh, file, fileReader, bufferSize, current)

		current += bufferSize + 1
		log.Printf("current size %d", current)
	}
	wg.Wait()
	close(dataCh)

	<-doneCh

}

func reader( dataCh chan int, file *os.File, fileReader *bufio.Reader, bufferSize, current int64) {

	file.Seek(current, 0)
	buf := make([]byte, bufferSize)
	log.Printf("len of buf %d current %d",len(buf), current)
	for {
		n, err := fileReader.Read(buf)
		if err != nil {
			if err == io.EOF {
				log.Printf("end of file")
				return
			}
			log.Printf("failed to read file %v", err)
			break
		}

		dataCh <- n

	}
	log.Printf("ending reader : %d", current)
}
