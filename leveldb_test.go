package kvbase_test

import (
	"github.com/Wolveix/kvbase"
	"os"
	"testing"
)

func NewLevelDBDisk() kvbase.Backend {
	_ = os.RemoveAll("testleveldbdata/")

	database, err := kvbase.NewLevelDB("testleveldbdata")
	if err != nil {
		panic("Error on database initialisation: " + err.Error())
	}

	return database
}

func TestLevelDBBackend(t *testing.T) {
	database := NewLevelDBDisk()

	if err := database.Create("bucket", "key", &exampleModel); err != nil {
		t.Fatal("Error on record creation:", err)
	}

	if _, err := kvbase.NewLevelDB("testleveldbdata"); err == nil {
		t.Fatal("Expected error on database initialisation:", err)
	}
}

func TestLevelDBBackend_Count(t *testing.T) {
	database := NewLevelDBDisk()
	testCount(t, database)
}

func TestLevelDBBackend_Create(t *testing.T) {
	database := NewLevelDBDisk()
	testCreate(t, database)
}

func TestLevelDBBackend_Delete(t *testing.T) {
	database := NewLevelDBDisk()
	testDelete(t, database)
}

func TestLevelDBBackend_Drop(t *testing.T) {
	database := NewLevelDBDisk()
	testDrop(t, database)
}

func TestLevelDBBackend_Get(t *testing.T) {
	database := NewLevelDBDisk()
	testGet(t, database)
}

func TestLevelDBBackend_Read(t *testing.T) {
	database := NewLevelDBDisk()
	testRead(t, database)
}

func TestLevelDBBackend_Update(t *testing.T) {
	database := NewLevelDBDisk()
	testUpdate(t, database)
}

func BenchmarkLevelDBBackend_Count(b *testing.B) {
	database := NewLevelDBDisk()
	benchmarkCount(b, database)
}

func BenchmarkLevelDBBackend_Create(b *testing.B) {
	database := NewLevelDBDisk()
	benchmarkCreate(b, database)
}

func BenchmarkLevelDBBackend_Delete(b *testing.B) {
	database := NewLevelDBDisk()
	benchmarkDelete(b, database)
}

func BenchmarkLevelDBBackend_Drop(b *testing.B) {
	database := NewLevelDBDisk()
	benchmarkDrop(b, database)
}

func BenchmarkLevelDBBackend_Get(b *testing.B) {
	database := NewLevelDBDisk()
	benchmarkGet(b, database)
}

func BenchmarkLevelDBBackend_Read(b *testing.B) {
	database := NewLevelDBDisk()
	benchmarkRead(b, database)
}

func BenchmarkLevelDBBackend_Update(b *testing.B) {
	database := NewLevelDBDisk()
	benchmarkUpdate(b, database)
}
