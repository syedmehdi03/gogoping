package main

import (
	"sync"
)

type Task interface {
	Process() Task
}

type Queue struct {
	tasks       chan Task
	results     chan Task
	concurrency int
	wg          sync.WaitGroup
	producerWg  sync.WaitGroup
}

func NewQueueClient(tasksBuffer, resultsBuffer, concurrency int) *Queue {
	return &Queue{
		tasks:       make(chan Task, tasksBuffer),
		results:     make(chan Task, resultsBuffer),
		concurrency: concurrency,
		wg:          sync.WaitGroup{},
		producerWg:  sync.WaitGroup{},
	}
}

func (q *Queue) EnQueue(task Task) {
	q.tasks <- task
}

func (q *Queue) Run() {
	defer q.closeResults()

	for range q.concurrency {
		q.wg.Go(func() {
			for task := range q.tasks {
				processedTask := task.Process()
				q.pushToResults(processedTask)
			}
		})
	}

	q.wg.Go(func() {
		q.producerWg.Wait()
		close(q.tasks)
	})
}

func (q *Queue) closeResults() {
	wg := sync.WaitGroup{} // using q.wg would result in a deadlock
	wg.Go(func() {
		q.wg.Wait()
		close(q.results)
	})
}

func (q *Queue) pushToResults(result Task) {
	q.results <- result
}

func (q *Queue) GetResults() chan Task {
	return q.results
}

func (q *Queue) AddToProducer() {
	q.producerWg.Add(1)
}

func (q *Queue) ProducerDone() {
	q.producerWg.Done()
}
