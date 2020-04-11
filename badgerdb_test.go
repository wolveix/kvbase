package kvbase_test

import (
	"github.com/Wolveix/kvbase"
	"os"
	"testing"
)

func BadgerDBInit() (kvbase.Backend, error) {
	_ = os.RemoveAll("testbadgerdbdata/")

	if database, err := kvbase.NewBadgerDB("testbadgerdbdata", false); err != nil {
		return nil, err
	} else {
		return database, nil
	}
}

func TestBadgerDBBackend(t *testing.T) {
	database, err := BadgerDBInit()
	if err != nil {
		t.Fatal("Error on database initialisation:", err)
	}

	exampleModel := model{
		"John Smith",
	}

	if err := database.Create("bucket", "key", &exampleModel); err != nil {
		t.Fatal("Error on record creation:", err)
	}

	if _, err = kvbase.NewBadgerDB("testbadgerdbdata", false); err == nil {
		t.Fatal("Expected error on database initialisation:", err)
	}
}

func TestBadgerDBBackendMem(t *testing.T) {
	_ = os.RemoveAll("testbadgerdbdata/")
	exampleModel := model{
		"John Smith",
	}

	database, err := kvbase.NewBadgerDB("testbadgerdbdata", true)
	if err != nil {
		t.Fatal("Error on database initialisation:", err)
	}

	if err := database.Create("bucket", "key", &exampleModel); err != nil {
		t.Fatal("Error on record creation:", err)
	}

	database, err = kvbase.NewBadgerDB("testbadgerdbdata", true)
	if err != nil {
		t.Fatal("Error on database initialisation:", err)
	}

	if err = database.Read("bucket", "key", &exampleModel); err == nil {
		t.Fatal("Error expected for missing key.")
	}
}

func TestBadgerDBBackend_Count(t *testing.T) {
	database, err := BadgerDBInit()
	if err != nil {
		t.Fatal("Error on database initialisation:", err)
	}

	Count(t, database)
}

func TestBadgerDBBackend_Create(t *testing.T) {
	database, err := BadgerDBInit()
	if err != nil {
		t.Fatal("Error on database initialisation:", err)
	}

	Create(t, database)
}

func TestBadgerDBBackend_Delete(t *testing.T) {
	database, err := BadgerDBInit()
	if err != nil {
		t.Fatal("Error on database initialisation:", err)
	}

	Delete(t, database)
}

func TestBadgerDBBackend_Get(t *testing.T) {
	database, err := BadgerDBInit()
	if err != nil {
		t.Fatal("Error on database initialisation:", err)
	}

	Get(t, database)
}

func TestBadgerDBBackend_Read(t *testing.T) {
	database, err := BadgerDBInit()
	if err != nil {
		t.Fatal("Error on database initialisation:", err)
	}

	Read(t, database)
}

func TestBadgerDBBackend_Update(t *testing.T) {
	database, err := BadgerDBInit()
	if err != nil {
		t.Fatal("Error on database initialisation:", err)
	}

	Update(t, database)
}
