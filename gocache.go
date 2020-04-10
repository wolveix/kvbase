package kvbase

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/patrickmn/go-cache"
	"strings"
)

type GoCacheBackend struct {
	Backend
	Connection *cache.Cache
	Memory     bool
	Source     string
}

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
		Memory:     memory,
		Source:     source,
	}

	return &database, nil
}

func (database *GoCacheBackend) Count(bucket string) (int, error) {
	db := database.Connection
	counter := 0

	data := db.Items()

	for key, _ := range data {
		if strings.HasPrefix(key, bucket+"_") {
			counter++
		}
	}

	return counter, nil
}

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

func (database *GoCacheBackend) Delete(bucket string, key string) error {
	db := database.Connection

	db.Delete(bucket + "_" + key)

	if !database.Memory {
		if err := db.SaveFile(database.Source); err != nil {
			return err
		}
	}

	return nil
}

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

func (database *GoCacheBackend) Read(bucket string, key string, model interface{}) error {
	db := database.Connection

	data, found := db.Get(bucket + "_" + key)
	if !found {
		return errors.New("key could not be found")
	}

	return json.Unmarshal(data.([]byte), &model)
}

func (database *GoCacheBackend) Update(bucket string, key string, model interface{}) error {
	db := database.Connection

	if err := db.Replace(bucket+"_"+key, model, cache.NoExpiration); err != nil {
		return err
	}

	if !database.Memory {
		if err := db.SaveFile(database.Source); err != nil {
			return err
		}
	}

	return nil
}
