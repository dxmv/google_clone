package main

import (
	"encoding/json"

	badger "github.com/dgraph-io/badger/v4"
)

const DB_PATH = "./tmp/badger"

// openDB initializes and retqurns a database connection
func openDB() (*badger.DB, error) {
	// Open the Badger database located in the /tmp/badger directory.
	// It is created if it doesn't exist.
	opts := badger.DefaultOptions(DB_PATH)
	opts.Logger = nil // Disable logging
	db, err := badger.Open(opts)

	return db, err

}

// savePosting saves a posting to the database
func savePostings(db *badger.DB, postings map[string][]Posting) error {
	err := db.Update(func(txn *badger.Txn) error {
		for term, postings := range postings {
			postingsBytes, err := json.Marshal(postings)
			if err != nil {
				return err
			}
			err = txn.Set([]byte(term), postingsBytes)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

// saveDocMeta saves document metadata to the database
func saveDocMeta(db *badger.DB, docID []byte, docMeta DocMeta) error {
	err := db.Update(func(txn *badger.Txn) error {
		docMetaBytes, err := json.Marshal(docMeta)
		if err != nil {
			return err
		}
		err = txn.Set(docID, docMetaBytes)
		return err
	})
	return err
}

func getPostings(db *badger.DB, term string) []Posting {
	postings := []Posting{}
	db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(term))
		if err != nil {
			return err
		}

		// Get the value from the item
		err = item.Value(func(val []byte) error {
			// Unmarshal the JSON bytes back to postings slice
			return json.Unmarshal(val, &postings)
		})
		return err
	})
	return postings
}
