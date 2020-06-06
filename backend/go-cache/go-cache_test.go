package kvbaseBackendGoCache_test

import (
	_ "github.com/Wolveix/kvbase/backend/go-cache"
	"github.com/Wolveix/kvbase/pkg/kvbaseBackendTest"
	"testing"
)

func Test_Disk(t *testing.T) {
	kvbaseBackendTest.RunTests(t, "go-cache", "testdata", false)
}

func Test_Memory(t *testing.T) {
	kvbaseBackendTest.RunTests(t, "go-cache", "testdata", true)
}

func Benchmark_Disk(b *testing.B) {
	kvbaseBackendTest.RunBenches(b, "go-cache", "testdata", false)
}

func Benchmark_Memory(b *testing.B) {
	kvbaseBackendTest.RunBenches(b, "go-cache", "testdata", true)
}