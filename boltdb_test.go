package kvbase_test

import (
	"github.com/Wolveix/kvbase"
	"os"
	"testing"
)

func BoltDBInit() (kvbase.Backend, error) {
	_ = os.Remove("testboltdbdata")

	if database, err := kvbase.NewBoltDB("testboltdbdata"); err != nil {
		return nil, err
	} else {
		return database, nil
	}
}

func TestBoltDBBackend(t *testing.T) {
	database, err := BoltDBInit()
	if err != nil {
		t.Fatal("Error on database initialisation:", err)
	}

	exampleModel := model{
		"John Smith",
	}

	if err := database.Create("bucket", "key", &exampleModel); err != nil {
		t.Fatal("Error on record creation:", err)
	}

	if _, err = kvbase.NewBoltDB("testboltdbdata"); err == nil {
		t.Fatal("Expected error on database initialisation:", err)
	}
}

func TestBoltDBBackend_Count(t *testing.T) {
	database, err := BoltDBInit()
	if err != nil {
		t.Fatal("Error on database initialisation:", err)
	}

	Count(t, database)
}

func TestBoltDBBackend_Create(t *testing.T) {
	database, err := BoltDBInit()
	if err != nil {
		t.Fatal("Error on database initialisation:", err)
	}

	Create(t, database)
}

func TestBoltDBBackend_Delete(t *testing.T) {
	database, err := BoltDBInit()
	if err != nil {
		t.Fatal("Error on database initialisation:", err)
	}

	Delete(t, database)
}

func TestBoltDBBackend_Get(t *testing.T) {
	database, err := BoltDBInit()
	if err != nil {
		t.Fatal("Error on database initialisation:", err)
	}

	Get(t, database)
}

func TestBoltDBBackend_Read(t *testing.T) {
	database, err := BoltDBInit()
	if err != nil {
		t.Fatal("Error on database initialisation:", err)
	}

	Read(t, database)
}

func TestBoltDBBackend_Update(t *testing.T) {
	database, err := BoltDBInit()
	if err != nil {
		t.Fatal("Error on database initialisation:", err)
	}

	Update(t, database)
}
