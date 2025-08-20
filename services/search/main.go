package main

// shared "github.com/dxmv/google_clone/shared"

func main() {
	// // t := time.Now()
	// storage := shared.NewStorage(shared.NewMinoMongoCorpus())
	// // results := search("logic math", storage, 0, 0)
	// // fmt.Println(len(results))
	// // fmt.Println(time.Since(t))
	// startServer(storage)
	cache := NewLruCache[int](1)
	cache.Put("1", 1)
	cache.Put("2", 2)
	cache.Put("3", 3)
	cache.Put("4", 4)
	cache.Put("5", 5)
	cache.Put("4", 6)
	cache.printList()
}
