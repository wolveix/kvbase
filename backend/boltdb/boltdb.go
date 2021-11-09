package kvbaseBackendBoltDB

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/Wolveix/kvbase"
	"github.com/boltdb/bolt"
)

type backend struct {
	kvbase.Backend
	Connection *bolt.DB
	Memory     bool
	Source     string
}

func init() {
	store := backend{
		Connection: nil,
		Memory:     false,
		Source:     "data.db",
	}

	if err := kvbase.Register("boltdb", &store); err != nil {
		panic(err)
	}
}

// Initialize initialises a new store using the BoltDB backend
func (store *backend) Initialize(source string, memory bool) error {
	if memory {
		return errors.New("kvbase: boltdb doesn't support memory-only")
	}

	if source == "" {
		source = "data.db"
	}

	db, err := bolt.Open(source, 0600, &bolt.Options{Timeout: 1 * time.Second})
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

	if err := store.checkBucket(bucket); err != nil {
		return 0, err
	}

	return counter, db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))

		counter = b.Stats().KeyN

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

	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))

		return b.Delete([]byte(key))
	})
}

// Drop deletes a bucket (and all of its contents) from the backend
func (store *backend) Drop(bucket string) error {
	db := store.Connection

	return db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte(bucket))
	})
}

// Get returns all records inside of the provided bucket
func (store *backend) Get(bucket string, model interface{}) (*map[string]interface{}, error) {
	db := store.Connection
	results := make(map[string]interface{})

	err := store.checkBucket(bucket)
	if err != nil {
		return nil, err
	}

	return &results, db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))

		return b.ForEach(func(key, value []byte) error {
			if err := json.Unmarshal(value, &model); err != nil {
				return err
			}

			results[string(key)] = model

			return nil
		})
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

func (store *backend) checkBucket(bucket string) error {
	db := store.Connection

	return db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(bucket)); err != nil {
			return err
		}

		return nil
	})
}

func (store *backend) view(bucket, key string) ([]byte, error) {
	db := store.Connection
	var data []byte

	if err := store.checkBucket(bucket); err != nil {
		return nil, err
	}

	return data, db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))

		data = b.Get([]byte(key))

		if data == nil {
			return errors.New("key does not exist")
		}

		return nil
	})
}

func (store *backend) write(bucket, key string, model interface{}) error {
	db := store.Connection

	data, err := json.Marshal(&model)
	if err != nil {
		return err
	}

	if err = store.checkBucket(bucket); err != nil {
		return err
	}

	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))

		return b.Put([]byte(key), data)
	})
}
