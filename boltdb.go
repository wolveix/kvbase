package kvbase

import (
	"encoding/json"
	"errors"
	"github.com/boltdb/bolt"
	"time"
)

type BoltBackend struct {
	Backend
	Connection *bolt.DB
}

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
	}

	return &database, nil
}

func (database *BoltBackend) Count(bucket string) (int, error) {
	db := database.Connection
	counter := 0

	err := database.checkBucket(bucket)
	if err != nil {
		return 0, err
	}

	return counter, db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))

		counter = b.Stats().KeyN

		return nil
	})
}

func (database *BoltBackend) Create(bucket string, key string, model interface{}) error {
	if _, err := database.view(bucket, key); err == nil {
		return errors.New("key already exists")
	}

	return database.write(bucket, key, model)
}

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
			err := json.Unmarshal(value, &model)
			if err != nil {
				return err
			}

			results[string(key)] = model

			return nil
		})
	})
}

func (database *BoltBackend) Read(bucket string, key string, model interface{}) error {
	data, err := database.view(bucket, key)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &model)
}

func (database *BoltBackend) Update(bucket string, key string, model interface{}) error {
	if _, err := database.view(bucket, key); err != nil {
		return err
	}

	return database.write(bucket, key, model)
}

func (database *BoltBackend) checkBucket(bucket string) error {
	db := database.Connection

	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}

		return nil
	})
}

func (database *BoltBackend) view(bucket string, key string) ([]byte, error) {
	db := database.Connection
	var data []byte

	err := database.checkBucket(bucket)
	if err != nil {
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

	err = database.checkBucket(bucket)
	if err != nil {
		return err
	}

	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))

		return b.Put([]byte(key), data)
	})
}