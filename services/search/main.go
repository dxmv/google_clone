package main

import "fmt"

// shared "github.com/dxmv/google_clone/shared"

func main() {
	// // t := time.Now()
	// storage := shared.NewStorage(shared.NewMinoMongoCorpus())
	// // results := search("logic math", storage, 0, 0)
	// // fmt.Println(len(results))
	// // fmt.Println(time.Since(t))
	// startServer(storage)
	cache := NewLRUCache(2)
	cache.Put(1)
	cache.Put(2)
	cache.Put(3)
	cache.Get(1)
	cache.Get(3)
	cache.Get(2)
	fmt.Println(cache.String())
	// cache.Put("2", 2)
	// cache.Put("3", 3)
	// cache.Put("4", 4)
	// cache.Put("5", 5)
	// cache.Put("4", 6)
	// cache.printList()
}
