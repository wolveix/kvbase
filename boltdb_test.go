package kvbase_test

import (
	"github.com/Wolveix/kvbase"
	"os"
	"testing"
)

func TestBoltDBBackend(t *testing.T) {
	_ = os.Remove("testboltdbdata")
	exampleModel := model{
		"John Smith",
	}

	database, err := kvbase.NewBoltDB("testboltdbdata")
	if err != nil {
		t.Fatal("Error on database initialisation:", err)
	}

	if err := database.Create("bucket", "key", &exampleModel); err != nil {
		t.Fatal("Error on record creation:", err)
	}

	if _, err = kvbase.NewBoltDB("testboltdbdata"); err == nil {
		t.Fatal("Expected error on database initialisation:", err)
	}
}

func TestBoltDBBackend_Count(t *testing.T) {
	_ = os.Remove("testboltdbdata")
	exampleModel := model{
		"John Smith",
	}

	database, err := kvbase.NewBoltDB("testboltdbdata")
	if err != nil {
		t.Fatal("Error on database initialisation:", err)
	}

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

func TestBoltDBBackend_Create(t *testing.T) {
	_ = os.Remove("testboltdbdata")
	exampleModel := model{
		"John Smith",
	}

	database, err := kvbase.NewBoltDB("testboltdbdata")
	if err != nil {
		t.Fatal("Error on database initialisation:", err)
	}

	if err := database.Create("bucket", "key", &exampleModel); err != nil {
		t.Fatal("Error on record creation:", err)
	}

	if err = database.Create("bucket", "key", &exampleModel); err == nil {
		t.Fatal("Error expected for existing key.")
	}
}

func TestBoltDBBackend_Delete(t *testing.T) {
	_ = os.Remove("testboltdbdata")
	exampleModel := model{
		"John Smith",
	}

	database, err := kvbase.NewBoltDB("testboltdbdata")
	if err != nil {
		t.Fatal("Error on database initialisation:", err)
	}

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

func TestBoltDBBackend_Get(t *testing.T) {
	_ = os.Remove("testboltdbdata")

	database, err := kvbase.NewBoltDB("testboltdbdata")
	if err != nil {
		t.Fatal("Error on database initialisation:", err)
	}

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

func TestBoltDBBackend_Read(t *testing.T) {
	_ = os.Remove("testboltdbdata")
	exampleModel := model{
		"John Smith",
	}

	database, err := kvbase.NewBoltDB("testboltdbdata")
	if err != nil {
		t.Fatal("Error on database initialisation:", err)
	}

	if err := database.Create("bucket", "key", &exampleModel); err != nil {
		t.Fatal("Error on record creation:", err)
	}

	emptyModel := model{}

	if err = database.Read("bucket", "key", &emptyModel); err != nil {
		t.Fatal("Error on database read:", err)
	}

	if emptyModel.Name != "John Smith" {
		t.Fatal("Expected John Smith for returned struct.Name, got:", emptyModel.Name)
	}
}

func TestBoltDBBackend_Update(t *testing.T) {
	_ = os.Remove("testboltdbdata")
	exampleModel := model{
		"John Smith",
	}

	database, err := kvbase.NewBoltDB("testboltdbdata")
	if err != nil {
		t.Fatal("Error on database initialisation:", err)
	}

	if err := database.Update("bucket", "key", &exampleModel); err == nil {
		t.Fatal("Error expected for missing key.")
	}

	if err := database.Create("bucket", "key", &exampleModel); err != nil {
		t.Fatal("Error on record creation:", err)
	}

	exampleModel.Name = "Updated John Smith"

	if err = database.Update("bucket", "key", &exampleModel); err != nil {
		t.Fatal("Error on record update:", err)
	}

	emptyModel := model{}

	if err = database.Read("bucket", "key", &emptyModel); err != nil {
		t.Fatal("Error on database read:", err)
	}

	if emptyModel.Name != "Updated John Smith" {
		t.Fatal("Expected Updated John Smith for returned struct.Name, got:", emptyModel.Name)
	}
}
