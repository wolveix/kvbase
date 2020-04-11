package kvbase_test

import (
	"github.com/Wolveix/kvbase"
	"os"
	"testing"
)

func GoCacheInit() (kvbase.Backend, error) {
	_ = os.Remove("testgocachedata")

	if database, err := kvbase.NewGoCache("testgocachedata", false); err != nil {
		return nil, err
	} else {
		return database, nil
	}
}

func TestGoCacheBackend(t *testing.T) {
	database, err := GoCacheInit()
	if err != nil {
		t.Fatal("Error on database initialisation:", err)
	}

	exampleModel := model{
		"John Smith",
	}

	if err := database.Create("bucket", "key", &exampleModel); err != nil {
		t.Fatal("Error on record creation:", err)
	}

	database, err = kvbase.NewGoCache("testgocachedata", false)
	if err != nil {
		t.Fatal("Error on database initialisation:", err)
	}

	if err = database.Read("bucket", "key", &exampleModel); err != nil {
		t.Fatal("Error on database read:", err)
	}
}

func TestGoCacheBackendMem(t *testing.T) {
	_ = os.Remove("testgocachedata")
	exampleModel := model{
		"John Smith",
	}

	database, err := kvbase.NewGoCache("testgocachedata", true)
	if err != nil {
		t.Fatal("Error on database initialisation:", err)
	}

	if err := database.Create("bucket", "key", &exampleModel); err != nil {
		t.Fatal("Error on record creation:", err)
	}

	database, err = kvbase.NewGoCache("testgocachedata", true)
	if err != nil {
		t.Fatal("Error on database initialisation:", err)
	}

	if err = database.Read("bucket", "key", &exampleModel); err == nil {
		t.Fatal("Error expected for missing key.")
	}
}

func TestGoCacheBackend_Count(t *testing.T) {
	database, err := GoCacheInit()
	if err != nil {
		t.Fatal("Error on database initialisation:", err)
	}

	Count(t, database)
}

func TestGoCacheBackend_Create(t *testing.T) {
	database, err := GoCacheInit()
	if err != nil {
		t.Fatal("Error on database initialisation:", err)
	}

	Create(t, database)
}

func TestGoCacheBackend_Delete(t *testing.T) {
	database, err := GoCacheInit()
	if err != nil {
		t.Fatal("Error on database initialisation:", err)
	}

	Delete(t, database)
}

func TestGoCacheBackend_Get(t *testing.T) {
	database, err := GoCacheInit()
	if err != nil {
		t.Fatal("Error on database initialisation:", err)
	}

	Get(t, database)
}

func TestGoCacheBackend_Read(t *testing.T) {
	database, err := GoCacheInit()
	if err != nil {
		t.Fatal("Error on database initialisation:", err)
	}

	Read(t, database)
}

func TestGoCacheBackend_Update(t *testing.T) {
	database, err := GoCacheInit()
	if err != nil {
		t.Fatal("Error on database initialisation:", err)
	}

	Update(t, database)
}
