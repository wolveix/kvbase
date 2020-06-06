package kvbaseBackendLevelDB

import (
	"encoding/json"
	"errors"
	"github.com/Wolveix/kvbase"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"strings"
)

type backend struct {
	kvbase.Backend
	Connection *leveldb.DB
	Memory     bool
	Source     string
}

func init() {
	store := backend{
		Connection: nil,
		Memory:     false,
		Source:     "data",
	}

	if err := kvbase.Register("leveldb", &store); err != nil {
		panic(err)
	}
}

// Initialize initialises a new store using the LevelDB backend
func (store *backend) Initialize(source string, memory bool) error {
	if memory {
		return errors.New("kvbase: leveldb doesn't support memory-only")
	}

	if source == "" {
		source = "data"
	}

	db, err := leveldb.OpenFile(source, nil)
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

	iter := db.NewIterator(util.BytesPrefix([]byte(bucket+"_")), nil)
	for iter.Next() {
		counter++
	}
	iter.Release()

	if err := iter.Error(); err != nil {
		return 0, err
	}

	return counter, nil
}

// Create inserts a record into the backend
func (store *backend) Create(bucket string, key string, model interface{}) error {
	db := store.Connection

	if _, err := db.Get([]byte(bucket+"_"+key), nil); err == nil {
		return errors.New("key already exists")
	}

	data, err := json.Marshal(&model)
	if err != nil {
		return err
	}

	if err := db.Put([]byte(bucket+"_"+key), data, nil); err != nil {
		return err
	}

	return nil
}

// Delete removes a record from the backend
func (store *backend) Delete(bucket string, key string) error {
	db := store.Connection

	if _, err := db.Get([]byte(bucket+"_"+key), nil); err != nil {
		return err
	}

	if err := db.Delete([]byte(bucket+"_"+key), nil); err != nil {
		return err
	}

	return nil
}

// Drop deletes a bucket (and all of its contents) from the backend
func (store *backend) Drop(bucket string) error {
	db := store.Connection

	iter := db.NewIterator(util.BytesPrefix([]byte(bucket+"_")), nil)
	for iter.Next() {
		if err := db.Delete(iter.Key(), nil); err != nil {
			return err
		}
	}
	iter.Release()

	if err := iter.Error(); err != nil {
		return err
	}

	return nil
}

// Get returns all records inside of the provided bucket
func (store *backend) Get(bucket string, model interface{}) (*map[string]interface{}, error) {
	db := store.Connection
	results := make(map[string]interface{})

	iter := db.NewIterator(util.BytesPrefix([]byte(bucket+"_")), nil)
	for iter.Next() {
		key := strings.TrimPrefix(string(iter.Key()), bucket+"_")

		if err := json.Unmarshal(iter.Value(), &model); err != nil {
			return nil, err
		}

		results[key] = model
	}
	iter.Release()

	if err := iter.Error(); err != nil {
		return nil, err
	}

	return &results, nil
}

// Read returns a single struct from the provided bucket, using the provided key
func (store *backend) Read(bucket string, key string, model interface{}) error {
	db := store.Connection

	data, err := db.Get([]byte(bucket+"_"+key), nil)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &model)
}

// Update modifies an existing record from the backend, inside of the provided bucket, using the provided key
func (store *backend) Update(bucket string, key string, model interface{}) error {
	db := store.Connection

	if _, err := db.Get([]byte(bucket+"_"+key), nil); err != nil {
		return err
	}

	data, err := json.Marshal(&model)
	if err != nil {
		return err
	}

	if err := db.Put([]byte(bucket+"_"+key), data, nil); err != nil {
		return err
	}

	return nil
}
