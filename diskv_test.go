package kvbase_test

import (
	"github.com/Wolveix/kvbase"
	"os"
	"testing"
)

func NewDiskvDBDisk() kvbase.Backend {
	_ = os.RemoveAll("testdiskvdbdata/")

	database, err := kvbase.NewDiskvDB("testdiskvdbdata")
	if err != nil {
		panic("Error on database initialisation: " + err.Error())
	}

	return database
}

func TestDiskvDBBackend(t *testing.T) {
	database := NewDiskvDBDisk()

	if err := database.Create("bucket", "key", &exampleModel); err != nil {
		t.Fatal("Error on record creation:", err)
	}

	if _, err := kvbase.NewDiskvDB("testdiskvdbdata"); err != nil {
		t.Fatal("Error on database initialisation:", err)
	}

	if err := database.Read("bucket", "key", &exampleModel); err != nil {
		t.Fatal("Error on database read:", err)
	}
}

func TestDiskvDBBackend_Count(t *testing.T) {
	database := NewDiskvDBDisk()
	testCount(t, database)
}

func TestDiskvDBBackend_Create(t *testing.T) {
	database := NewDiskvDBDisk()
	testCreate(t, database)
}

func TestDiskvDBBackend_Delete(t *testing.T) {
	database := NewDiskvDBDisk()
	testDelete(t, database)
}

func TestDiskvDBBackend_Drop(t *testing.T) {
	database := NewDiskvDBDisk()
	testDrop(t, database)
}

func TestDiskvDBBackend_Get(t *testing.T) {
	database := NewDiskvDBDisk()
	testGet(t, database)
}

func TestDiskvDBBackend_Read(t *testing.T) {
	database := NewDiskvDBDisk()
	testRead(t, database)
}

func TestDiskvDBBackend_Update(t *testing.T) {
	database := NewDiskvDBDisk()
	testUpdate(t, database)
}

func BenchmarkDiskvDBBackend_Count(b *testing.B) {
	database := NewDiskvDBDisk()
	benchmarkCount(b, database)
}

func BenchmarkDiskvDBBackend_Create(b *testing.B) {
	database := NewDiskvDBDisk()
	benchmarkCreate(b, database)
}

func BenchmarkDiskvDBBackend_Delete(b *testing.B) {
	database := NewDiskvDBDisk()
	benchmarkDelete(b, database)
}

func BenchmarkDiskvDBBackend_Drop(b *testing.B) {
	database := NewDiskvDBDisk()
	benchmarkDrop(b, database)
}

func BenchmarkDiskvDBBackend_Get(b *testing.B) {
	database := NewDiskvDBDisk()
	benchmarkGet(b, database)
}

func BenchmarkDiskvDBBackend_Read(b *testing.B) {
	database := NewDiskvDBDisk()
	benchmarkRead(b, database)
}

func BenchmarkDiskvDBBackend_Update(b *testing.B) {
	database := NewDiskvDBDisk()
	benchmarkUpdate(b, database)
}
