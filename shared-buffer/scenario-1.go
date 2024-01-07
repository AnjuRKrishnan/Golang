package main

import (
	"fmt"
	"sync"
)

const bufferSize = 10

func main() {
	var buffer []byte
	var rwMutex sync.RWMutex

	readers := make(chan struct{}, 8)
	writers := make(chan struct{}, 2)

	for i := 0; i < 8; i++ {
		readers <- struct{}{}
		go func() {
			for {
				rwMutex.RLock()
				
				_ = buffer
				rwMutex.RUnlock()
			}
		}()
	}

	for i := 0; i < 2; i++ {
		writers <- struct{}{}
		go func() {
			for {
				rwMutex.Lock()
				
				buffer = make([]byte, bufferSize)
				rwMutex.Unlock()
			}
		}()
	}

	
	select {}
}
