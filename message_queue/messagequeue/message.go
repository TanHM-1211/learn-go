package messagequeue

import (
	"errors"
	"sync"
)

type Message struct {
	Content string `json:"content"`
}

type Node struct {
	data *Message
	next *Node
}

type Queue struct {
	mu   *sync.Mutex
	head *Node
	tail *Node
	cap  int
	len  int
}

var (
	errEmptyQueue = errors.New("can not get from into empty queue")
	errFullQueue  = errors.New("can not insert into full queue")
)

func NewQueue(capacity int) *Queue {
	return &Queue{new(sync.Mutex), nil, nil, capacity, 0}
}

func (q *Queue) Put(node *Node) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.Full() {
		return errFullQueue
	}

	if q.Empty() {
		q.head = node
	} else {
		q.tail.next = node
	}
	q.tail = node
	q.len++
	return nil
}

func (q *Queue) Get() (*Node, error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.Empty() {
		return nil, errEmptyQueue
	}
	res := q.head
	q.head = q.head.next
	q.len--

	if q.Empty() {
		q.tail = nil
	}
	res.next = nil
	return res, nil
}

func (q *Queue) Empty() bool {
	return q.len == 0
}

func (q *Queue) Full() bool {
	return q.len == q.cap
}
