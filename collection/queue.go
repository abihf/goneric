package collection

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type queueItem[V any] struct {
	v    V
	next *queueItem[V]
}
type Queue[V any] struct {
	first   *queueItem[V]
	last    *queueItem[V]
	len     uint32
	m       sync.RWMutex
	c       sync.Cond
	stopped int32
}

var errStopped = fmt.Errorf("queue has been stopped")

func (q *Queue[V]) Len() uint32 {
	return atomic.LoadUint32(&q.len)
}

func (q *Queue[V]) Stop() {
	if atomic.CompareAndSwapInt32(&q.stopped, 0, 1) {
		q.c.Broadcast()
	}
}

func (q *Queue[V]) Enqueue(v V) error {
	if atomic.LoadInt32(&q.stopped) == 1 {
		return errStopped
	}

	q.m.Lock()
	defer q.m.Unlock()

	item := &queueItem[V]{v: v}
	if q.first == nil {
		q.first = item
	} else if q.last != nil {
		q.last.next = item
	}

	q.last = item
	atomic.AddUint32(&q.len, 1)
	go q.c.Signal()
	return nil
}

func (q *Queue[V]) Dequeue() (v V, ok bool) {
	q.c.L.Lock()
	defer q.c.L.Unlock()

	for {
		if atomic.LoadInt32(&q.stopped) == 1 {
			return
		}

		// TODO: how to handle concurrency on this
		if q.first != nil {
			break
		}
		q.c.Wait()
	}

	q.m.Lock()
	defer q.m.Unlock()

	item := q.first
	q.first = item.next
	return item.v, true
}

func (q *Queue[V]) Consume(f func(V) error) error {
	for {
		v, ok := q.Dequeue()
		if !ok {
			break
		}
		err := f(v)
		if err != nil {
			return err
		}
	}
	return nil
}
