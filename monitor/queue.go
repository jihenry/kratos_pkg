package monitor

import (
	"container/list"
	"sync"
)

// Queue 加锁的队列
type Queue struct {
	l *list.List
	m sync.Mutex
}

// NewQueue 声明新队列
func NewQueue() *Queue {
	return &Queue{l: list.New()}
}

// PushBack 插入新队列
func (q *Queue) PushBack(v interface{}) {
	if v == nil {
		return
	}
	q.m.Lock()
	defer q.m.Unlock()
	q.l.PushBack(v)
}

// Front 取栈顶
func (q *Queue) Front() *list.Element {
	q.m.Lock()
	defer q.m.Unlock()
	return q.l.Front()
}

// Remove 去除元素
func (q *Queue) Remove(e *list.Element) {
	if e == nil {
		return
	}
	q.m.Lock()
	defer q.m.Unlock()
	q.l.Remove(e)
}

// Len 长度
func (q *Queue) Len() int {
	q.m.Lock()
	defer q.m.Unlock()
	return q.l.Len()
}
