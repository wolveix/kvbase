package kvbaseBackendBitcask

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/Wolveix/kvbase"
	"github.com/prologic/bitcask"
)

type backend struct {
	kvbase.Backend
	Connection *bitcask.Bitcask
	Memory     bool
	Source     string
}

func init() {
	store := backend{
		Connection: nil,
		Memory:     false,
		Source:     "data",
	}

	if err := kvbase.Register("bitcask", &store); err != nil {
		panic(err)
	}
}

// Initialize initialises a new store using the Bitcask backend
func (store *backend) Initialize(source string, memory bool) error {
	if memory {
		return errors.New("kvbase: bitcask doesn't support memory-only")
	}

	if source == "" {
		source = "data"
	}

	db, err := bitcask.Open(source)
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

	return counter, db.Scan([]byte(bucket+"_"), func(key []byte) error {
		counter++
		return nil
	})
}

// Create inserts a record into the backend
func (store *backend) Create(bucket, key string, model interface{}) error {
	db := store.Connection

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
func (store *backend) Delete(bucket, key string) error {
	db := store.Connection

	if !db.Has([]byte(bucket + "_" + key)) {
		return errors.New("key doesn't exist")
	}

	return db.Delete([]byte(bucket + "_" + key))
}

// Drop deletes a bucket (and all of its contents) from the backend
func (store *backend) Drop(bucket string) error {
	db := store.Connection

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
func (store *backend) Get(bucket string, model interface{}) (*map[string]interface{}, error) {
	db := store.Connection
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
func (store *backend) Read(bucket, key string, model interface{}) error {
	db := store.Connection

	data, err := db.Get([]byte(bucket + "_" + key))
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &model)
}

// Update modifies an existing record from the backend, inside of the provided bucket, using the provided key
func (store *backend) Update(bucket, key string, model interface{}) error {
	db := store.Connection

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
