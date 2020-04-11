package kvbase

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/patrickmn/go-cache"
	"strings"
)

// GoCacheBackend acts as a wrapper around a Backend interface
type GoCacheBackend struct {
	Backend
	Connection *cache.Cache
	Driver string
	Memory     bool
	Source     string
}

// NewGoCache initialises a new database using the GoCache driver
func NewGoCache(source string, memory bool) (Backend, error) {
	if source == "" {
		source = "data"
	}

	db := cache.New(cache.NoExpiration, 0)

	if !memory {
		if err := db.LoadFile(source); err != nil {
			fmt.Println("Creating new database...")
		}
	}

	database := GoCacheBackend{
		Connection: db,
		Driver: 	"gocache",
		Memory:     memory,
		Source:     source,
	}

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

	if !database.Memory {
		if err := db.SaveFile(database.Source); err != nil {
			return err
		}
	}

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

	if !database.Memory {
		if err := db.SaveFile(database.Source); err != nil {
			return err
		}
	}

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

	if !database.Memory {
		if err := db.SaveFile(database.Source); err != nil {
			return err
		}
	}

	return nil
}
