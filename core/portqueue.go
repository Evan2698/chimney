package core

import (
	"sync"
)

// Item port type
type Item uint16 //port

// ItemQueue ..
type ItemQueue struct {
	items []Item
	lock  sync.RWMutex
}

// New ..
func (q *ItemQueue) New() *ItemQueue {
	q.items = []Item{}
	return q
}

// Enqueue ..
func (q *ItemQueue) Enqueue(t Item) {
	q.lock.Lock()
	q.items = append(q.items, t)
	q.lock.Unlock()
}

// Dequeue ..
func (q *ItemQueue) Dequeue() Item {
	q.lock.Lock()
	item := q.items[0]
	q.items = q.items[1:len(q.items)]
	q.lock.Unlock()
	return item
}

// Front ..
func (q *ItemQueue) Front() Item {
	q.lock.Lock()
	item := q.items[0]
	q.lock.Unlock()
	return item
}

// IsEmpty ..
func (q *ItemQueue) IsEmpty() bool {
	return len(q.items) == 0
}

// Size ..
func (q *ItemQueue) Size() int {
	return len(q.items)
}
