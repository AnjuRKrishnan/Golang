package main

import (
    "fmt"
    "sync"
)

var (
    buffer []byte
    mutex  sync.Mutex
)

func reader(id int) {
    for {
        mutex.Lock()
        
        fmt.Printf("Reader %d read: %s\n", id, string(buffer))
        mutex.Unlock()
    
    }
}

func writer(id int) {
    for {
        mutex.Lock()
        
        buffer = []byte(fmt.Sprintf("Data written by writer %d", id))
        fmt.Printf("Writer %d wrote: %s\n", id, string(buffer))
        mutex.Unlock()
        
    }
}

func main() {
    M := 8 
    N := 16 

    
    for i := 0; i < M; i++ {
        go reader(i)
    }

    
    for i := 0; i < N; i++ {
        go writer(i)
    }

    
    select {}
}
