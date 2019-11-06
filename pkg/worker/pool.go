package worker

import (
	"context"
	"errors"
	"sync/atomic"

	"github.com/sirupsen/logrus"
)

var ErrAllWorkerAreBusy = errors.New("all worker is busy")

type Pool struct {
	maxWorkerCount    int64
	activeWorkerCount int64
	stopped           bool

	loggerCreator func(context.Context) logrus.FieldLogger

	workerChan chan *worker
}

func NewPool(maxWorkerCount int, loggerCreator func(context.Context) logrus.FieldLogger) *Pool {
	return &Pool{
		maxWorkerCount: int64(maxWorkerCount),
		loggerCreator:  loggerCreator,
		workerChan:     make(chan *worker, maxWorkerCount),
	}
}

func (p *Pool) ActiveWorkerCount() int64 {
	return atomic.LoadInt64(&p.activeWorkerCount)
}

func (p *Pool) AddJob(ctx context.Context, fn func() error) error {
	if p.stopped {
		return errors.New("pool stopped")
	}
	workerJob := &job{
		ctx: ctx,
		fn:  fn,
	}
	select {
	case worker := <-p.workerChan:
		worker.jobChan <- workerJob
	default:
		if atomic.LoadInt64(&p.activeWorkerCount) >= p.maxWorkerCount {
			return ErrAllWorkerAreBusy
		}
		worker := &worker{
			parentPool: p,
			jobChan:    make(chan *job),
		}
		worker.Run()
		worker.jobChan <- workerJob
		atomic.AddInt64(&p.activeWorkerCount, 1)
	}
	return nil
}

func (p *Pool) Stop() {
	p.stopped = true
	for worker := range p.workerChan {
		close(worker.jobChan)
		atomic.AddInt64(&p.activeWorkerCount, -1)
		if atomic.LoadInt64(&p.activeWorkerCount) == 0 {
			return
		}
	}
}
