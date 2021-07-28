package pool

import "sync"

type WorkerPool struct {
	capacity int
	workers  chan WorkerFunc
	wg       *sync.WaitGroup
}

func NewWorkerPool(capacity int) *WorkerPool {
	return &WorkerPool{
		capacity: capacity,
		workers:  make(chan WorkerFunc, capacity),
		wg:       &sync.WaitGroup{},
	}
}

func (this *WorkerPool) AddWorker(work WorkerFunc) {
	this.workers <- work

	this.wg.Add(1)

	go func() {
		work()
		<-this.workers
		this.wg.Done()
	}()
}

func (this *WorkerPool) Wait() {
	this.wg.Wait()
	close(this.workers)
}

type WorkerFunc func()
