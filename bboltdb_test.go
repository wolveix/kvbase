package kvbase_test

import (
	"github.com/Wolveix/kvbase"
	"os"
	"testing"
)

func BboltDBInit() (kvbase.Backend, error) {
	_ = os.Remove("testbboltdbdata")

	if database, err := kvbase.NewBboltDB("testbboltdbdata"); err != nil {
		return nil, err
	} else {
		return database, nil
	}
}

func TestBboltDBBackend(t *testing.T) {
	database, err := BboltDBInit()
	if err != nil {
		t.Fatal("Error on database initialisation:", err)
	}

	exampleModel := model{
		"John Smith",
	}

	if err := database.Create("bucket", "key", &exampleModel); err != nil {
		t.Fatal("Error on record creation:", err)
	}

	if _, err = kvbase.NewBboltDB("testbboltdbdata"); err == nil {
		t.Fatal("Expected error on database initialisation:", err)
	}
}

func TestBboltDBBackend_Count(t *testing.T) {
	database, err := BboltDBInit()
	if err != nil {
		t.Fatal("Error on database initialisation:", err)
	}

	Count(t, database)
}

func TestBboltDBBackend_Create(t *testing.T) {
	database, err := BboltDBInit()
	if err != nil {
		t.Fatal("Error on database initialisation:", err)
	}

	Create(t, database)
}

func TestBboltDBBackend_Delete(t *testing.T) {
	database, err := BboltDBInit()
	if err != nil {
		t.Fatal("Error on database initialisation:", err)
	}

	Delete(t, database)
}

func TestBboltDBBackend_Get(t *testing.T) {
	database, err := BboltDBInit()
	if err != nil {
		t.Fatal("Error on database initialisation:", err)
	}

	Get(t, database)
}

func TestBboltDBBackend_Read(t *testing.T) {
	database, err := BboltDBInit()
	if err != nil {
		t.Fatal("Error on database initialisation:", err)
	}

	Read(t, database)
}

func TestBboltDBBackend_Update(t *testing.T) {
	database, err := BboltDBInit()
	if err != nil {
		t.Fatal("Error on database initialisation:", err)
	}

	Update(t, database)
}
