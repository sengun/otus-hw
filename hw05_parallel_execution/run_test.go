package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})
}

func TestRunTaskExecutors(t *testing.T) {
	defer goleak.VerifyNone(t)
	taskChannel := make(chan Task, 3)
	tasks := []Task{
		func() error {
			return nil
		},
		func() error {
			return nil
		},
		func() error {
			return errors.New("some error")
		},
	}
	taskChannel <- tasks[0]
	taskChannel <- tasks[1]
	taskChannel <- tasks[2]
	close(taskChannel)

	wg := sync.WaitGroup{}
	wg.Add(3)
	var errorCounter int64

	runTaskExecutors(3, &wg, &taskChannel, &errorCounter)

	require.Eventually(t, func() bool {
		return errorCounter == int64(1) && len(taskChannel) == 0
	}, time.Second*3, time.Second)
}

func TestPullTasks(t *testing.T) {
	defer goleak.VerifyNone(t)
	tasks := []Task{
		func() error {
			return nil
		},
		func() error {
			return nil
		},
		func() error {
			return errors.New("some error")
		},
	}
	dataSet := []struct {
		taskChannel               chan Task
		m                         int64
		errorCounter              int64
		expectedTaskChannelLength int
	}{
		{make(chan Task, 3), int64(1), int64(1), 0},
		{make(chan Task, 3), int64(2), int64(1), 3},
		{make(chan Task, 3), int64(0), int64(1), 3},
	}

	for _, ds := range dataSet {
		ds := ds
		pullTasks(tasks, &ds.taskChannel, ds.m, &ds.errorCounter)

		require.Eventually(t, func() bool {
			return len(ds.taskChannel) == ds.expectedTaskChannelLength
		}, time.Second*3, time.Second)
	}
}
