package main

import "sync"

type Job interface{}

type Result interface{}

type Worker[T Job, R Result] func(workerID int, job T) R

type WorkerPool[T Job, R Result] struct {
	numOfWorkers int
	jobs         chan T
	results      chan R
	done         chan struct{}
	worker       Worker[T, R]
}

func NewWorkerPool[T Job, R Result](numOfWorkers, bufferSize int, worker Worker[T, R]) *WorkerPool[T, R] {
	return &WorkerPool[T, R]{
		numOfWorkers: numOfWorkers,
		jobs:         make(chan T, bufferSize),
		results:      make(chan R, bufferSize),
		done:         make(chan struct{}),
		worker:       worker,
	}
}

func (pool *WorkerPool[T, R]) Start() {
	var wg sync.WaitGroup
	done := make(chan struct{})

	for i := 1; i <= pool.numOfWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for job := range pool.jobs {
				select {
				case <-pool.done:
					return
				default:
					result := pool.worker(workerID, job)
					pool.results <- result
				}
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(done)
		close(pool.results)
	}()
}

func (pool *WorkerPool[T, R]) Submit(job T) {
	pool.jobs <- job

}

func (pool *WorkerPool[T, R]) Results() <-chan R {
	return pool.results
}

func (pool *WorkerPool[T, R]) Stop() {
	close(pool.done)
}

func (pool *WorkerPool[T, R]) CloseQueue() {
	close(pool.jobs)
}
