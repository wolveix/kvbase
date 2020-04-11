package kvbase_test

import (
	"github.com/Wolveix/kvbase"
	"os"
	"testing"
)

func NewGoCacheDisk() kvbase.Backend {
	_ = os.Remove("testgocachedata")

	database, err := kvbase.NewGoCache("testgocachedata", false)
	if err != nil {
		panic("Error on database initialisation: " + err.Error())
	}

	return database
}

func NewGoCacheMemory() kvbase.Backend {
	_ = os.Remove("testgocachedata")

	database, err := kvbase.NewGoCache("testgocachedata", true)
	if err != nil {
		panic("Error on database initialisation: " + err.Error())
	}

	return database
}

func TestGoCacheBackend(t *testing.T) {
	database := NewGoCacheDisk()

	if err := database.Create("bucket", "key", &exampleModel); err != nil {
		t.Fatal("Error on record creation:", err)
	}

	database, err := kvbase.NewGoCache("testgocachedata", false)
	if err != nil {
		t.Fatal("Error on database initialisation:", err)
	}

	if err := database.Read("bucket", "key", &exampleModel); err != nil {
		t.Fatal("Error on database read:", err)
	}
}

func TestGoCacheBackendMem(t *testing.T) {
	database := NewGoCacheMemory()

	if err := database.Create("bucket", "key", &exampleModel); err != nil {
		t.Fatal("Error on record creation:", err)
	}

	database, err := kvbase.NewGoCache("testgocachedata", true)
	if err != nil {
		t.Fatal("Error on database initialisation:", err)
	}

	if err := database.Read("bucket", "key", &exampleModel); err == nil {
		t.Fatal("Error expected for missing key.")
	}
}

func TestGoCacheBackend_Count(t *testing.T) {
	t.Run("Disk", func(t *testing.T) {
		database := NewGoCacheDisk()
		testCount(t, database)
	})

	t.Run("Memory", func(t *testing.T) {
		database := NewGoCacheMemory()
		testCount(t, database)
	})
}

func TestGoCacheBackend_Create(t *testing.T) {
	t.Run("Disk", func(t *testing.T) {
		database := NewGoCacheDisk()
		testCreate(t, database)
	})

	t.Run("Memory", func(t *testing.T) {
		database := NewGoCacheMemory()
		testCreate(t, database)
	})
}

func TestGoCacheBackend_Delete(t *testing.T) {
	t.Run("Disk", func(t *testing.T) {
		database := NewGoCacheDisk()
		testDelete(t, database)
	})

	t.Run("Memory", func(t *testing.T) {
		database := NewGoCacheMemory()
		testDelete(t, database)
	})
}

func BenchmarkGoCacheBackend_Drop(b *testing.B) {
	b.Run("Disk", func(b *testing.B) {
		database := NewGoCacheDisk()
		benchmarkDrop(b, database)
	})

	b.Run("Memory", func(b *testing.B) {
		database := NewGoCacheMemory()
		benchmarkDrop(b, database)
	})
}

func TestGoCacheBackend_Drop(t *testing.T) {
	t.Run("Disk", func(t *testing.T) {
		database := NewGoCacheDisk()
		testDrop(t, database)
	})

	t.Run("Memory", func(t *testing.T) {
		database := NewGoCacheDisk()
		testDrop(t, database)
	})
}

func TestGoCacheBackend_Get(t *testing.T) {
	t.Run("Disk", func(t *testing.T) {
		database := NewGoCacheDisk()
		testGet(t, database)
	})

	t.Run("Memory", func(t *testing.T) {
		database := NewGoCacheMemory()
		testGet(t, database)
	})
}

func TestGoCacheBackend_Read(t *testing.T) {
	t.Run("Disk", func(t *testing.T) {
		database := NewGoCacheDisk()
		testRead(t, database)
	})

	t.Run("Memory", func(t *testing.T) {
		database := NewGoCacheMemory()
		testRead(t, database)
	})
}

func TestGoCacheBackend_Update(t *testing.T) {
	t.Run("Disk", func(t *testing.T) {
		database := NewGoCacheDisk()
		testUpdate(t, database)
	})

	t.Run("Memory", func(t *testing.T) {
		database := NewGoCacheMemory()
		testUpdate(t, database)
	})
}

func BenchmarkGoCacheBackend_Count(b *testing.B) {
	b.Run("Disk", func(b *testing.B) {
		database := NewGoCacheDisk()
		benchmarkCount(b, database)
	})

	b.Run("Memory", func(b *testing.B) {
		database := NewGoCacheMemory()
		benchmarkCount(b, database)
	})
}

func BenchmarkGoCacheBackend_Create(b *testing.B) {
	b.Run("Disk", func(b *testing.B) {
		database := NewGoCacheDisk()
		benchmarkCreate(b, database)
	})

	b.Run("Memory", func(b *testing.B) {
		database := NewGoCacheMemory()
		benchmarkCreate(b, database)
	})
}

func BenchmarkGoCacheBackend_Delete(b *testing.B) {
	b.Run("Disk", func(b *testing.B) {
		database := NewGoCacheDisk()
		benchmarkDelete(b, database)
	})

	b.Run("Memory", func(b *testing.B) {
		database := NewGoCacheMemory()
		benchmarkDelete(b, database)
	})
}

func BenchmarkGoCacheBackend_Get(b *testing.B) {
	b.Run("Disk", func(b *testing.B) {
		database := NewGoCacheDisk()
		benchmarkGet(b, database)
	})

	b.Run("Memory", func(b *testing.B) {
		database := NewGoCacheMemory()
		benchmarkGet(b, database)
	})
}

func BenchmarkGoCacheBackend_Read(b *testing.B) {
	b.Run("Disk", func(b *testing.B) {
		database := NewGoCacheDisk()
		benchmarkRead(b, database)
	})

	b.Run("Memory", func(b *testing.B) {
		database := NewGoCacheMemory()
		benchmarkRead(b, database)
	})
}

func BenchmarkGoCacheBackend_Update(b *testing.B) {
	b.Run("Disk", func(b *testing.B) {
		database := NewGoCacheDisk()
		benchmarkUpdate(b, database)
	})

	b.Run("Memory", func(b *testing.B) {
		database := NewGoCacheMemory()
		benchmarkUpdate(b, database)
	})
}
