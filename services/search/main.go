package main

import (
	"fmt"

	shared "github.com/dxmv/google_clone/shared"
)

func main() {
	storage := shared.NewStorage(shared.NewMinoMongoCorpus())
	results := search("hello", storage, 0, 0)
	fmt.Println(results)
}
