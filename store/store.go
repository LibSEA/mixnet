package store

import (
	"errors"
	"fmt"
	"time"

	badger "github.com/dgraph-io/badger/v4"
)

type Store struct {
	db *badger.DB
}

type Storage interface {
	Get(key []byte) ([]byte, error)
	Put(key []byte, value []byte, ttl time.Duration) error
	Update(key []byte, ttl time.Duration) error
}

var ErrKeyMissing = errors.New("key does not exist")

func Open(path string) (*Store, error) {
	db, err := badger.Open(badger.DefaultOptions(path))
	if err != nil {
		return nil, fmt.Errorf("failed to open database at path %s", path)
	}

	return &Store{db}, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}
