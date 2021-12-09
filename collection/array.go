package collection

import (
	"sort"
	"sync"
)

type Array[V any] struct {
	rw   sync.RWMutex
	list []V
}

func (a *Array[V]) Len() int {
	a.rw.RLock()
	defer a.rw.RUnlock()
	return len(a.list)
}

func (a *Array[V]) Append(v V) {
	a.rw.Lock()
	defer a.rw.Unlock()
	a.list = append(a.list, v)
}

func (a *Array[V]) Get(i int) V {
	a.rw.RLock()
	defer a.rw.RUnlock()
	return a.list[i]
}

func (a *Array[V]) ForEach(f func(i int, v V) bool) {
	a.rw.RLock()
	defer a.rw.RUnlock()
	for i, v := range a.list {
		if !f(i, v) {
			break
		}
	}
}

func (a *Array[V]) Find(f func(i int, v V) bool) (index int, value V) {
	a.rw.RLock()
	defer a.rw.RUnlock()
	index = -1
	for i, v := range a.list {
		if f(i, v) {
			return i, v
		}
	}
	return
}

func (a *Array[V]) All(f func(v V) bool) bool {
	a.rw.RLock()
	defer a.rw.RUnlock()
	for _, v := range a.list {
		if !f(v) {
			return false
		}
	}
	return true
}

func (a *Array[V]) Some(f func(v V) bool) bool {
	a.rw.RLock()
	defer a.rw.RUnlock()
	for _, v := range a.list {
		if f(v) {
			return true
		}
	}
	return false
}

func (a *Array[V]) Sort(c func(v1, v2 V) bool) {
	a.rw.Lock()
	defer a.rw.Unlock()
	sort.SliceStable(a.list, func(i, j int) bool {
		return c(a.list[i], a.list[j])
	})
}

func (a *Array[V]) Slice(start, end int) *Array[V] {
	a.rw.RLock()
	defer a.rw.RUnlock()
	return &Array[V]{
		list: a.list[start:end],
	}
}

func ArrayMap[I, R any](input *Array[I], f func(I) R) *Array[R] {
	input.rw.RLock()
	defer input.rw.RUnlock()

	res := &Array[R]{
		list: make([]R, input.Len()),
	}

	for i, v := range input.list {
		res.list[i] = f(v)
	}

	return res
}

func ArrayReduce[V, R any](a *Array[V], init R, f func(V, R) R) R {
	res := init
	a.ForEach(func(i int, v V) bool {
		res = f(v, res)
		return true
	})
	return res
}
