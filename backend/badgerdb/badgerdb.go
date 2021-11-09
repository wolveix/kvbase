package kvbaseBackendBadgerDB

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/Wolveix/kvbase"
	"github.com/dgraph-io/badger/v2"
)

type backend struct {
	kvbase.Backend
	Connection *badger.DB
	Memory     bool
	Source     string
}

func init() {
	store := backend{
		Connection: nil,
		Memory:     false,
		Source:     "data",
	}

	if err := kvbase.Register("badgerdb", &store); err != nil {
		panic(err)
	}
}

// Initialize initialises a new store using the BadgerDB backend
func (store *backend) Initialize(source string, memory bool) error {
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
		return err
	}

	store.Connection = db
	store.Memory = memory
	store.Source = source

	return nil
}

// Count returns the total number of records inside of the provided bucket
func (store *backend) Count(bucket string) (int, error) {
	db := store.Connection
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
func (store *backend) Create(bucket, key string, model interface{}) error {
	if _, err := store.view(bucket, key); err == nil {
		return errors.New("key already exists")
	}

	return store.write(bucket, key, model)
}

// Delete removes a record from the backend
func (store *backend) Delete(bucket, key string) error {
	db := store.Connection

	if _, err := store.view(bucket, key); err != nil {
		return err
	}

	return db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(bucket + "_" + key))
	})
}

// Drop deletes a bucket (and all of its contents) from the backend
func (store *backend) Drop(bucket string) error {
	db := store.Connection

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
func (store *backend) Get(bucket string, model interface{}) (*map[string]interface{}, error) {
	db := store.Connection
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
func (store *backend) Read(bucket, key string, model interface{}) error {
	data, err := store.view(bucket, key)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &model)
}

// Update modifies an existing record from the backend, inside of the provided bucket, using the provided key
func (store *backend) Update(bucket, key string, model interface{}) error {
	if _, err := store.view(bucket, key); err != nil {
		return err
	}

	return store.write(bucket, key, model)
}

func (store *backend) view(bucket, key string) ([]byte, error) {
	db := store.Connection
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

func (store *backend) write(bucket, key string, model interface{}) error {
	db := store.Connection

	data, err := json.Marshal(&model)
	if err != nil {
		return err
	}

	return db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(bucket+"_"+key), data)
	})
}
