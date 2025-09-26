package go_functional

import (
	"cmp"
	"iter"
)

func Map[T, R any](s iter.Seq[T], f func(T) R) iter.Seq[R] {
	return func(yield func(R) bool) {
		for v := range s {
			if !yield(f(v)) {
				return
			}
		}
	}
}

func Filter[T any](s iter.Seq[T], pred func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range s {
			if pred(v) && !yield(v) {
				return
			}
		}
	}
}

func DistinctFunc[T any, K comparable](s iter.Seq[T], key func(T) K) iter.Seq[T] {
	return func(yield func(T) bool) {
		seen := make(map[K]struct{})
		for v := range s {
			k := key(v)
			if _, ok := seen[k]; ok {
				continue
			}
			seen[k] = struct{}{}
			if !yield(v) {
				return
			}
		}
	}
}

func Distinct[T comparable](s iter.Seq[T]) iter.Seq[T] {
	return DistinctFunc(s, func(a T) T {
		return a
	})
}

func Take[T any](s iter.Seq[T], n int) iter.Seq[T] {
	return func(yield func(T) bool) {
		count := 0
		for v := range s {
			if count >= n || !yield(v) {
				return
			}
			count++
		}
	}
}

func Reduce[T, R any](s iter.Seq[T], init R, f func(R, T) R) R {
	acc := init
	for v := range s {
		acc = f(acc, v)
	}
	return acc

}

func Sum[T cmp.Ordered](s iter.Seq[T]) T {
	var acc T
	for v := range s {
		acc += v
	}
	return acc
}

func GroupBy[T any, K comparable](s iter.Seq[T], key func(T) K) map[K][]T {
	groups := make(map[K][]T)
	for v := range s {
		k := key(v)
		groups[k] = append(groups[k], v)
	}
	return groups
}

func All[T any](s iter.Seq[T], pred func(T) bool) bool {
	for v := range s {
		if !pred(v) {
			return false
		}
	}
	return true
}

func Any[T any](s iter.Seq[T], pred func(T) bool) bool {
	for v := range s {
		if pred(v) {
			return true
		}
	}
	return false
}
