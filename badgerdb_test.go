package kvbase_test

import (
	"github.com/Wolveix/kvbase"
	"os"
	"testing"
)

func NewBadgerDBDisk() kvbase.Backend {
	_ = os.RemoveAll("testbadgerdbdata/")

	database, err := kvbase.NewBadgerDB("testbadgerdbdata", false)
	if err != nil {
		panic("Error on database initialisation: " + err.Error())
	}

	return database
}

func NewBadgerDBMemory() kvbase.Backend {
	_ = os.RemoveAll("testbadgerdbdata/")

	database, err := kvbase.NewBadgerDB("testbadgerdbdata", true)
	if err != nil {
		panic("Error on database initialisation: " + err.Error())
	}

	return database
}

func TestBadgerDBBackendDisk(t *testing.T) {
	database := NewBadgerDBDisk()

	if err := database.Create("bucket", "key", &exampleModel); err != nil {
		t.Fatal("Error on record creation:", err)
	}

	database, err := kvbase.NewBadgerDB("testbadgerdbdata", false)
	if err != nil {
		t.Fatal("Error on database initialisation:", err)
	}

	if err := database.Read("bucket", "key", &exampleModel); err != nil {
		t.Fatal("Error on database read:", err)
	}
}

func TestBadgerDBBackendMemory(t *testing.T) {
	database := NewBadgerDBMemory()

	if err := database.Create("bucket", "key", &exampleModel); err != nil {
		t.Fatal("Error on record creation:", err)
	}

	database, err := kvbase.NewBadgerDB("testbadgerdbdata", true)
	if err != nil {
		t.Fatal("Error on database initialisation:", err)
	}

	if err := database.Read("bucket", "key", &exampleModel); err == nil {
		t.Fatal("Error expected for missing key.")
	}
}

func TestBadgerDBBackend_Count(t *testing.T) {
	t.Run("Disk", func(t *testing.T) {
		database := NewBadgerDBDisk()
		testCount(t, database)
	})

	t.Run("Memory", func(t *testing.T) {
		database := NewBadgerDBMemory()
		testCount(t, database)
	})
}

func TestBadgerDBBackend_Create(t *testing.T) {
	t.Run("Disk", func(t *testing.T) {
		database := NewBadgerDBDisk()
		testCreate(t, database)
	})

	t.Run("Memory", func(t *testing.T) {
		database := NewBadgerDBMemory()
		testCreate(t, database)
	})
}

func TestBadgerDBBackend_Delete(t *testing.T) {
	t.Run("Disk", func(t *testing.T) {
		database := NewBadgerDBDisk()
		testDelete(t, database)
	})

	t.Run("Memory", func(t *testing.T) {
		database := NewBadgerDBMemory()
		testDelete(t, database)
	})
}

func TestBadgerDBBackend_Drop(t *testing.T) {
	t.Run("Disk", func(t *testing.T) {
		database := NewBadgerDBDisk()
		testDrop(t, database)
	})

	t.Run("Memory", func(t *testing.T) {
		database := NewBadgerDBMemory()
		testDrop(t, database)
	})
}

func TestBadgerDBBackend_Get(t *testing.T) {
	t.Run("Disk", func(t *testing.T) {
		database := NewBadgerDBDisk()
		testGet(t, database)
	})

	t.Run("Memory", func(t *testing.T) {
		database := NewBadgerDBMemory()
		testGet(t, database)
	})
}

func TestBadgerDBBackend_Read(t *testing.T) {
	t.Run("Disk", func(t *testing.T) {
		database := NewBadgerDBDisk()
		testRead(t, database)
	})

	t.Run("Memory", func(t *testing.T) {
		database := NewBadgerDBMemory()
		testRead(t, database)
	})
}

func TestBadgerDBBackend_Update(t *testing.T) {
	t.Run("Disk", func(t *testing.T) {
		database := NewBadgerDBDisk()
		testUpdate(t, database)
	})

	t.Run("Memory", func(t *testing.T) {
		database := NewBadgerDBMemory()
		testUpdate(t, database)
	})
}

func BenchmarkBadgerDBBackend_Count(b *testing.B) {
	b.Run("Disk", func(b *testing.B) {
		database := NewBadgerDBDisk()
		benchmarkCount(b, database)
	})

	b.Run("Memory", func(b *testing.B) {
		database := NewBadgerDBMemory()
		benchmarkCount(b, database)
	})
}

func BenchmarkBadgerDBBackend_Create(b *testing.B) {
	b.Run("Disk", func(b *testing.B) {
		database := NewBadgerDBDisk()
		benchmarkCreate(b, database)
	})

	b.Run("Memory", func(b *testing.B) {
		database := NewBadgerDBMemory()
		benchmarkCreate(b, database)
	})
}

func BenchmarkBadgerDBBackend_Delete(b *testing.B) {
	b.Run("Disk", func(b *testing.B) {
		database := NewBadgerDBDisk()
		benchmarkDelete(b, database)
	})

	b.Run("Memory", func(b *testing.B) {
		database := NewBadgerDBMemory()
		benchmarkDelete(b, database)
	})
}

func BenchmarkBadgerDBBackend_Drop(b *testing.B) {
	b.Run("Disk", func(b *testing.B) {
		database := NewBadgerDBDisk()
		benchmarkDrop(b, database)
	})

	b.Run("Memory", func(b *testing.B) {
		database := NewBadgerDBMemory()
		benchmarkDrop(b, database)
	})
}

func BenchmarkBadgerDBBackend_Get(b *testing.B) {
	b.Run("Disk", func(b *testing.B) {
		database := NewBadgerDBDisk()
		benchmarkGet(b, database)
	})

	b.Run("Memory", func(b *testing.B) {
		database := NewBadgerDBMemory()
		benchmarkGet(b, database)
	})
}

func BenchmarkBadgerDBBackend_Read(b *testing.B) {
	b.Run("Disk", func(b *testing.B) {
		database := NewBadgerDBDisk()
		benchmarkRead(b, database)
	})

	b.Run("Memory", func(b *testing.B) {
		database := NewBadgerDBMemory()
		benchmarkRead(b, database)
	})
}

func BenchmarkBadgerDBBackend_Update(b *testing.B) {
	b.Run("Disk", func(b *testing.B) {
		database := NewBadgerDBDisk()
		benchmarkUpdate(b, database)
	})

	b.Run("Memory", func(b *testing.B) {
		database := NewBadgerDBMemory()
		benchmarkUpdate(b, database)
	})
}
