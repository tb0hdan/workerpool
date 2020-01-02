package workerpool

import "sync"

// Logger - worker pool logger interface. Can be stdlib log, logrus etc
type Logger interface {
	Printf(fmt string, args ...interface{})
	Debug(args ...interface{})
	Debugf(fmt string, args ...interface{})
}

// WorkerPool - straightforward worker pool implementation with bound methods
type WorkerPool struct {
	Jobs        chan interface{}
	Results     chan interface{}
	WaitGroup   *sync.WaitGroup
	WaitCh      chan struct{}
	AllItemsCh  chan struct{}
	WorkerCount int
	Function    func(interface{}) interface{}
	logger      Logger
}
