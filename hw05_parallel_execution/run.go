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
	if n <= 0 || len(tasks) == 0 {
		return errors.New("invalid arguments")
	}

	taskChan := make(chan Task)
	errorsChan := make(chan error)
	doneChan := make(chan struct{})
	wg := sync.WaitGroup{}

	for i := 0; i < n; i++ {
		wg.Add(1)
		go worker(taskChan, errorsChan, doneChan, &wg)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, task := range tasks {
			select {
			case <-doneChan:
				return
			case taskChan <- task:
			}
		}
		close(taskChan)
	}()

	var errorsCount int32
	go func() {
		for range errorsChan {
			if m > 0 {
				if atomic.AddInt32(&errorsCount, 1) >= int32(m) {
					close(doneChan)
					return
				}
			}
		}
	}()

	wg.Wait()
	close(errorsChan)

	if m > 0 && atomic.LoadInt32(&errorsCount) >= int32(m) {
		return ErrErrorsLimitExceeded
	}
	return nil
}

func worker(tasksChan <-chan Task, errsChan chan<- error, doneChan <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case task, ok := <-tasksChan:
			if !ok {
				return
			}
			if err := task(); err != nil {
				select {
				case errsChan <- err:
				case <-doneChan:
					return
				}
			}
		case <-doneChan:
			return
		}
	}
}
