package sum

import (
	"runtime"
	"testing"

	"github.com/efficientgo/tools/core/pkg/testutil"
)

func TestSum(t *testing.T) {
	ret, err := Sum("input.txt") // 3.55 MB 1mln lines
	testutil.Ok(t, err)
	testutil.Equals(t, int64(242028430), ret)
}

func TestSum2(t *testing.T) {
	ret, err := Sum2("input.txt") // 3.55 MB 1mln lines
	testutil.Ok(t, err)
	testutil.Equals(t, int64(242028430), ret)
}

func TestConcurrentSum1(t *testing.T) {
	ret, err := ConcurrentSum1("input.txt") // 3.55 MB 1mln lines
	testutil.Ok(t, err)
	testutil.Equals(t, int64(242028430), ret)

	ret, err = ConcurrentSum1("input.txt") // 3.55 MB 1mln lines
	testutil.Ok(t, err)
	testutil.Equals(t, int64(242028430), ret)
}

func TestConcurrentSum2(t *testing.T) {
	ret, err := ConcurrentSum2("input.txt", 4) // 3.55 MB 1mln lines
	testutil.Ok(t, err)
	testutil.Equals(t, int64(242028430), ret)

	ret, err = ConcurrentSum2("input.txt", 11) // 3.55 MB 1mln lines
	testutil.Ok(t, err)
	testutil.Equals(t, int64(242028430), ret)
}

func TestConcurrentSum3(t *testing.T) {
	ret, err := ConcurrentSum3("input.txt", 4) // 3.55 MB 1mln lines
	testutil.Ok(t, err)
	testutil.Equals(t, int64(242028430), ret)

	ret, err = ConcurrentSum3("input.txt", 6) // 3.55 MB 1mln lines
	testutil.Ok(t, err)
	testutil.Equals(t, int64(242028430), ret)
}

func TestConcurrentSumOpt(t *testing.T) {
	ret, err := ConcurrentSumOpt("input.txt", 4) // 3.55 MB 1mln lines
	testutil.Ok(t, err)
	testutil.Equals(t, int64(242028430), ret)

	ret, err = ConcurrentSumOpt("input.txt", 11) // 3.55 MB 1mln lines
	testutil.Ok(t, err)
	testutil.Equals(t, int64(242028430), ret)
}

var Answer int64

// export var=v1 && go test -count 5 -benchtime 5s -run '^$' -bench . -memprofile=${var}.mem.pprof -cpuprofile=${var}.cpu.pprof > ${var}.txt
func BenchmarkSum(b *testing.B) {
	runtime.GOMAXPROCS(4)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Answer, _ = ConcurrentSum3("input.txt", 8)
	}
}
