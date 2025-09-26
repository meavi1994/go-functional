package go_functional

import (
	"iter"
	"sync"
	"testing"
)

func sliceToSeq[T any](s []T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, v := range s {
			if !yield(v) {
				return
			}
		}
	}
}

func seqToSlice[T any](s iter.Seq[T]) []T {
	var out []T
	for v := range s {
		out = append(out, v)
	}
	return out
}

func TestMap(t *testing.T) {
	s := sliceToSeq([]int{1, 2, 3})
	res := seqToSlice(Map(s, func(x int) int { return x * x }))
	expected := []int{1, 4, 9}
	if len(res) != len(expected) {
		t.Fatalf("expected %v, got %v", expected, res)
	}
	for i := range res {
		if res[i] != expected[i] {
			t.Errorf("expected %v, got %v", expected, res)
		}
	}
}

func TestFilter(t *testing.T) {
	s := sliceToSeq([]int{1, 2, 3, 4})
	res := seqToSlice(Filter(s, func(x int) bool { return x%2 == 0 }))
	expected := []int{2, 4}
	if len(res) != len(expected) {
		t.Fatalf("expected %v, got %v", expected, res)
	}
	for i := range res {
		if res[i] != expected[i] {
			t.Errorf("expected %v, got %v", expected, res)
		}
	}
}

func TestDistinct(t *testing.T) {
	s := sliceToSeq([]int{1, 2, 2, 3, 1})
	res := seqToSlice(Distinct(s))
	expected := []int{1, 2, 3}
	if len(res) != len(expected) {
		t.Fatalf("expected %v, got %v", expected, res)
	}
	for i := range res {
		if res[i] != expected[i] {
			t.Errorf("expected %v, got %v", expected, res)
		}
	}
}

func TestTake(t *testing.T) {
	s := sliceToSeq([]int{10, 20, 30, 40})
	res := seqToSlice(Take(s, 2))
	expected := []int{10, 20}
	if len(res) != len(expected) {
		t.Fatalf("expected %v, got %v", expected, res)
	}
	for i := range res {
		if res[i] != expected[i] {
			t.Errorf("expected %v, got %v", expected, res)
		}
	}
}

func TestReduce(t *testing.T) {
	s := sliceToSeq([]int{1, 2, 3, 4})
	sum := Reduce(s, 0, func(acc, x int) int { return acc + x })
	if sum != 10 {
		t.Errorf("expected 10, got %d", sum)
	}
}

func TestSum(t *testing.T) {
	s := sliceToSeq([]float64{1.5, 2.5, 3.0})
	total := Sum(s)
	if total != 7.0 {
		t.Errorf("expected 7.0, got %f", total)
	}
}

func TestGroupBy(t *testing.T) {
	s := sliceToSeq([]string{"apple", "ant", "banana", "bat"})
	groups := GroupBy(s, func(v string) string { return string(v[0]) })
	if len(groups["a"]) != 2 || len(groups["b"]) != 2 {
		t.Errorf("expected 2 each, got %v", groups)
	}
}

func TestAll(t *testing.T) {
	s := sliceToSeq([]int{2, 4, 6})
	if !All(s, func(x int) bool { return x%2 == 0 }) {
		t.Errorf("expected all even")
	}
}

func TestAny(t *testing.T) {
	s := sliceToSeq([]int{1, 3, 5, 6})
	if !Any(s, func(x int) bool { return x%2 == 0 }) {
		t.Errorf("expected at least one even")
	}
}

func TestSyncKeys(t *testing.T) {
	m := &sync.Map{}
	m.Store("a", 1)
	m.Store("b", 2)

	keys := seqToSlice(SyncKeys[string, int](m))

	expected := map[string]bool{"a": true, "b": true}
	if len(keys) != len(expected) {
		t.Fatalf("expected %d keys, got %d", len(expected), len(keys))
	}
	for _, k := range keys {
		if !expected[k] {
			t.Errorf("unexpected key %s", k)
		}
	}
}

func TestSyncValues(t *testing.T) {
	m := &sync.Map{}
	m.Store("a", 1)
	m.Store("b", 2)

	vals := seqToSlice(SyncValues[string, int](m))

	expected := map[int]bool{1: true, 2: true}
	if len(vals) != len(expected) {
		t.Fatalf("expected %d values, got %d", len(expected), len(vals))
	}
	for _, v := range vals {
		if !expected[v] {
			t.Errorf("unexpected value %d", v)
		}
	}
}

func TestSyncAll(t *testing.T) {
	m := &sync.Map{}
	m.Store("a", 1)
	m.Store("b", 2)

	pairs := []struct {
		Key string
		Val int
	}{}

	SyncAll[string, int](m)(func(k string, v int) bool {
		pairs = append(pairs, struct {
			Key string
			Val int
		}{k, v})
		return true
	})

	expected := map[string]int{"a": 1, "b": 2}
	if len(pairs) != len(expected) {
		t.Fatalf("expected %d pairs, got %d", len(expected), len(pairs))
	}
	for _, p := range pairs {
		if expected[p.Key] != p.Val {
			t.Errorf("for key %s, expected value %d, got %d", p.Key, expected[p.Key], p.Val)
		}
	}
}

func TestSyncGetAllByKeys(t *testing.T) {
	m := &sync.Map{}
	m.Store("a", 1)
	m.Store("b", 2)
	m.Store("c", 3)

	keys := []string{"a", "c", "x"} // "x" does not exist

	pairs := []struct {
		Key string
		Val int
	}{}

	SyncGetAllByKeys[string, int](m, keys)(func(k string, v int) bool {
		pairs = append(pairs, struct {
			Key string
			Val int
		}{k, v})
		return true
	})

	expected := map[string]int{"a": 1, "c": 3} // only existing keys
	if len(pairs) != len(expected) {
		t.Fatalf("expected %d pairs, got %d", len(expected), len(pairs))
	}

	for _, p := range pairs {
		if expected[p.Key] != p.Val {
			t.Errorf("for key %s, expected value %d, got %d", p.Key, expected[p.Key], p.Val)
		}
	}
}

func TestGetAllByKeys_Map(t *testing.T) {
	m := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}

	keys := []string{"a", "c", "x"} // "x" does not exist

	collected := []struct {
		Key string
		Val int
	}{}

	GetAllByKeys[string, int](m, keys)(func(k string, v int) bool {
		collected = append(collected, struct {
			Key string
			Val int
		}{k, v})
		return true
	})

	// Expected: only existing keys "a" and "c"
	expected := map[string]int{"a": 1, "c": 3}

	if len(collected) != len(expected) {
		t.Fatalf("expected %d pairs, got %d", len(expected), len(collected))
	}

	for _, p := range collected {
		if expected[p.Key] != p.Val {
			t.Errorf("for key %s, expected value %d, got %d", p.Key, expected[p.Key], p.Val)
		}
	}
}
