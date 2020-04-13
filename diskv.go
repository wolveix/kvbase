package kvbase

import (
	"encoding/json"
	"errors"
	"github.com/peterbourgon/diskv"
	"strings"
)

// DiskvBackend acts as a wrapper around a Backend interface
type DiskvBackend struct {
	Backend
	Connection *diskv.Diskv
	Source     string
}

// NewDiskvDB initialises a new database using the DiskvDB driver
func NewDiskvDB(source string) (Backend, error) {
	if source == "" {
		source = "data"
	}

	db := diskv.New(diskv.Options{
		BasePath:     source,
		Transform:    func(s string) []string { return []string{} },
		CacheSizeMax: 1024 * 1024,
	})

	database := DiskvBackend{
		Connection: db,
		Source:     source,
	}

	return &database, nil
}

// Count returns the total number of records inside of the provided bucket
func (database *DiskvBackend) Count(bucket string) (int, error) {
	db := database.Connection
	counter := 0

	keys := db.KeysPrefix(bucket+"_", nil)
	for range keys {
		counter++
	}

	return counter, nil
}

// Create inserts a record into the backend
func (database *DiskvBackend) Create(bucket string, key string, model interface{}) error {
	db := database.Connection

	if db.Has(bucket + "_" + key) {
		return errors.New("key already exists")
	}

	data, err := json.Marshal(&model)
	if err != nil {
		return err
	}

	if err := db.Write(bucket+"_"+key, data); err != nil {
		return err
	}

	return nil
}

// Delete removes a record from the backend
func (database *DiskvBackend) Delete(bucket string, key string) error {
	db := database.Connection

	if !db.Has(bucket + "_" + key) {
		return errors.New("key doesn't exist")
	}

	if err := db.Erase(bucket + "_" + key); err != nil {
		return err
	}

	return nil
}

// Drop deletes a bucket (and all of its contents) from the backend
func (database *DiskvBackend) Drop(bucket string) error {
	db := database.Connection

	keys := db.KeysPrefix(bucket+"_", nil)
	for key := range keys {
		if err := db.Erase(key); err != nil {
			return err
		}
	}

	return nil
}

// Get returns all records inside of the provided bucket
func (database *DiskvBackend) Get(bucket string, model interface{}) (*map[string]interface{}, error) {
	db := database.Connection
	results := make(map[string]interface{})

	keys := db.KeysPrefix(bucket+"_", nil)
	for rawKey := range keys {
		value, err := db.Read(rawKey)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(value, &model); err != nil {
			return nil, err
		}

		key := strings.TrimPrefix(string(rawKey), bucket+"_")

		results[key] = model
	}

	return &results, nil
}

// Read returns a single struct from the provided bucket, using the provided key
func (database *DiskvBackend) Read(bucket string, key string, model interface{}) error {
	db := database.Connection

	data, err := db.Read(bucket + "_" + key)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &model)
}

// Update modifies an existing record from the backend, inside of the provided bucket, using the provided key
func (database *DiskvBackend) Update(bucket string, key string, model interface{}) error {
	db := database.Connection

	if !db.Has(bucket + "_" + key) {
		return errors.New("key doesn't exist")
	}

	data, err := json.Marshal(&model)
	if err != nil {
		return err
	}

	if err := db.Write(bucket+"_"+key, data); err != nil {
		return err
	}

	return nil
}
