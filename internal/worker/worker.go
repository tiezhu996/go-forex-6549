package worker

import (
	"context"
	"sync"

	"forex/internal/model"
	"forex/internal/service"
)

type Pool struct {
	svc     *service.Service
	workers int
}

func New(svc *service.Service, workers int) *Pool {
	if workers <= 0 {
		workers = 1
	}
	return &Pool{svc: svc, workers: workers}
}

func (p *Pool) Check(ctx context.Context) model.Summary {
	batches := p.svc.SubscriptionBatches()

	var wg sync.WaitGroup
	ch := make(chan []*model.Subscription, len(batches))

	go func() {
		defer close(ch)
		for _, b := range batches {
			if ctx.Err() != nil {
				close(ch)
				return
			}
			select {
			case <-ctx.Done():
				return
			case ch <- b:
			}
		}
	}()

	var mu sync.Mutex
	var sum model.Summary

	for i := 0; i < p.workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for batch := range ch {
				var local model.Summary
				for _, sub := range batch {
					select {
					case <-ctx.Done():
						return
					default:
					}
					triggered, err := p.svc.EvaluateSub(sub)
					local.Checked++
					if err != nil {
						local.Failed++
						continue
					}
					if triggered {
						local.Triggered++
					}
				}
				mu.Lock()
				sum = model.MergeSummary(sum, local)
				mu.Unlock()
			}
		}()
	}

	wg.Wait()
	return sum
}
