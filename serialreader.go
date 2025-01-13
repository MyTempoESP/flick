package flick

import (
	"bufio"
	"fmt"
	"sync"
)

type SerialReader struct {
	mu       sync.Mutex
	latest   string
	stopChan chan struct{}

	Available chan struct{}
}

func NewSerialReader() *SerialReader {

	return &SerialReader{
		Available: make(chan struct{}),
		stopChan:  make(chan struct{}),
	}
}

func (sr *SerialReader) Start(scanner *bufio.Scanner) {

	go func() {

		for {
			select {
			case <-sr.stopChan:
				return
			default:
				if scanner.Scan() {

					sr.Available <- struct{}{}

					sr.mu.Lock()
					sr.latest = scanner.Text()
					sr.mu.Unlock()
				} else {

					err := scanner.Err()

					if err != nil {

						fmt.Printf("Scanner error: %v\n", err)
					}

					fmt.Println("Scanner finished (EOF/Error).")

					return
				}
			}
		}
	}()
}

func (sr *SerialReader) Wait() {

	if sr.Available == nil {

		return
	}

	<-sr.Available
}

func (sr *SerialReader) Stop() {

	close(sr.stopChan)
}

func (sr *SerialReader) GetLatest() (latest string) {

	sr.mu.Lock()

	latest = sr.latest

	sr.mu.Unlock()

	return
}
