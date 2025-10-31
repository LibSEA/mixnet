/*
mixnet - tool to create and manage LibSEA mixnets
Copyright (C) 2025  Liberatory Sofware Engineering Association

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

// Package store implements the storage for the DHT
package store

import (
	"errors"
	"fmt"
	"time"

	badger "github.com/dgraph-io/badger/v4"
)

type Options struct {
	// Frequency for GC to be requested
	GCFrequency time.Duration
}

type Store struct {
	db *badger.DB
}

type Storage interface {
	Get(out []byte, key []byte) ([]byte, error)
	Put(key []byte, value []byte, ttl time.Duration) error
}

func (s *Store) Put(key []byte, value []byte, ttl time.Duration) error {
	return s.db.Update(func(txn *badger.Txn) error {
		return txn.SetEntry(badger.NewEntry(key, value).WithTTL(ttl))
	})
}

func (s *Store) Get(out []byte, key []byte) ([]byte, error) {
	var val []byte

	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}

		val, err = item.ValueCopy(out)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return val, nil
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
