package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	taskChannel := make(chan Task)
	var errorCounter int64
	wg := sync.WaitGroup{}
	wg.Add(n)

	runTaskExecutors(n, &wg, &taskChannel, &errorCounter)
	pullTasks(tasks, &taskChannel, int64(m), &errorCounter)

	close(taskChannel)

	wg.Wait()

	if m > 0 && errorCounter >= int64(m) {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func runTaskExecutors(n int, wg *sync.WaitGroup, taskChannel *chan Task, errorCounter *int64) {
	for i := 0; i < n; i++ {
		go (func() {
			defer wg.Done()

			for task := range *taskChannel {
				if err := task(); err != nil {
					atomic.AddInt64(errorCounter, 1)
				}
			}
		})()
	}
}

func pullTasks(tasks []Task, taskChannel *chan Task, m int64, errorCounter *int64) {
	for _, task := range tasks {
		if m > 0 && atomic.LoadInt64(errorCounter) >= m {
			break
		}

		*taskChannel <- task
	}
}
