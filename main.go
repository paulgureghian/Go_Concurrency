package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var cache = map[int]Book{}
var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

func main() {
	wg := &sync.WaitGroup{}
	m := &sync.RWMutex{}

	// Create channels
	cacheCh := make(chan Book)
	dbCh := make(chan Book)

	for i := 0; i < 10; i++ {

		id := rnd.Intn(10) + 1
		wg.Add(2)

		go func(id int, wg *sync.WaitGroup, m *sync.RWMutex, ch chan<- Book) {
			if b, ok := queryCache(id, m); ok {

				ch <- b
			}
			wg.Done()
		}(id, wg, m, cacheCh)

		go func(id int, wg *sync.WaitGroup, m *sync.RWMutex, ch chan<- Book) {
			if b, ok := queryDatabase(id, m); ok {

				m.Lock()
				cache[id] = b
				m.Unlock()
				ch <- b
			}
			wg.Done()
		}(id, wg, m, dbCh)

		// Create a goroutine to recieve massages from the channels
		go func(cacheCh, dbCh <-chan Book) {

			select {

			case b := <-cacheCh:

				fmt.Println("Found in cache")
				fmt.Println(b)
				<-dbCh

			case b := <-dbCh:

				fmt.Printf("\n")
				fmt.Println("Found in database")
				fmt.Println(b)
			}

		}(cacheCh, dbCh)

		time.Sleep(150 * time.Millisecond)
	}

	wg.Wait()
}

// Add query functions
func queryCache(id int, m *sync.RWMutex) (Book, bool) {

	m.RLock()
	b, ok := cache[id]
	m.RUnlock()

	return b, ok
}

func queryDatabase(id int, m *sync.RWMutex) (Book, bool) {

	time.Sleep(100 * time.Millisecond)

	for _, b := range books {

		if b.ID == id {

			m.Lock()
			cache[id] = b
			m.Unlock()

			return b, true
		}
	}

	return Book{}, false
}
