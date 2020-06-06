package kvbase_test

import (
	"github.com/Wolveix/kvbase"
	_ "github.com/Wolveix/kvbase/backend/badgerdb"
	_ "github.com/Wolveix/kvbase/backend/bboltdb"
	_ "github.com/Wolveix/kvbase/backend/bitcask"
	_ "github.com/Wolveix/kvbase/backend/boltdb"
	_ "github.com/Wolveix/kvbase/backend/diskv"
	_ "github.com/Wolveix/kvbase/backend/go-cache"
	_ "github.com/Wolveix/kvbase/backend/leveldb"
	"github.com/Wolveix/kvbase/pkg/kvbaseBackendTest"
	"testing"
)

type backend struct {
	kvbase.Backend
	Connection bool
	Memory     bool
	Source     string
}

func TestBackends(t *testing.T) {
	for _, backend := range kvbase.Backends() {
		kvbaseBackendTest.RunTests(t, backend, "testData", false)
	}
}

func TestGetBackends(t *testing.T) {
	backends := kvbase.Backends()
	if len(backends) != 7 {
		t.Fatal("Expected 1 backend")
	}
}

func TestNewNonRegisteredBackend(t *testing.T) {
	if _, err := kvbase.New("kvbasemissingtestdb", "data", false); err == nil {
		t.Fatal("Expected 'backend not registered' error")
	}
}

func TestRegister(t *testing.T) {
	store := backend{}

	if err := kvbase.Register("", nil); err == nil {
		t.Fatal("Expected 'kvbase: backend is nil'")
	}

	if err := kvbase.Register("kvbasetestdb", store); err != nil {
		t.Fatal("Unexpected error")
	}

	if err := kvbase.Register("kvbasetestdb", store); err == nil {
		t.Fatal("Expected 'backend already registered' error")
	}
}