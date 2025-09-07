package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("maxErrorsCount <=0", func(t *testing.T) {
		// Добавил кейсы, когда количество горутин подается <=0 (в ТЗ не зафиксировано, что n>0)
		cases := []struct {
			name            string
			goroutinesCount int
			maxErrorsCount  int
			expected        error
		}{
			{"N<0, M<0", -1, -1, ErrErrorsLimitExceeded},
			{"N=0, M=0", 0, 0, ErrErrorsLimitExceeded},
			{"N>0, M=0", 2, 0, ErrErrorsLimitExceeded},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				tasks := make([]Task, 0, 10)
				err := Run(tasks, tc.goroutinesCount, tc.maxErrorsCount)
				require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
			})
		}
	})

	t.Run("no tasks at all", func(t *testing.T) {
		err := Run(nil, 5, 2)
		require.NoError(t, err)
	})

	t.Run("all tasks fail but m > tasksCount", func(t *testing.T) {
		tasks := []Task{
			func() error { return errors.New("fail") },
			func() error { return errors.New("fail") },
		}
		err := Run(tasks, 3, 10)
		require.NoError(t, err)
	})

	t.Run("first task fails with m=1", func(t *testing.T) {
		tasks := []Task{
			func() error { return errors.New("fail") },
			func() error { return nil },
			func() error { return nil },
		}
		err := Run(tasks, 5, 1)
		require.True(t, errors.Is(err, ErrErrorsLimitExceeded))
	})

	// тест дополнен кейсами без изменения логики
	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		cases := []struct {
			name           string
			tasksCount     int
			workersCount   int
			maxErrorsCount int
		}{
			{"Tasks more then workers", 50, 5, 1},
			{"Tasks count is equal workers count ", 10, 10, 1},
			{"Tasks less then workers", 2, 5, 1},
			{"Tasks = workers = 1", 1, 1, 1},
			{"Workers = 1", 10, 1, 1},
		}
		for _, tc := range cases {
			tasksCount := tc.tasksCount
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

			workersCount := tc.workersCount
			maxErrorsCount := tc.maxErrorsCount
			err := Run(tasks, workersCount, maxErrorsCount)

			require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
			require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
		}
	})

	// тест дополнен кейсами без изменения логики (кроме workersCount=1)
	t.Run("tasks without errors", func(t *testing.T) {
		cases := []struct {
			name           string
			tasksCount     int
			workersCount   int
			maxErrorsCount int
		}{
			{"Tasks more then workers", 50, 5, 1},
			{"Tasks count is equal workers count ", 10, 10, 1},
			{"Tasks less then workers", 2, 5, 1},
			{"Tasks = workers = 1", 1, 1, 1},
			{"Workers = 1", 10, 1, 1},
			{"No tasks at running moment", 0, 5, 1},
			{"A lot of tasks - stress test", 100000, 5, 1},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
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

				workersCount := tc.workersCount
				maxErrorsCount := tc.maxErrorsCount

				start := time.Now()
				err := Run(tasks, workersCount, maxErrorsCount)
				elapsedTime := time.Since(start)
				require.NoError(t, err)

				require.Equal(t, int32(tasksCount), runTasksCount, "not all tasks were completed")
				if workersCount != 1 {
					require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
				}
			})
		}
	})
}
