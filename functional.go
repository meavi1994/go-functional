package go_functional

import (
	"cmp"
	"iter"
	"sync"
)

func AnyAs[T any](v any) T {
	return v.(T)
}

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

// Keys extracts the first element of a Seq2 (the keys).
func Keys[K comparable, V any](s iter.Seq2[K, V]) iter.Seq[K] {
	return func(yield func(K) bool) {
		for k, _ := range s {
			if !yield(k) {
				return
			}
		}
	}
}

// Values extracts the second element of a Seq2 (the values).
func Values[K comparable, V any](s iter.Seq2[K, V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range s {
			if !yield(v) {
				return
			}
		}
	}
}

// SyncKeys converts a sync.Map to an iter.Seq of keys
func SyncKeys[K comparable, V any](m *sync.Map) iter.Seq[K] {
	return func(yield func(K) bool) {
		m.Range(func(k, _ any) bool {
			return yield(k.(K))
		})
	}
}

// SyncValues converts a sync.Map to an iter.Seq of values
func SyncValues[K comparable, V any](m *sync.Map) iter.Seq[V] {
	return func(yield func(V) bool) {
		m.Range(func(_, v any) bool {
			return yield(v.(V))
		})
	}
}

// SyncAll converts a sync.Map to an iter.Seq2 of keys and values
func SyncAll[K comparable, V any](m *sync.Map) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		m.Range(func(k, v any) bool {
			return yield(k.(K), v.(V))
		})
	}
}

// SyncGetAllByKeys returns an iter.Seq2[K, V] for the given keys in the sync.Map.
// Skips keys that are not present.
func SyncGetAllByKeys[K comparable, V any](m *sync.Map, keys []K) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, k := range keys {
			if v, ok := m.Load(k); ok {
				if !yield(k, v.(V)) {
					return
				}
			}
		}
	}
}
