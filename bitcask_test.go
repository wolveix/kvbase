package kvbase_test

import (
	"github.com/Wolveix/kvbase"
	"os"
	"testing"
)

func NewBitcaskDBDisk() kvbase.Backend {
	_ = os.RemoveAll("testbitcaskdbdata/")

	database, err := kvbase.NewBitcaskDB("testbitcaskdbdata")
	if err != nil {
		panic("Error on database initialisation: " + err.Error())
	}

	return database
}

func TestBitcaskDBBackend(t *testing.T) {
	database := NewBitcaskDBDisk()

	if err := database.Create("bucket", "key", &exampleModel); err != nil {
		t.Fatal("Error on record creation:", err)
	}

	if _, err := kvbase.NewBitcaskDB("testbitcaskdbdata"); err == nil {
		t.Fatal("Expected error on database initialisation:", err)
	}
}

func TestBitcaskDBBackend_Count(t *testing.T) {
	database := NewBitcaskDBDisk()
	testCount(t, database)
}

func TestBitcaskDBBackend_Create(t *testing.T) {
	database := NewBitcaskDBDisk()
	testCreate(t, database)
}

func TestBitcaskDBBackend_Delete(t *testing.T) {
	database := NewBitcaskDBDisk()
	testDelete(t, database)
}

func TestBitcaskDBBackend_Drop(t *testing.T) {
	database := NewBitcaskDBDisk()
	testDrop(t, database)
}

func TestBitcaskDBBackend_Get(t *testing.T) {
	database := NewBitcaskDBDisk()
	testGet(t, database)
}

func TestBitcaskDBBackend_Read(t *testing.T) {
	database := NewBitcaskDBDisk()
	testRead(t, database)
}

func TestBitcaskDBBackend_Update(t *testing.T) {
	database := NewBitcaskDBDisk()
	testUpdate(t, database)
}

func BenchmarkBitcaskDBBackend_Count(b *testing.B) {
	database := NewBitcaskDBDisk()
	benchmarkCount(b, database)
}

func BenchmarkBitcaskDBBackend_Create(b *testing.B) {
	database := NewBitcaskDBDisk()
	benchmarkCreate(b, database)
}

func BenchmarkBitcaskDBBackend_Delete(b *testing.B) {
	database := NewBitcaskDBDisk()
	benchmarkDelete(b, database)
}

func BenchmarkBitcaskDBBackend_Drop(b *testing.B) {
	database := NewBitcaskDBDisk()
	benchmarkDrop(b, database)
}

func BenchmarkBitcaskDBBackend_Get(b *testing.B) {
	database := NewBitcaskDBDisk()
	benchmarkGet(b, database)
}

func BenchmarkBitcaskDBBackend_Read(b *testing.B) {
	database := NewBitcaskDBDisk()
	benchmarkRead(b, database)
}

func BenchmarkBitcaskDBBackend_Update(b *testing.B) {
	database := NewBitcaskDBDisk()
	benchmarkUpdate(b, database)
}
