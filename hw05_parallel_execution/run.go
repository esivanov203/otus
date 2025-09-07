package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func Run(tasks []Task, n, m int) error {
	if n <= 0 {
		n = 1
	}
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	var errCount int64
	limit := int64(m) // используем int64 для безопасного сравнения

	taskCh := make(chan Task)
	done := make(chan struct{})
	var wg sync.WaitGroup
	var once sync.Once

	worker := func() {
		defer wg.Done()
		for {
			select {
			case task, ok := <-taskCh:
				if !ok {
					return
				}
				if task() != nil {
					if atomic.AddInt64(&errCount, 1) >= limit {
						once.Do(func() { close(done) })
						return
					}
				}
			case <-done:
				return
			}
		}
	}

	// запуск воркеров
	for i := 0; i < n; i++ {
		wg.Add(1)
		go worker()
	}

	// рассылка задач
sendTasks:
	for _, t := range tasks {
		select {
		case <-done:
			break sendTasks
		case taskCh <- t:
		}
	}

	close(taskCh)
	wg.Wait()

	if atomic.LoadInt64(&errCount) >= limit {
		return ErrErrorsLimitExceeded
	}

	return nil
}
