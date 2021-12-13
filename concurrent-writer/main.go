package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
)

func main() {

	dataCh := make(chan int)

	doneCh := make(chan struct{})

	file, err := os.OpenFile("text.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
	if err != nil {
		log.Printf("failed to open file %v", err)
	}

	wg := new(sync.WaitGroup)

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go producer(wg, dataCh)
	}

	go writer(file, dataCh, doneCh)
	wg.Wait()

	go func() {
		for d := range doneCh {
			log.Printf("don ch %v", d)
		}
		close(doneCh)
	}()

	close(dataCh)

}

func producer(wg *sync.WaitGroup, dataCh chan int) {
	defer wg.Done()
	num := rand.Intn(999)
	dataCh <- num
}

func writer(f *os.File, dataCh chan int, doneCh chan struct{}) {
	for data := range dataCh {
		_, err := fmt.Fprintln(f, data)
		if err != nil {
			log.Printf("failed to write data to file %v", err)
			f.Close()
			doneCh <- struct{}{}
			return

		}
	}
	f.Close()

	doneCh <- struct{}{}
}
