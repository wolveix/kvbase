package kvbaseBackendDiskv_test

import (
	_ "github.com/Wolveix/kvbase/backend/diskv"
	"github.com/Wolveix/kvbase/pkg/kvbaseBackendTest"
	"testing"
)

func Test_Disk(t *testing.T) {
	kvbaseBackendTest.RunTests(t, "diskv", "testdata", false)
}

func Benchmark_Disk(b *testing.B) {
	kvbaseBackendTest.RunBenches(b, "diskv", "testdata", false)
}