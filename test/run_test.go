package ldtest

import (
	"context"
	"flag"
	"os"
	"testing"
	"time"
	// "io/fs"
	"fmt"

	"runtime"
	"runtime/debug"

	"unsafe"

	_ "github.com/stretchr/testify/assert"

	lb "gihtub.com/kepkin/leaderboard"
)

var testFiles []string = []string{"./data/test_10.csv"}

var (
	ldtestInitTime time.Duration
)

func init() {
	flag.DurationVar(&ldtestInitTime, "ldtestInitDuration", 0, "for how long Initialize data. Usually used when changing the ldtest package for debug purposes")
}

func logMem(b *testing.B, prefix string) {
	var memStat runtime.MemStats
	runtime.ReadMemStats(&memStat)
	b.Logf("%v mem %v", prefix, memStat.HeapAlloc)
}

func runDataTest[K any](b *testing.B, prefix string, factoryFunc func(*testing.B) K, insertFunc func(sut K, score float64, user string), updateFunc func(sut K, score float64, user string)) {
	testData, err := NewTestData()
	if err != nil {
		b.Fatal(err)
	}

	variants, err := testData.GetVariants()
	if err != nil {
		b.Fatal(err)
	}

	for i, v := range variants {
		testData.Choose(prefix, i)
		sut := factoryFunc(b)

		ctx := context.Background()
		if ldtestInitTime > 0 {
			var ctxCancel context.CancelFunc
			ctx, ctxCancel = context.WithTimeout(context.Background(), ldtestInitTime)
			defer ctxCancel()
		}

		lastNewUser, err := testData.Initialize(ctx, i+1, func(score float64, user string) {
			insertFunc(sut, score, user)
		})
		if err != nil && err != context.DeadlineExceeded {
			b.Fatal(err)
		}

		b.Log("Initialized done ", lastNewUser)

		testData.Close()
		debug.FreeOSMemory()
		logMem(b, "after close")

		_ = b.Run(fmt.Sprint(prefix, "-", v), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				var req Request
				req, lastNewUser = MakeRequest(lastNewUser)
				// b.Log(req, lastNewUser)
				// req.User = "372"
				// req.Action = ActionUpdate

				if req.Action == ActionInsert {
					insertFunc(sut, req.Value, req.User)
				} else if req.Action == ActionUpdate {
					updateFunc(sut, req.Value, req.User)
				}
			}
		})
	}
}

func BenchmarkStores(b *testing.B) {
	runDataTest(b, "btree",
		func(b *testing.B) *lb.BtreeStore[float64, string] {
			res := lb.NewBtreeStore[float64, string](
				lb.StdLess[float64],
				lb.StdEquals[float64],
				lb.StdLess[string],
				lb.StdEquals[string],
			)

			b.Logf("btree sizeof empty: %v", unsafe.Sizeof(*res))

			return res
		},
		func(sut *lb.BtreeStore[float64, string], score float64, user string) {
			it, _ := sut.Insert(lb.BTreeLeaf[float64, string]{OrderKey: score, Value: user})
			if it != nil {
				it.Close()
			}
		},
		func(sut *lb.BtreeStore[float64, string], score float64, user string) {
			it, _ := sut.Upsert(lb.BTreeLeaf[float64, string]{OrderKey: score, Value: user})
			if it != nil {
				it.Close()
			}
		},
	)

	runDataTest(b, "heap",
		func(b *testing.B) *lb.HeapStore[float64, string] {
			return lb.NewHeapStore[float64, string](func(a, b float64) bool { return a < b })
		},
		func(sut *lb.HeapStore[float64, string], score float64, user string) {
			sut.Insert(score, user)
		},
		func(sut *lb.HeapStore[float64, string], score float64, user string) {
			sut.Update(score, user)
		},
	)

}

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	flag.Parse()
	os.Exit(m.Run())
}
