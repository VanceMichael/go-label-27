package workerpool

import (
	"context"

	"github.com/VanceMichael/go-base-airbridge/internal/domain"
)

func (p *Pool) SubmitGroup(ctx context.Context, jobs []Job) (Submission, error) {
	result := Submission{Total: len(jobs)}
	if len(jobs) == 0 {
		return result, domain.ErrInvalid
	}
	for _, job := range jobs {
		if job == nil {
			return result, domain.ErrInvalid
		}
		select {
		case p.jobs <- job:
			result.Accepted++
		case <-ctx.Done():
			return result, ctx.Err()
		case <-p.done:
			return result, domain.ErrState
		}
	}
	return result, nil
}
