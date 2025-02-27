package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if m == 0 {
		return ErrErrorsLimitExceeded
	}
	var mu sync.Mutex
	var wg sync.WaitGroup
	var stopOnce sync.Once

	var errorsLimit int

	chTasks := make(chan Task)
	stop := make(chan struct{})

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(chTasks)
		for _, task := range tasks {
			select {
			case <-stop:
				return
			case chTasks <- task:
			}
		}
	}()

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range chTasks {
				err := task()
				if err != nil {
					mu.Lock()
					errorsLimit++
					if errorsLimit >= m {
						stopOnce.Do(
							func() {
								close(stop)
							},
						)
					}
					mu.Unlock()
				}
				select {
				case <-stop:
					return
				default:
				}
			}
		}()
	}
	wg.Wait()
	if errorsLimit >= m {
		return ErrErrorsLimitExceeded
	}
	return nil
}
