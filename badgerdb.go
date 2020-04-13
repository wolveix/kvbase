package kvbase

import (
	"encoding/json"
	"errors"
	"github.com/dgraph-io/badger/v2"
	"strings"
)

// BadgerBackend acts as a wrapper around a Backend interface
type BadgerBackend struct {
	Backend
	Connection *badger.DB
	Memory     bool
	Source     string
}

// NewBadgerDB initialises a new database using the BadgerDB driver
func NewBadgerDB(source string, memory bool) (Backend, error) {
	if source == "" {
		source = "data"
	}

	opts := badger.DefaultOptions(source)
	opts.BypassLockGuard = true
	opts.Logger = nil
	opts.SyncWrites = true

	if memory {
		opts.InMemory = true
		opts.Dir = ""
		opts.ValueDir = ""
	}

	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}

	database := BadgerBackend{
		Connection: db,
		Memory:     memory,
		Source:     source,
	}

	return &database, nil
}

// Count returns the total number of records inside of the provided bucket
func (database *BadgerBackend) Count(bucket string) (int, error) {
	db := database.Connection
	counter := 0

	return counter, db.View(func(txn *badger.Txn) error {
		prefix := []byte(bucket + "_")
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			counter++
		}
		return nil
	})
}

// Create inserts a record into the backend
func (database *BadgerBackend) Create(bucket string, key string, model interface{}) error {
	if _, err := database.view(bucket, key); err == nil {
		return errors.New("key already exists")
	}

	return database.write(bucket, key, model)
}

// Delete removes a record from the backend
func (database *BadgerBackend) Delete(bucket string, key string) error {
	db := database.Connection

	if _, err := database.view(bucket, key); err != nil {
		return err
	}

	return db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(bucket + "_" + key))
	})
}

// Drop deletes a bucket (and all of its contents) from the backend
func (database *BadgerBackend) Drop(bucket string) error {
	db := database.Connection

	return db.Update(func(txn *badger.Txn) error {
		prefix := []byte(bucket + "_")
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			if err := txn.Delete(it.Item().Key()); err != nil {
				return err
			}
		}
		return nil
	})
}

// Get returns all records inside of the provided bucket
func (database *BadgerBackend) Get(bucket string, model interface{}) (*map[string]interface{}, error) {
	db := database.Connection
	results := make(map[string]interface{})

	return &results, db.View(func(txn *badger.Txn) error {
		prefix := []byte(bucket + "_")
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			key := strings.TrimPrefix(string(item.Key()), bucket+"_")

			if err := item.Value(func(value []byte) error {
				if err := json.Unmarshal(value, &model); err != nil {
					return err
				}

				results[key] = model

				return nil
			}); err != nil {
				return err
			}
		}
		return nil
	})
}

// Read returns a single struct from the provided bucket, using the provided key
func (database *BadgerBackend) Read(bucket string, key string, model interface{}) error {
	data, err := database.view(bucket, key)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &model)
}

// Update modifies an existing record from the backend, inside of the provided bucket, using the provided key
func (database *BadgerBackend) Update(bucket string, key string, model interface{}) error {
	if _, err := database.view(bucket, key); err != nil {
		return err
	}

	return database.write(bucket, key, model)
}

func (database *BadgerBackend) view(bucket string, key string) ([]byte, error) {
	db := database.Connection
	var data []byte

	return data, db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(bucket + "_" + key))
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			data = append([]byte{}, val...)

			return nil
		})
	})
}

func (database *BadgerBackend) write(bucket string, key string, model interface{}) error {
	db := database.Connection

	data, err := json.Marshal(&model)
	if err != nil {
		return err
	}

	return db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(bucket+"_"+key), data)
	})
}
