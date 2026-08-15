package worker

import (
	"context"
	"testing"

	"forex/internal/config"
	"forex/internal/model"
	"forex/internal/service"
	"forex/internal/store"
)

func newPool() (*service.Service, *Pool) {
	s := store.New()
	svc := service.New(s, config.Load())
	return svc, New(svc, 4)
}

func TestCheckSummary(t *testing.T) {
	svc, p := newPool()
	_ = svc.UpdateRates([]*model.Rate{{Pair: "USD/CNY", Value: 7.2}})
	for i := 0; i < 10; i++ {
		if _, err := svc.Subscribe("USD/CNY", 7.0, model.DirAbove); err != nil {
			t.Fatal(err)
		}
	}
	sum := p.Check(context.Background())
	if sum.Checked != 10 || sum.Triggered != 10 || sum.Failed != 0 {
		t.Fatalf("summary=%+v", sum)
	}
}

func TestCheckNoTrigger(t *testing.T) {
	svc, p := newPool()
	_ = svc.UpdateRates([]*model.Rate{{Pair: "USD/CNY", Value: 6.9}})
	_, _ = svc.Subscribe("USD/CNY", 7.0, model.DirAbove)
	sum := p.Check(context.Background())
	if sum.Triggered != 0 {
		t.Fatalf("summary=%+v", sum)
	}
}

func TestCheckCancel(t *testing.T) {
	svc, p := newPool()
	_ = svc.UpdateRates([]*model.Rate{{Pair: "USD/CNY", Value: 7.2}})
	_, _ = svc.Subscribe("USD/CNY", 7.0, model.DirAbove)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	sum := p.Check(ctx)
	if sum.Triggered != 0 {
		t.Fatalf("triggered=%d want 0", sum.Triggered)
	}
}
