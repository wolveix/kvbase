package kvbaseBackendBitcask_test

import (
	_ "github.com/Wolveix/kvbase/backend/bitcask"
	"github.com/Wolveix/kvbase/pkg/kvbaseBackendTest"
	"testing"
)

func Test_Disk(t *testing.T) {
	kvbaseBackendTest.RunTests(t, "bitcask", "testdata", false)
}

func Benchmark_Disk(b *testing.B) {
	kvbaseBackendTest.RunBenches(b, "bitcask", "testdata", false)
}
