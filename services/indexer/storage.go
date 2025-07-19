package main

import (
	"encoding/json"
	"fmt"

	badger "github.com/dgraph-io/badger/v4"
)

const DB_PATH = "./tmp/badger"

// openDB initializes and returns a database connection
func openDB() (*badger.DB, error) {
	// Open the Badger database located in the /tmp/badger directory.
	// It is created if it doesn't exist.
	db, err := badger.Open(badger.DefaultOptions(DB_PATH))

	return db, err

}

// savePosting saves a posting to the database
func savePosting(db *badger.DB, term string, posting Posting) error {
	// TODO: Implementation for saving posting to database
	return nil
}

// saveDocMeta saves document metadata to the database
func saveDocMeta(db *badger.DB, docID int64, docMeta DocMeta) error {
	err := db.Update(func(txn *badger.Txn) error {
		docMetaBytes, err := json.Marshal(docMeta)
		if err != nil {
			return err
		}
		err = txn.Set([]byte(fmt.Sprintf("docmeta:%d", docID)), docMetaBytes)
		return err
	})
	return err
}
