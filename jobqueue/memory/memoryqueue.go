package memoryqueue

import (
	log "github.com/sirupsen/logrus"
)

type MemoryQueue struct {
	jobChannel chan int64
}

func NewMemoryQueue(size int) *MemoryQueue {
	ch := make(chan int64, size)
	return &MemoryQueue{
		jobChannel: ch,
	}
}

func (q *MemoryQueue) Enqueue(BlockNum int64) bool {
	select {
	case q.jobChannel <- BlockNum:
		return true
	default:
		log.Error("jobqueue buffer full")
		return false
	}
}

func (q *MemoryQueue) Subscribe() chan int64 {
	return q.jobChannel
}
