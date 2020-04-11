package kvbase

import (
	"encoding/json"
	"errors"
	"github.com/patrickmn/go-cache"
	"log"
	"strings"
	"sync"
)

// GoCacheBackend acts as a wrapper around a Backend interface
type GoCacheBackend struct {
	Backend
	Connection *cache.Cache
	Memory     bool
	Mux        sync.RWMutex
	Source     string
}

// NewGoCache initialises a new database using the GoCache driver
func NewGoCache(source string, memory bool) (Backend, error) {
	if source == "" {
		source = "data"
	}

	db := cache.New(cache.NoExpiration, 0)

	if !memory {
		_ = db.LoadFile(source)
	}

	database := GoCacheBackend{
		Connection: db,
		Memory:     memory,
		Source:     source,
	}

	database.save()

	return &database, nil
}

// Count returns the total number of records inside of the provided bucket
func (database *GoCacheBackend) Count(bucket string) (int, error) {
	db := database.Connection
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
func (database *GoCacheBackend) Create(bucket string, key string, model interface{}) error {
	db := database.Connection

	data, err := json.Marshal(&model)
	if err != nil {
		return err
	}

	if err = db.Add(bucket+"_"+key, data, cache.NoExpiration); err != nil {
		return err
	}

	database.save()

	return nil
}

// Delete removes a record from the backend
func (database *GoCacheBackend) Delete(bucket string, key string) error {
	db := database.Connection

	_, found := db.Get(bucket + "_" + key)
	if !found {
		return errors.New("key could not be found")
	}

	db.Delete(bucket + "_" + key)

	database.save()

	return nil
}

// Drop deletes a bucket (and all of its contents) from the backend
func (database *GoCacheBackend) Drop(bucket string) error {
	db := database.Connection

	data := db.Items()

	for key := range data {
		if strings.HasPrefix(key, bucket+"_") {
			db.Delete(key)
		}
	}

	database.save()

	return nil
}

// Get returns all records inside of the provided bucket
func (database *GoCacheBackend) Get(bucket string, model interface{}) (*map[string]interface{}, error) {
	db := database.Connection
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
func (database *GoCacheBackend) Read(bucket string, key string, model interface{}) error {
	db := database.Connection

	data, found := db.Get(bucket + "_" + key)
	if !found {
		return errors.New("key could not be found")
	}

	return json.Unmarshal(data.([]byte), &model)
}

// Update modifies an existing record from the backend, inside of the provided bucket, using the provided key
func (database *GoCacheBackend) Update(bucket string, key string, model interface{}) error {
	db := database.Connection

	data, err := json.Marshal(&model)
	if err != nil {
		return err
	}

	if err := db.Replace(bucket+"_"+key, data, cache.NoExpiration); err != nil {
		return err
	}

	database.save()

	return nil
}

func (database *GoCacheBackend) save() {
	if !database.Memory {
		database.Mux.RLock()
		if err := database.Connection.SaveFile(database.Source); err != nil {
			log.Fatal(err)
		}
		database.Mux.RUnlock()
	}
}
