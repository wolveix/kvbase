package kvbaseBackendBadgerDB_test

import (
	_ "github.com/Wolveix/kvbase/backend/badgerdb"
	"github.com/Wolveix/kvbase/pkg/kvbaseBackendTest"
	"testing"
)

func Test_Disk(t *testing.T) {
	kvbaseBackendTest.RunTests(t, "badgerdb", "testdata", false)
}

func Test_Memory(t *testing.T) {
	kvbaseBackendTest.RunTests(t, "badgerdb", "testdata", true)
}

func Benchmark_Disk(b *testing.B) {
	kvbaseBackendTest.RunBenches(b, "badgerdb", "testdata", false)
}

func Benchmark_Memory(b *testing.B) {
	kvbaseBackendTest.RunBenches(b, "badgerdb", "testdata", true)
}