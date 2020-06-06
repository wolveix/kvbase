package kvbaseBackendLevelDB_test

import (
	_ "github.com/Wolveix/kvbase/backend/leveldb"
	"github.com/Wolveix/kvbase/pkg/kvbaseBackendTest"
	"testing"
)

func Test_Disk(t *testing.T) {
	kvbaseBackendTest.RunTests(t, "leveldb", "testdata", false)
}

func Benchmark_Disk(b *testing.B) {
	kvbaseBackendTest.RunBenches(b, "leveldb", "testdata", false)
}
