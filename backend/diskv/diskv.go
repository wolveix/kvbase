package kvbaseBackendDiskv

import (
	"encoding/json"
	"errors"
	"github.com/Wolveix/kvbase"
	"github.com/peterbourgon/diskv"
	"strings"
)

type backend struct {
	kvbase.Backend
	Connection *diskv.Diskv
	Memory     bool
	Source     string
}

func init() {
	store := backend{
		Connection: nil,
		Memory:     false,
		Source:     "data",
	}

	if err := kvbase.Register("diskv", &store); err != nil {
		panic(err)
	}
}

// Initialize initialises a new store using the Diskv backend
func (store *backend) Initialize(source string, memory bool) error {
	if memory {
		return errors.New("kvbase: diskv doesn't support memory-only")
	}

	if source == "" {
		source = "data"
	}

	db := diskv.New(diskv.Options{
		BasePath:     source,
		Transform:    func(s string) []string { return []string{} },
		CacheSizeMax: 1024 * 1024,
	})

	store.Connection = db
	store.Memory = memory
	store.Source = source

	return nil
}

// Count returns the total number of records inside of the provided bucket
func (store *backend) Count(bucket string) (int, error) {
	db := store.Connection
	counter := 0

	keys := db.KeysPrefix(bucket+"_", nil)
	for range keys {
		counter++
	}

	return counter, nil
}

// Create inserts a record into the backend
func (store *backend) Create(bucket string, key string, model interface{}) error {
	db := store.Connection

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
func (store *backend) Delete(bucket string, key string) error {
	db := store.Connection

	if !db.Has(bucket + "_" + key) {
		return errors.New("key doesn't exist")
	}

	if err := db.Erase(bucket + "_" + key); err != nil {
		return err
	}

	return nil
}

// Drop deletes a bucket (and all of its contents) from the backend
func (store *backend) Drop(bucket string) error {
	db := store.Connection

	keys := db.KeysPrefix(bucket+"_", nil)
	for key := range keys {
		if err := db.Erase(key); err != nil {
			return err
		}
	}

	return nil
}

// Get returns all records inside of the provided bucket
func (store *backend) Get(bucket string, model interface{}) (*map[string]interface{}, error) {
	db := store.Connection
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
func (store *backend) Read(bucket string, key string, model interface{}) error {
	db := store.Connection

	data, err := db.Read(bucket + "_" + key)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &model)
}

// Update modifies an existing record from the backend, inside of the provided bucket, using the provided key
func (store *backend) Update(bucket string, key string, model interface{}) error {
	db := store.Connection

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
