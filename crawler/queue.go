package crawler

import "container/list"

type Queue struct {
	l *list.List
}

func NewQueue() *Queue {
	return &Queue{l: list.New()}
}

func (q *Queue) Push(s string) {
	q.l.PushBack(s)
}

func (q *Queue) Pop() string {
	e := q.l.Front()
	q.l.Remove(e)
	return e.Value.(string)
}

func (q *Queue) IsEmpty() bool {
	return q.l.Len() == 0
}
