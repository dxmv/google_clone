package shared

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"log"

	badger "github.com/dgraph-io/badger/v4"
)

const DB_PATH = "./app/badger"

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
		DB:     db,
		Corpus: corpus,
	}
}

// savePosting saves a posting to the database
func (s *Storage) SavePostings(postings map[string][]Posting) error {
	err := s.DB.Update(func(txn *badger.Txn) error {
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

func (s *Storage) SavePosting(term []byte, posting Posting) error {
	err := s.DB.Update(func(txn *badger.Txn) error {
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
	return s.Corpus.GetMetadata(context.Background(), docID)
}

func (s *Storage) ListMetadata() ([]DocMetadata, error) {
	return s.Corpus.ListMetadata(context.Background())
}

func (s *Storage) GetHTML(docID string) ([]byte, error) {
	return s.Corpus.GetHTML(context.Background(), docID)
}

func (s *Storage) GetPostings(term string) []Posting {
	postings := []Posting{}
	s.DB.View(func(txn *badger.Txn) error {
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
	err := s.DB.Update(func(txn *badger.Txn) error {
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
	s.DB.View(func(txn *badger.Txn) error {
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

func (s *Storage) SaveDocLength(docID string, length uint32) error {
	err := s.DB.Update(func(txn *badger.Txn) error {
		dst := make([]byte, 4)
		binary.BigEndian.PutUint32(dst, length)
		return txn.Set([]byte(docID), dst)
	})
	return err
}

func (s *Storage) GetDocLength(docID string) (uint32, error) {
	length := 0
	err := s.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(docID))
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			length = int(binary.BigEndian.Uint32(val))
			return nil
		})
		return err
	})
	return uint32(length), err
}

// returns postings for a term and then returns the positions of the term in the document
func (s *Storage) GetPositions(term string, docID string) []int {
	postings := s.GetPostings(term)
	for _, posting := range postings {
		if string(posting.DocID) == docID {
			return posting.Positions
		}
	}
	return []int{}
}
