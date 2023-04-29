package bridge

import (
	"github.com/lemon-mint/go-datastructures/queue"
)

type Bridge[T any] struct {
	queue *queue.RingBuffer[T]
	size  uint64
}

func New[T any](size uint64) *Bridge[T] {
	return &Bridge[T]{
		queue: queue.NewRingBuffer[T](size),
		size:  size,
	}
}

func (b *Bridge[T]) IsFull() bool {
	return b.queue.Len() == b.size
}

func (b *Bridge[T]) Push(data T) error {
	if b.IsFull() {
		return errIsFull
	}
	return b.queue.Put(data)
}

func (b *Bridge[T]) Pop() (T, error) {
	return b.queue.Get()
}
