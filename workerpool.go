package workerpool

import (
	"sync"
	"time"
)

// Worker - runs worker with specified integer id
func (wp *WorkerPool) Worker(id int) {
	for j := range wp.Jobs {
		wp.logger.Debugf("worker %d started job %v", id, j)
		//
		result := wp.Function(j)
		wp.WaitGroup.Done()
		//
		wp.logger.Debugf("worker %d finished job %v", id, j)
		wp.Results <- result
	}
}

// ProcessItems - process multiple times in goroutines
func (wp *WorkerPool) ProcessItems(items []interface{}) (results []interface{}) {
	go func() {
		var needBreak bool
		for {
			select {
			case item := <-wp.Results:
				results = append(results, item)
				wp.logger.Debug("result!")
			case <-wp.AllItemsCh:
				wp.logger.Debug("All items!")
				wp.WaitCh <- struct{}{}
				needBreak = true
				break
			case <-wp.WaitCh:
				wp.logger.Debug("reader got waitch")
				needBreak = true
				break
			}

			if needBreak {
				break
			}
		}
		wp.logger.Debug("reader exit")
	}()

	// All items done, signal exit
	go func() {
		wp.WaitGroup.Wait()
		wp.AllItemsCh <- struct{}{}
		wp.logger.Debug("Wooo!")
	}()

	wp.WaitGroup.Add(len(items))

	for w := 1; w <= wp.WorkerCount; w++ {
		go wp.Worker(w)
	}

	for _, item := range items {
		wp.Jobs <- item
	}

	wp.logger.Debug("close reached")
	close(wp.Jobs)

	select {
	case <-wp.WaitCh:
		break
	case <-time.After(30 * time.Second):
		break
	}

	wp.logger.Debug("return")

	return results
}

// New - create worker pool instance properly
func New(workerCount int, fn func(item interface{}) interface{}, logger Logger) *WorkerPool {
	return &WorkerPool{
		Jobs:        make(chan interface{}),
		Results:     make(chan interface{}),
		WaitCh:      make(chan struct{}),
		AllItemsCh:  make(chan struct{}),
		WaitGroup:   &sync.WaitGroup{},
		WorkerCount: workerCount,
		Function:    fn,
		logger:      logger,
	}
}
