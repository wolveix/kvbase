package kvbaseBackendBoltDB_test

import (
	_ "github.com/Wolveix/kvbase/backend/boltdb"
	"github.com/Wolveix/kvbase/pkg/kvbaseBackendTest"
	"testing"
)

func Test_Disk(t *testing.T) {
	kvbaseBackendTest.RunTests(t, "boltdb", "testdata", false)
}

func Benchmark_Disk(b *testing.B) {
	kvbaseBackendTest.RunBenches(b, "boltdb", "testdata", false)
}
