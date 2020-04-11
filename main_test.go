package kvbase_test

import (
	"fmt"
	"github.com/Wolveix/kvbase"
	"strconv"
	"testing"
)

type model struct {
	Name string
}

var exampleModel = model{
	"John Smith",
}

func testCount(t *testing.T, database kvbase.Backend) {
	if err := database.Create("bucket", "key", &exampleModel); err != nil {
		t.Fatal("Error on record creation:", err)
	}

	counter, err := database.Count("bucket")
	if err != nil {
		t.Fatal("Error on record count:", err)
	}

	if counter != 1 {
		t.Fatal("Expected 1 from counter, got", counter)
	}
}

func testCreate(t *testing.T, database kvbase.Backend) {
	if err := database.Create("bucket", "key", &exampleModel); err != nil {
		t.Fatal("Error on record creation:", err)
	}

	if err := database.Create("bucket", "key", &exampleModel); err == nil {
		t.Fatal("Error expected for existing key.")
	}
}

func testDelete(t *testing.T, database kvbase.Backend) {
	if err := database.Create("bucket", "key", &exampleModel); err != nil {
		t.Fatal("Error on record creation:", err)
	}

	if err := database.Delete("bucket", "key"); err != nil {
		t.Fatal("Error on record deletion:", err)
	}

	if err := database.Delete("bucket", "key"); err == nil {
		t.Fatal("Error expected for missing key.")
	}
}

func testDrop(t *testing.T, database kvbase.Backend) {
	if err := database.Create("bucket", "k0", &exampleModel); err != nil {
		t.Fatal("Error on record creation:", err)
	}

	newModel := exampleModel
	newModel.Name = "Updated John Smith"

	if err := database.Create("bucket", "k1", &newModel); err != nil {
		t.Fatal("Error on record creation:", err)
	}

	if err := database.Drop("bucket"); err != nil {
		t.Fatal("Error on bucket drop:", err)
	}

	if err := database.Read("bucket", "k0", &model{}); err == nil {
		t.Fatal("Error expected for missing key.")
	}
}

func testGet(t *testing.T, database kvbase.Backend) {
	kO := model{
		"John Smith",
	}

	if err := database.Create("bucket", "keyOne", &kO); err != nil {
		t.Fatal("Error on record creation:", err)
	}

	k1 := model{
		"James Green",
	}

	if err := database.Create("bucket", "keyTwo", &k1); err != nil {
		t.Fatal("Error on record creation:", err)
	}

	results, err := database.Get("bucket", model{})
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

func testRead(t *testing.T, database kvbase.Backend) {
	if err := database.Create("bucket", "key", &exampleModel); err != nil {
		t.Fatal("Error on record creation:", err)
	}

	emptyModel := model{}

	if err := database.Read("bucket", "key", &emptyModel); err != nil {
		t.Fatal("Error on database read:", err)
	}

	if emptyModel.Name != "John Smith" {
		t.Fatal("Expected John Smith for returned struct.Name, got:", emptyModel.Name)
	}
}

func testUpdate(t *testing.T, database kvbase.Backend) {
	if err := database.Update("bucket", "key", &exampleModel); err == nil {
		t.Fatal("Error expected for missing key.")
	}

	if err := database.Create("bucket", "key", &exampleModel); err != nil {
		t.Fatal("Error on record creation:", err)
	}

	newModel := exampleModel
	newModel.Name = "Updated John Smith"

	if err := database.Update("bucket", "key", &newModel); err != nil {
		t.Fatal("Error on record update:", err)
	}

	emptyModel := model{}

	if err := database.Read("bucket", "key", &emptyModel); err != nil {
		t.Fatal("Error on database read:", err)
	}

	if emptyModel.Name != "Updated John Smith" {
		t.Fatal("Expected Updated John Smith for returned struct.Name, got:", emptyModel.Name)
	}
}

func benchmarkCount(b *testing.B, database kvbase.Backend) {
	for i := 0; i < b.N; i++ {
		if err := database.Create("bucket", strconv.Itoa(i), &exampleModel); err != nil {
			b.Error("Error on record creation:", err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := database.Count("bucket"); err != nil {
			b.Error("Error on record creation:", err)
		}
	}
}

func benchmarkCreate(b *testing.B, database kvbase.Backend) {
	for i := 0; i < b.N; i++ {
		if err := database.Create("bucket", strconv.Itoa(i), &exampleModel); err != nil {
			b.Error("Error on record creation:", err)
		}
	}
}

func benchmarkDelete(b *testing.B, database kvbase.Backend) {
	for i := 0; i < b.N; i++ {
		if err := database.Create("bucket", strconv.Itoa(i), &exampleModel); err != nil {
			b.Error("Error on record creation:", err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := database.Delete("bucket", strconv.Itoa(i)); err != nil {
			b.Error("Error on record delete:", err)
		}
	}
}

func benchmarkDrop(b *testing.B, database kvbase.Backend) {
	for i := 0; i < b.N; i++ {
		if err := database.Create("bucket"+strconv.Itoa(i), strconv.Itoa(i), &exampleModel); err != nil {
			b.Error("Error on record creation:", err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := database.Drop("bucket" + strconv.Itoa(i)); err != nil {
			b.Error("Error on bucket drop:", err)
		}
	}
}

func benchmarkGet(b *testing.B, database kvbase.Backend) {
	for i := 0; i < b.N; i++ {
		if err := database.Create("bucket", strconv.Itoa(i), &exampleModel); err != nil {
			b.Error("Error on record creation:", err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		newModel := model{}
		if _, err := database.Get("bucket", &newModel); err != nil {
			b.Error("Error on record get:", err)
		}
	}
}

func benchmarkRead(b *testing.B, database kvbase.Backend) {
	for i := 0; i < b.N; i++ {
		fmt.Println(i)
		if err := database.Create("bucket", strconv.Itoa(i), &exampleModel); err != nil {
			b.Error("Error on record creation:", err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		newModel := model{}
		if err := database.Read("bucket", strconv.Itoa(i), &newModel); err != nil {
			b.Error("Error on record update:", err)
		}
	}
}

func benchmarkUpdate(b *testing.B, database kvbase.Backend) {
	for i := 0; i < b.N; i++ {
		if err := database.Create("bucket", strconv.Itoa(i), &exampleModel); err != nil {
			b.Error("Error on record creation:", err)
		}
	}

	b.ResetTimer()
	newModel := exampleModel
	newModel.Name = "Updated John Smith"
	for i := 0; i < b.N; i++ {
		if err := database.Update("bucket", strconv.Itoa(i), &newModel); err != nil {
			b.Error("Error on record update:", err)
		}
	}
}
