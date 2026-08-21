package workerpool

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/VanceMichael/go-base-airbridge/internal/domain"
)

func TestSubmitGroupValidatesEveryJobBeforeQueueAdmission(t *testing.T) {
	pool := New(1, 4)
	var executed atomic.Int32
	job := Job(func(context.Context) error {
		executed.Add(1)
		return nil
	})
	submission, err := pool.SubmitGroup(context.Background(), []Job{job, nil})
	if !errors.Is(err, domain.ErrInvalid) {
		t.Errorf("submit error = %v", err)
	}
	if submission.Accepted != 0 {
		t.Errorf("failed group accepted %d jobs", submission.Accepted)
	}
	ctx, cancel := context.WithCancel(context.Background())
	pool.Start(ctx)
	time.Sleep(50 * time.Millisecond)
	cancel()
	pool.Stop()
	if got := executed.Load(); got != 0 {
		t.Fatalf("failed group executed %d jobs", got)
	}
}
