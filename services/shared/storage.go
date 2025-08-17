package shared

import (
	"context"
	"encoding/json"
	"log"

	badger "github.com/dgraph-io/badger/v4"
)

const DB_PATH = "/Users/dimitrijestepanovic/Projects/google_clone/services/indexer/tmp/badger"

// openDB initializes and retqurns a database connection
func openDB() (*badger.DB, error) {
	// Open the Badger database located in the /tmp/badger directory.
	// It is created if it doesn't exist.
	opts := badger.DefaultOptions(DB_PATH)
	opts.Logger = nil // Disable logging
	db, err := badger.Open(opts)

	return db, err
}

func NewStorage(corpus Corpus) *Storage {
	db, err := openDB()
	if err != nil {
		log.Fatal(err)
	}
	return &Storage{
		db:     db,
		corpus: corpus,
	}
}

// savePosting saves a posting to the database
func (s *Storage) savePostings(postings map[string][]Posting) error {
	err := s.db.Update(func(txn *badger.Txn) error {
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

func (s *Storage) savePosting(term []byte, posting Posting) error {
	err := s.db.Update(func(txn *badger.Txn) error {
		postingBytes, err := json.Marshal(posting)
		if err != nil {
			return err
		}
		err = txn.Set(term, postingBytes)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

// saveMetadata saves document metadata to the database
func (s *Storage) GetMetadata(docID string) (DocMetadata, error) {
	return s.corpus.GetMetadata(context.Background(), docID)
}

func (s *Storage) ListMetadata() ([]DocMetadata, error) {
	return s.corpus.ListMetadata(context.Background())
}

func (s *Storage) GetHTML(docID string) ([]byte, error) {
	return s.corpus.GetHTML(context.Background(), docID)
}

func (s *Storage) GetPostings(term string) []Posting {
	postings := []Posting{}
	s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(term))
		if err != nil {
			log.Println("Error getting postings: ", err)
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

func (s *Storage) SaveStats(stats Stats) error {
	err := s.db.Update(func(txn *badger.Txn) error {
		statsBytes, err := json.Marshal(stats)
		if err != nil {
			return err
		}
		return txn.Set([]byte("stats"), statsBytes)
	})
	return err
}

func (s *Storage) GetStats() (Stats, error) {
	stats := Stats{}
	s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("stats"))
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			return json.Unmarshal(val, &stats)
		})
		return err
	})
	return stats, nil
}
