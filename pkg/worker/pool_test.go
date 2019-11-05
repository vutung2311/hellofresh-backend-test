package worker_test

import (
	"bytes"
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"vutung2311-golang-test/pkg/worker"

	"github.com/sirupsen/logrus"
)

func TestPool_AddJob(t *testing.T) {
	var buf bytes.Buffer
	logger := logrus.New()
	logger.SetOutput(&buf)
	logger.SetFormatter(new(logrus.JSONFormatter))
	loggerCreator := func(_ context.Context) logrus.FieldLogger {
		return logger
	}
	ctx := context.Background()

	t.Run("job is fast enough and there is enough worker", func(t *testing.T) {
		counter := int64(0)
		pool := worker.NewPool(100, loggerCreator)
		jobCount := 200
		wg := new(sync.WaitGroup)
		wg.Add(jobCount)
		for i := 0; i < jobCount; i++ {
			err := pool.AddJob(ctx, func() error {
				defer wg.Done()
				atomic.AddInt64(&counter, 1)
				return nil
			})
			if err != nil {
				t.Fatalf("there shouldn't be error. Got %v", err)
			}
		}
		wg.Wait()
		if counter != int64(jobCount) {
			t.Errorf("counter value is wrong. Got %v", counter)
		}
	})
	t.Run("job is fast enough and there is no need to create more worker", func(t *testing.T) {
		counter := int64(0)
		pool := worker.NewPool(100, loggerCreator)
		jobCount := 150
		wg := new(sync.WaitGroup)
		wg.Add(jobCount)
		for i := 0; i < jobCount; i++ {
			err := pool.AddJob(ctx, func() error {
				defer wg.Done()
				atomic.AddInt64(&counter, 1)
				return nil
			})
			if err != nil {
				t.Fatalf("there shouldn't be error. Got %v", err)
			}
		}
		wg.Wait()
		if counter != int64(jobCount) {
			t.Errorf("counter value is wrong. Got %v", counter)
		}
		if pool.ActiveWorkerCount() == 100 {
			t.Errorf("there should be less active worker because job is fast")
		}
	})
	t.Run("job is not fast enough or there is not enough worker", func(t *testing.T) {
		counter := int64(0)
		pool := worker.NewPool(100, loggerCreator)
		jobCount := 200
		for i := 0; i < jobCount; i++ {
			err := pool.AddJob(ctx, func() error {
				atomic.AddInt64(&counter, 1)
				time.Sleep(time.Second)
				return nil
			})
			if i > 99 && err == nil {
				t.Fatal("there should be error")
			}
		}
	})
	t.Run("stop pool should stop all worker", func(t *testing.T) {
		counter := int64(0)
		pool := worker.NewPool(100, loggerCreator)
		jobCount := 100
		for i := 0; i < jobCount; i++ {
			err := pool.AddJob(ctx, func() error {
				atomic.AddInt64(&counter, 1)
				time.Sleep(time.Second)
				return nil
			})
			if err != nil {
				t.Fatalf("there shouldn't be error. Got %v", err)
			}
		}
		if pool.ActiveWorkerCount() == 0 {
			t.Error("active worker count should be positive")
		}
		pool.Stop()
		if pool.ActiveWorkerCount() > 0 {
			t.Error("active worker should be zero")
		}
	})
}
