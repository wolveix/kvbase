package kvbase

import (
	"encoding/json"
	"errors"
	"github.com/prologic/bitcask"
	"strings"
)

// BitcaskBackend acts as a wrapper around a Backend interface
type BitcaskBackend struct {
	Backend
	Connection *bitcask.Bitcask
	Source     string
}

// NewBitcaskDB initialises a new database using the BitcaskDB driver
func NewBitcaskDB(source string) (Backend, error) {
	if source == "" {
		source = "data"
	}

	db, err := bitcask.Open(source)
	if err != nil {
		return nil, err
	}

	database := BitcaskBackend{
		Connection: db,
		Source:     source,
	}

	return &database, nil
}

// Count returns the total number of records inside of the provided bucket
func (database *BitcaskBackend) Count(bucket string) (int, error) {
	db := database.Connection
	counter := 0

	return counter, db.Scan([]byte(bucket+"_"), func(key []byte) error {
		counter++
		return nil
	})
}

// Create inserts a record into the backend
func (database *BitcaskBackend) Create(bucket string, key string, model interface{}) error {
	db := database.Connection

	if db.Has([]byte(bucket + "_" + key)) {
		return errors.New("key already exists")
	}

	data, err := json.Marshal(&model)
	if err != nil {
		return err
	}

	if err := db.Put([]byte(bucket+"_"+key), data); err != nil {
		return err
	}

	return nil
}

// Delete removes a record from the backend
func (database *BitcaskBackend) Delete(bucket string, key string) error {
	db := database.Connection

	if !db.Has([]byte(bucket + "_" + key)) {
		return errors.New("key doesn't exist")
	}

	return db.Delete([]byte(bucket + "_" + key))
}

// Drop deletes a bucket (and all of its contents) from the backend
func (database *BitcaskBackend) Drop(bucket string) error {
	db := database.Connection

	var keys [][]byte
	if err := db.Scan([]byte(bucket+"_"), func(key []byte) error {
		keys = append(keys, key)
		return nil
	}); err != nil {
		return err
	}

	for _, key := range keys {
		if err := db.Delete(key); err != nil {
			return err
		}
	}

	return nil
}

// Get returns all records inside of the provided bucket
func (database *BitcaskBackend) Get(bucket string, model interface{}) (*map[string]interface{}, error) {
	db := database.Connection
	results := make(map[string]interface{})

	return &results, db.Scan([]byte(bucket+"_"), func(rawKey []byte) error {
		data, err := db.Get(rawKey)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(data, &model); err != nil {
			return err
		}

		key := strings.TrimPrefix(string(rawKey), bucket+"_")
		results[key] = model
		return nil
	})
}

// Read returns a single struct from the provided bucket, using the provided key
func (database *BitcaskBackend) Read(bucket string, key string, model interface{}) error {
	db := database.Connection

	data, err := db.Get([]byte(bucket + "_" + key))
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &model)
}

// Update modifies an existing record from the backend, inside of the provided bucket, using the provided key
func (database *BitcaskBackend) Update(bucket string, key string, model interface{}) error {
	db := database.Connection

	if !db.Has([]byte(bucket + "_" + key)) {
		return errors.New("key doesn't exist")
	}

	data, err := json.Marshal(&model)
	if err != nil {
		return err
	}

	if err := db.Put([]byte(bucket+"_"+key), data); err != nil {
		return err
	}

	return nil
}
