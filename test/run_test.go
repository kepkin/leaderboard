package ldtest

import (
	"context"
	"flag"
	"os"
	"strings"
	"testing"
	"time"
	"unsafe"

	// "io/fs"
	"fmt"

	"runtime"
	"runtime/debug"

	lb "gihtub.com/kepkin/leaderboard"
	rdstore "gihtub.com/kepkin/leaderboard/test/redis"
	_ "github.com/stretchr/testify/assert"
)

type inputInsertionsType []string

func (i *inputInsertionsType) String() string { return "aha" }

func (i *inputInsertionsType) Set(s string) error {
	*i = append(*i, strings.Split(s, ",")...)
	return nil
}

var (
	ldtestInitTime time.Duration
	inputList      inputInsertionsType
)

func init() {
	flag.DurationVar(&ldtestInitTime, "ldtestInitDuration", 0, "for how long Initialize data. Usually used when changing the ldtest package for debug purposes")
	flag.Var(&inputList, "inputList", "comma separated list of isertions input types: m01s,m,redis")
}

func logMem(b *testing.B, prefix string) {
	var memStat runtime.MemStats
	runtime.ReadMemStats(&memStat)
	b.Logf("%v mem %v", prefix, memStat.HeapAlloc)
}

func runDataTest[K any](b *testing.B, testData *TestData, prefix string, factoryFunc func(*testing.B) K, insertFunc func(sut K, score float64, user string), updateFunc func(sut K, score float64, user string)) {
	variants, err := testData.GetVariants()
	if err != nil {
		b.Fatal(err)
	}

	if len(inputList) != 0 {
		variants = inputList
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

func BenchmarkInsertionRate(b *testing.B) {
	testData, err := NewTestData()
	if err != nil {
		b.Fatal(err)
	}

	b.Run("btree", func(b *testing.B) {
		runDataTest(b, testData, "btree",
			func(b *testing.B) *lb.BtreeStore[float64, string] {
				res := lb.NewBtreeStore[float64, string](
					lb.StdLess[float64],
					lb.StdEquals[float64],
					lb.StdLess[string],
					lb.StdEquals[string],
					101,
				)

				b.Logf("btree sizeof empty: %v", unsafe.Sizeof(*res))

				return res
			},
			func(sut *lb.BtreeStore[float64, string], score float64, user string) {
				it, _ := sut.Insert(lb.Tuple[float64, string]{Key: score, Val: user})
				if it != nil {
					it.Close()
				}
			},
			func(sut *lb.BtreeStore[float64, string], score float64, user string) {
				it, _ := sut.Upsert(lb.Tuple[float64, string]{Key: score, Val: user})
				if it != nil {
					it.Close()
				}
			},
		)
	})

	b.Run("heap", func(b *testing.B) {
		runDataTest(b, testData, "heap",
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
	})

	b.Run("redis", func(b *testing.B) {
		runDataTest(b, testData, "redis",
			func(b *testing.B) *rdstore.RDStore {
				return rdstore.New()
			},
			func(sut *rdstore.RDStore, score float64, user string) {
				sut.Insert(score, user)
			},
			func(sut *rdstore.RDStore, score float64, user string) {
				sut.Insert(score, user)
			},
		)
	})

}

func BenchmarkBtreeSize(b *testing.B) {
	testData, err := NewTestData()
	if err != nil {
		b.Fatal(err)
	}

	for _, btreeSize := range []int{3, 5, 51, 101, 501, 1501} {
		testName := fmt.Sprintf("btree-%v", btreeSize)
		b.Run(testName, func(b *testing.B) {
			runDataTest(b, testData, testName,
				func(b *testing.B) *lb.BtreeStore[float64, string] {
					res := lb.NewBtreeStore[float64, string](
						lb.StdLess[float64],
						lb.StdEquals[float64],
						lb.StdLess[string],
						lb.StdEquals[string],
						btreeSize,
					)

					return res
				},
				func(sut *lb.BtreeStore[float64, string], score float64, user string) {
					it, _ := sut.Insert(lb.Tuple[float64, string]{Key: score, Val: user})
					if it != nil {
						it.Close()
					}
				},
				func(sut *lb.BtreeStore[float64, string], score float64, user string) {
					it, _ := sut.Upsert(lb.Tuple[float64, string]{Key: score, Val: user})
					if it != nil {
						it.Close()
					}
				},
			)
		})
	}
}

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	flag.Parse()
	os.Exit(m.Run())
}
