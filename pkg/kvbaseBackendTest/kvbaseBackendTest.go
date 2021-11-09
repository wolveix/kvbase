package kvbaseBackendTest

import (
	"os"
	"strconv"
	"testing"

	"github.com/Wolveix/kvbase"
	_ "github.com/Wolveix/kvbase/backend/badgerdb"
)

type model struct {
	Name string
}

var (
	exampleModel = model{
		"John Smith",
	}
	err   error
	store kvbase.Backend
)

func reset(backend, source string, memory bool) {
	if !memory {
		if err = os.RemoveAll(source); err != nil {
			panic(err)
		}
	}

	if store, err = kvbase.New(backend, source, memory); err != nil {
		panic(err)
	}
}

func RunTests(t *testing.T, backend, source string, memory bool) {
	t.Run(backend+"_Count", func(t *testing.T) {
		reset(backend, source, memory)
		testCount(t)
	})

	t.Run(backend+"_Create", func(t *testing.T) {
		reset(backend, source, memory)
		testCreate(t)
	})

	t.Run(backend+"_Delete", func(t *testing.T) {
		reset(backend, source, memory)
		testDelete(t)
	})

	t.Run(backend+"_Drop", func(t *testing.T) {
		reset(backend, source, memory)
		testDrop(t)
	})

	t.Run(backend+"_Get", func(t *testing.T) {
		reset(backend, source, memory)
		testGet(t)
	})

	t.Run(backend+"_Read", func(t *testing.T) {
		reset(backend, source, memory)
		testRead(t)
	})

	t.Run(backend+"_Update", func(t *testing.T) {
		reset(backend, source, memory)
		testUpdate(t)
	})

	if !memory {
		if err = os.RemoveAll(source); err != nil {
			panic(err)
		}
	}
}

func RunBenches(b *testing.B, backend, source string, memory bool) {
	b.Run(backend+"_Count", func(b *testing.B) {
		reset(backend, source, memory)
		benchmarkCount(b)
	})

	b.Run(backend+"_Create", func(b *testing.B) {
		reset(backend, source, memory)
		benchmarkCreate(b)
	})

	b.Run(backend+"_Delete", func(b *testing.B) {
		reset(backend, source, memory)
		benchmarkDelete(b)
	})

	b.Run(backend+"_Drop", func(b *testing.B) {
		reset(backend, source, memory)
		benchmarkDrop(b)
	})

	b.Run(backend+"_Get", func(b *testing.B) {
		reset(backend, source, memory)
		benchmarkGet(b)
	})

	b.Run(backend+"_Read", func(b *testing.B) {
		reset(backend, source, memory)
		benchmarkRead(b)
	})

	b.Run(backend+"_Update", func(b *testing.B) {
		reset(backend, source, memory)
		benchmarkUpdate(b)
	})

	if !memory {
		if err = os.RemoveAll(source); err != nil {
			panic(err)
		}
	}
}

func testCount(t *testing.T) {
	if err := store.Create("bucket", "key", &exampleModel); err != nil {
		t.Fatal("Error on record creation:", err)
	}

	counter, err := store.Count("bucket")
	if err != nil {
		t.Fatal("Error on record count:", err)
	}

	if counter != 1 {
		t.Fatal("Expected 1 from counter, got", counter)
	}
}

func testCreate(t *testing.T) {
	if err := store.Create("bucket", "key", &exampleModel); err != nil {
		t.Fatal("Error on record creation:", err)
	}

	if err := store.Create("bucket", "key", &exampleModel); err == nil {
		t.Fatal("Error expected for existing key.")
	}
}

func testDelete(t *testing.T) {
	if err := store.Create("bucket", "key", &exampleModel); err != nil {
		t.Fatal("Error on record creation:", err)
	}

	if err := store.Delete("bucket", "key"); err != nil {
		t.Fatal("Error on record deletion:", err)
	}

	if err := store.Delete("bucket", "key"); err == nil {
		t.Fatal("Error expected for missing key.")
	}
}

func testDrop(t *testing.T) {
	newModel := exampleModel
	newModel.Name = "Updated John Smith"

	if err := store.Create("bucket", "k0", &exampleModel); err != nil {
		t.Fatal("Error on record creation:", err)
	}

	if err := store.Create("bucket", "k1", &newModel); err != nil {
		t.Fatal("Error on record creation:", err)
	}

	if err := store.Drop("bucket"); err != nil {
		t.Fatal("Error on bucket drop:", err)
	}

	if err := store.Read("bucket", "k0", &model{}); err == nil {
		t.Fatal("Error expected for missing key.")
	}

	if err := store.Read("bucket", "k1", &model{}); err == nil {
		t.Fatal("Error expected for missing key.")
	}
}

