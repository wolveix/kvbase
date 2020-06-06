package kvbaseBackendBboltDB_test

import (
	_ "github.com/Wolveix/kvbase/backend/bboltdb"
	"github.com/Wolveix/kvbase/pkg/kvbaseBackendTest"
	"testing"
)

func Test_Disk(t *testing.T) {
	kvbaseBackendTest.RunTests(t, "bboltdb", "testdata", false)
}

func Benchmark_Disk(b *testing.B) {
	kvbaseBackendTest.RunBenches(b, "bboltdb", "testdata", false)
}
