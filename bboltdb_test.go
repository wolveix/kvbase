package kvbase_test

import (
	"github.com/Wolveix/kvbase"
	"os"
	"testing"
)

func NewBboltDBDisk() kvbase.Backend {
	_ = os.Remove("testbboltdbdata")

	database, err := kvbase.NewBboltDB("testbboltdbdata")
	if err != nil {
		panic("Error on database initialisation: " + err.Error())
	}

	return database
}

func TestBboltDBBackend(t *testing.T) {
	database := NewBboltDBDisk()

	if err := database.Create("bucket", "key", &exampleModel); err != nil {
		t.Fatal("Error on record creation:", err)
	}

	if _, err := kvbase.NewBboltDB("testbboltdbdata"); err == nil {
		t.Fatal("Expected error on database initialisation:", err)
	}
}

func TestBboltDBBackend_Count(t *testing.T) {
	database := NewBboltDBDisk()
	testCount(t, database)
}

func TestBboltDBBackend_Create(t *testing.T) {
	database := NewBboltDBDisk()
	testCreate(t, database)
}

func TestBboltDBBackend_Delete(t *testing.T) {
	database := NewBboltDBDisk()
	testDelete(t, database)
}

func TestBboltDBBackend_Drop(t *testing.T) {
	database := NewBboltDBDisk()
	testDrop(t, database)
}

func TestBboltDBBackend_Get(t *testing.T) {
	database := NewBboltDBDisk()
	testGet(t, database)
}

func TestBboltDBBackend_Read(t *testing.T) {
	database := NewBboltDBDisk()
	testRead(t, database)
}

func TestBboltDBBackend_Update(t *testing.T) {
	database := NewBboltDBDisk()
	testUpdate(t, database)
}

func BenchmarkBboltDBBackend_Count(b *testing.B) {
	database := NewBboltDBDisk()
	benchmarkCount(b, database)
}

func BenchmarkBboltDBBackend_Create(b *testing.B) {
	database := NewBboltDBDisk()
	benchmarkCreate(b, database)
}

func BenchmarkBboltDBBackend_Delete(b *testing.B) {
	database := NewBboltDBDisk()
	benchmarkDelete(b, database)
}

func BenchmarkBboltDBBackend_Drop(b *testing.B) {
	database := NewBboltDBDisk()
	benchmarkDrop(b, database)
}

func BenchmarkBboltDBBackend_Get(b *testing.B) {
	database := NewBboltDBDisk()
	benchmarkGet(b, database)
}

func BenchmarkBboltDBBackend_Read(b *testing.B) {
	database := NewBboltDBDisk()
	benchmarkRead(b, database)
}

func BenchmarkBboltDBBackend_Update(b *testing.B) {
	database := NewBboltDBDisk()
	benchmarkUpdate(b, database)
}
