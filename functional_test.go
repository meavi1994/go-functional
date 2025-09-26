package go_functional

import (
	"iter"
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
