package worker

import "context"

type job struct {
	ctx context.Context
	fn  func() error
}

type worker struct {
	parentPool *Pool

	jobChan chan *job
}

func (w *worker) Run() {
	go func() {
		for job := range w.jobChan {
			err := job.fn()
			if err != nil {
				w.parentPool.loggerCreator(job.ctx).Errorf("job return error: %v", err)
			}
			w.parentPool.workerChan <- w
		}
	}()
}
