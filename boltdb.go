package kvbase

import (
	"encoding/json"
	"errors"
	"github.com/boltdb/bolt"
	"time"
)

// BoltBackend acts as a wrapper around a Backend interface
type BoltBackend struct {
	Backend
	Connection *bolt.DB
	Source     string
}

// NewBoltDB initialises a new database using the BoltDB driver
func NewBoltDB(source string) (Backend, error) {
	if source == "" {
		source = "data.db"
	}

	db, err := bolt.Open(source, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}

	database := BoltBackend{
		Connection: db,
		Source:     source,
	}

	return &database, nil
}

// Count returns the total number of records inside of the provided bucket
func (database *BoltBackend) Count(bucket string) (int, error) {
	db := database.Connection
	counter := 0

	if err := database.checkBucket(bucket); err != nil {
		return 0, err
	}

	return counter, db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))

		counter = b.Stats().KeyN

		return nil
	})
}

// Create inserts a record into the backend
func (database *BoltBackend) Create(bucket string, key string, model interface{}) error {
	if _, err := database.view(bucket, key); err == nil {
		return errors.New("key already exists")
	}

	return database.write(bucket, key, model)
}

// Delete removes a record from the backend
func (database *BoltBackend) Delete(bucket string, key string) error {
	db := database.Connection

	if _, err := database.view(bucket, key); err != nil {
		return err
	}

	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))

		return b.Delete([]byte(key))
	})
}

// Drop deletes a bucket (and all of its contents) from the backend
func (database *BoltBackend) Drop(bucket string) error {
	db := database.Connection

	return db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte(bucket))
	})
}

// Get returns all records inside of the provided bucket
func (database *BoltBackend) Get(bucket string, model interface{}) (*map[string]interface{}, error) {
	db := database.Connection
	results := make(map[string]interface{})

	err := database.checkBucket(bucket)
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
func (database *BoltBackend) Read(bucket string, key string, model interface{}) error {
	data, err := database.view(bucket, key)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &model)
}

// Update modifies an existing record from the backend, inside of the provided bucket, using the provided key
func (database *BoltBackend) Update(bucket string, key string, model interface{}) error {
	if _, err := database.view(bucket, key); err != nil {
		return err
	}

	return database.write(bucket, key, model)
}

func (database *BoltBackend) checkBucket(bucket string) error {
	db := database.Connection

	return db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(bucket)); err != nil {
			return err
		}

		return nil
	})
}

func (database *BoltBackend) view(bucket string, key string) ([]byte, error) {
	db := database.Connection
	var data []byte

	if err := database.checkBucket(bucket); err != nil {
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

func (database *BoltBackend) write(bucket string, key string, model interface{}) error {
	db := database.Connection

	data, err := json.Marshal(&model)
	if err != nil {
		return err
	}

	if err = database.checkBucket(bucket); err != nil {
		return err
	}

	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))

		return b.Put([]byte(key), data)
	})
}
