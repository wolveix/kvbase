package kvbase_test

import (
	"github.com/Wolveix/kvbase"
	"os"
	"testing"
)

func NewBoltDBDisk() kvbase.Backend {
	_ = os.Remove("testboltdbdata")

	database, err := kvbase.NewBoltDB("testboltdbdata")
	if err != nil {
		panic("Error on database initialisation: " + err.Error())
	}

	return database
}

func TestBoltDBBackend(t *testing.T) {
	database := NewBoltDBDisk()

	if err := database.Create("bucket", "key", &exampleModel); err != nil {
		t.Fatal("Error on record creation:", err)
	}

	if _, err := kvbase.NewBoltDB("testboltdbdata"); err == nil {
		t.Fatal("Expected error on database initialisation:", err)
	}
}

func TestBoltDBBackend_Count(t *testing.T) {
	database := NewBoltDBDisk()
	testCount(t, database)
}

func TestBoltDBBackend_Create(t *testing.T) {
	database := NewBoltDBDisk()
	testCreate(t, database)
}

func TestBoltDBBackend_Delete(t *testing.T) {
	database := NewBoltDBDisk()
	testDelete(t, database)
}

func TestBoltDBBackend_Drop(t *testing.T) {
	database := NewBoltDBDisk()
	testDrop(t, database)
}

func TestBoltDBBackend_Get(t *testing.T) {
	database := NewBoltDBDisk()
	testGet(t, database)
}

func TestBoltDBBackend_Read(t *testing.T) {
	database := NewBoltDBDisk()
	testRead(t, database)
}

func TestBoltDBBackend_Update(t *testing.T) {
	database := NewBoltDBDisk()
	testUpdate(t, database)
}

func BenchmarkBoltDBBackend_Count(b *testing.B) {
	database := NewBoltDBDisk()
	benchmarkCount(b, database)
}

func BenchmarkBoltDBBackend_Create(b *testing.B) {
	database := NewBoltDBDisk()
	benchmarkCreate(b, database)
}

func BenchmarkBoltDBBackend_Delete(b *testing.B) {
	database := NewBoltDBDisk()
	benchmarkDelete(b, database)
}

func BenchmarkBoltDBBackend_Drop(b *testing.B) {
	database := NewBoltDBDisk()
	benchmarkDrop(b, database)
}

func BenchmarkBoltDBBackend_Get(b *testing.B) {
	database := NewBoltDBDisk()
	benchmarkGet(b, database)
}

func BenchmarkBoltDBBackend_Read(b *testing.B) {
	database := NewBoltDBDisk()
	benchmarkRead(b, database)
}

func BenchmarkBoltDBBackend_Update(b *testing.B) {
	database := NewBoltDBDisk()
	benchmarkUpdate(b, database)
}
