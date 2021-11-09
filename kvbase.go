package kvbase

import (
	"errors"
	"sort"
)

// Backend implements all common backend functions
type Backend interface {
	Count(bucket string) (int, error)
	Create(bucket string, key string, model interface{}) error
	Delete(bucket string, key string) error
	Drop(bucket string) error
	Get(bucket string, model interface{}) (*map[string]interface{}, error)
	Initialize(source string, memory bool) error
	Read(bucket string, key string, model interface{}) error
	Update(bucket string, key string, model interface{}) error
}

var (
	backends = make(map[string]Backend)
)

func Backends() []string {
	list := make([]string, 0, len(backends))

	for name := range backends {
		list = append(list, name)
	}

	sort.Strings(list)
	return list
}

func New(backend, source string, memory bool) (Backend, error) {
	store := backends[backend]

	if store == nil {
		return nil, errors.New("kvbase: " + backend + " backend not registered")
	}

	if err := store.Initialize(source, memory); err != nil {
		return nil, err
	}

	return store, nil
}

func Register(name string, backend Backend) error {
	if backend == nil {
		return errors.New("kvbase: backend is nil")
	}

	if backends[name] != nil {
		return errors.New("kvbase: backend already registered")
	}

	backends[name] = backend

	return nil
}
