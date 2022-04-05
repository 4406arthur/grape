package jobqueue

type JobQueue interface {
	Enqueue(BlockNum int64) bool
	Subscribe() chan int64
}