func testGet(t *testing.T) {
	kO := model{
		"John Smith",
	}

	k1 := model{
		"James Green",
	}

	if err := store.Create("bucket", "keyOne", &kO); err != nil {
		t.Fatal("Error on record creation:", err)
	}

	if err := store.Create("bucket", "keyTwo", &k1); err != nil {
		t.Fatal("Error on record creation:", err)
	}

	results, err := store.Get("bucket", model{})
	if err != nil {
		t.Fatal("Error on record get:", err)
	}

	for key, value := range *results {
		if key != "keyOne" && key != "keyTwo" {
			t.Fatal("Expected keyOne or keyTwo, got:", key)
		}

		if value == nil {
			t.Fatal("Expected non-nil value, got:", value)
		}
	}
}

func testRead(t *testing.T) {
	emptyModel := model{}

	if err := store.Create("bucket", "key", &exampleModel); err != nil {
		t.Fatal("Error on record creation:", err)
	}

	if err := store.Read("bucket", "key", &emptyModel); err != nil {
		t.Fatal("Error on store read:", err)
	}

	if emptyModel.Name != "John Smith" {
		t.Fatal("Expected John Smith for returned struct.Name, got:", emptyModel.Name)
	}
}

func testUpdate(t *testing.T) {
	emptyModel := model{}
	newModel := exampleModel
	newModel.Name = "Updated John Smith"

	if err := store.Update("bucket", "key", &exampleModel); err == nil {
		t.Fatal("Error expected for missing key.")
	}

	if err := store.Create("bucket", "key", &exampleModel); err != nil {
		t.Fatal("Error on record creation:", err)
	}

	if err := store.Update("bucket", "key", &newModel); err != nil {
		t.Fatal("Error on record update:", err)
	}

	if err := store.Read("bucket", "key", &emptyModel); err != nil {
		t.Fatal("Error on store read:", err)
	}

	if emptyModel.Name != "Updated John Smith" {
		t.Fatal("Expected Updated John Smith for returned struct.Name, got:", emptyModel.Name)
	}
}

func benchmarkCount(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if err := store.Create("bucket", strconv.Itoa(i), &exampleModel); err != nil {
			b.Error("Error on record creation:", err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := store.Count("bucket"); err != nil {
			b.Error("Error on record creation:", err)
		}
	}
}

func benchmarkCreate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if err := store.Create("bucket", strconv.Itoa(i), &exampleModel); err != nil {
			b.Error("Error on record creation:", err)
		}
	}
}

func benchmarkDelete(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if err := store.Create("bucket", strconv.Itoa(i), &exampleModel); err != nil {
			b.Error("Error on record creation:", err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := store.Delete("bucket", strconv.Itoa(i)); err != nil {
			b.Error("Error on record delete:", err)
		}
	}
}

func benchmarkDrop(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if err := store.Create("bucket"+strconv.Itoa(i), strconv.Itoa(i), &exampleModel); err != nil {
			b.Error("Error on record creation:", err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := store.Drop("bucket" + strconv.Itoa(i)); err != nil {
			b.Error("Error on bucket drop:", err)
		}
	}
}

func benchmarkGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if err := store.Create("bucket", strconv.Itoa(i), &exampleModel); err != nil {
			b.Error("Error on record creation:", err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		newModel := model{}
		if _, err := store.Get("bucket", &newModel); err != nil {
			b.Error("Error on record get:", err)
		}
	}
}

func benchmarkRead(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if err := store.Create("bucket", strconv.Itoa(i), &exampleModel); err != nil {
			b.Error("Error on record creation:", err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		newModel := model{}
		if err := store.Read("bucket", strconv.Itoa(i), &newModel); err != nil {
			b.Error("Error on record update:", err)
		}
	}
}

func benchmarkUpdate(b *testing.B) {
	newModel := exampleModel
	newModel.Name = "Updated John Smith"
	for i := 0; i < b.N; i++ {
		if err := store.Create("bucket", strconv.Itoa(i), &exampleModel); err != nil {
			b.Error("Error on record creation:", err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := store.Update("bucket", strconv.Itoa(i), &newModel); err != nil {
			b.Error("Error on record update:", err)
		}
	}
}
