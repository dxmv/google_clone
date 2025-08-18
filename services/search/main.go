package main

import (
	"fmt"
	"time"

	shared "github.com/dxmv/google_clone/shared"
)

func main() {
	t := time.Now()
	storage := shared.NewStorage(shared.NewMinoMongoCorpus())
	results := search("logic math", storage, 0, 0)
	fmt.Println(len(results))
	fmt.Println(time.Since(t))

}
