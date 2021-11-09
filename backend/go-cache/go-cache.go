package kvbaseBackendGoCache

import (
	"encoding/json"
	"errors"
	"strings"
	"sync"

	"github.com/Wolveix/kvbase"
	"github.com/patrickmn/go-cache"
)

type backend struct {
	kvbase.Backend
	Connection *cache.Cache
	Memory     bool
	Mux        sync.RWMutex
	Source     string
}

func init() {
	store := backend{
		Connection: nil,
		Memory:     false,
		Source:     "data",
	}

	if err := kvbase.Register("go-cache", &store); err != nil {
		panic(err)
	}
}

// Initialize initialises a new store using the Go-Cache backend
func (store *backend) Initialize(source string, memory bool) error {
	if source == "" {
		source = "data"
	}

	db := cache.New(cache.NoExpiration, 0)

	if !memory {
		_ = db.LoadFile(source)
	}

	store.Connection = db
	store.Memory = memory
	store.Source = source

	if err := store.save(); err != nil {
		return err
	}

	return nil
}

// Count returns the total number of records inside of the provided bucket
func (store *backend) Count(bucket string) (int, error) {
	db := store.Connection
	counter := 0

	data := db.Items()

	for key := range data {
		if strings.HasPrefix(key, bucket+"_") {
			counter++
		}
	}

	return counter, nil
}

// Create inserts a record into the backend
func (store *backend) Create(bucket, key string, model interface{}) error {
	db := store.Connection

	data, err := json.Marshal(&model)
	if err != nil {
		return err
	}

	if err = db.Add(bucket+"_"+key, data, cache.NoExpiration); err != nil {
		return err
	}

	if err := store.save(); err != nil {
		return err
	}

	return nil
}

// Delete removes a record from the backend
func (store *backend) Delete(bucket, key string) error {
	db := store.Connection

	_, found := db.Get(bucket + "_" + key)
	if !found {
		return errors.New("key could not be found")
	}

	db.Delete(bucket + "_" + key)

	if err := store.save(); err != nil {
		return err
	}

	return nil
}

// Drop deletes a bucket (and all of its contents) from the backend
func (store *backend) Drop(bucket string) error {
	db := store.Connection

	data := db.Items()

	for key := range data {
		if strings.HasPrefix(key, bucket+"_") {
			db.Delete(key)
		}
	}

	if err := store.save(); err != nil {
		return err
	}

	return nil
}

// Get returns all records inside of the provided bucket
func (store *backend) Get(bucket string, model interface{}) (*map[string]interface{}, error) {
	db := store.Connection
	results := make(map[string]interface{})

	data := db.Items()

	for key, value := range data {
		if strings.HasPrefix(key, bucket+"_") {
			if err := json.Unmarshal(value.Object.([]byte), &model); err != nil {
				return nil, err
			}

			results[strings.TrimPrefix(key, bucket+"_")] = model
		}
	}

	return &results, nil
}

// Read returns a single struct from the provided bucket, using the provided key
func (store *backend) Read(bucket, key string, model interface{}) error {
	db := store.Connection

	data, found := db.Get(bucket + "_" + key)
	if !found {
		return errors.New("key could not be found")
	}

	return json.Unmarshal(data.([]byte), &model)
}

// Update modifies an existing record from the backend, inside of the provided bucket, using the provided key
func (store *backend) Update(bucket, key string, model interface{}) error {
	db := store.Connection

	data, err := json.Marshal(&model)
	if err != nil {
		return err
	}

	if err := db.Replace(bucket+"_"+key, data, cache.NoExpiration); err != nil {
		return err
	}

	if err := store.save(); err != nil {
		return err
	}

	return nil
}

func (store *backend) save() error {
	if !store.Memory {
		store.Mux.RLock()
		if err := store.Connection.SaveFile(store.Source); err != nil {
			return err
		}
		store.Mux.RUnlock()
	}

	return nil
}
