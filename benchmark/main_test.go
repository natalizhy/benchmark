package benchmark

import (
	"strconv"
	"sync"
	"testing"
)

type testStruct struct {
	X int
	Y string
}

func (t *testStruct) ToJSON() ([]byte, error) {
	return []byte(`{"X": ` + strconv.Itoa(t.X) + `, "Y": "` + t.Y + `"}`), nil
}

func BenchmarkToJSON(b *testing.B) {
	tmp := &testStruct{X: 1, Y: "string"}
	js, err := tmp.ToJSON()
	if err != nil {
		b.Fatal(err)
	}
	b.SetBytes(int64(len(js)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := tmp.ToJSON(); err != nil {
			b.Fatal(err)
		}
	}
}

type Set struct {
	set map[interface{}]struct{}
	mu  sync.Mutex
}

func (s *Set) Add(x interface{}) {
	s.mu.Lock()
	s.set[x] = struct{}{}
	s.mu.Unlock()
}

func (s *Set) Delete(x interface{}) {
	s.mu.Lock()
	delete(s.set, x)
	s.mu.Unlock()
}

func BenchmarkSetDelete(b *testing.B) {
	var testSet []string
	for i := 0; i < 1024; i++ {
		testSet = append(testSet, strconv.Itoa(i))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		set := Set{set: make(map[interface{}]struct{})}
		for _, elem := range testSet {
			set.Add(elem)
		}
		b.StartTimer()
		for _, elem := range testSet {
			set.Delete(elem)
		}
	}
}

//func BenchmarkSample(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		if x := fmt.Sprintf("%d", 42); x != "42" {
//			b.Fatalf("Unexpected string: %s", x)
//		}
//	}
//}

//test -bench=. <name_test.go>
//go test -benchmem <name_test.go>
