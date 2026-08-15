package service

import (
	"errors"
	"testing"

	"forex/internal/config"
	"forex/internal/model"
	"forex/internal/store"
)

func newSvc() (*store.Store, *Service) {
	s := store.New()
	return s, New(s, config.Load())
}

func TestQueryRates(t *testing.T) {
	_, svc := newSvc()
	_ = svc.UpdateRates([]*model.Rate{
		{Pair: "USD/CNY", Value: 7.2},
		{Pair: "USD/EUR", Value: 0.9},
	})
	rates, missing := svc.QueryRates([]string{"USD/CNY", "EUR/JPY"})
	if len(rates) != 1 || len(missing) != 1 || missing[0] != "EUR/JPY" {
		t.Fatalf("rates=%v missing=%v", rates, missing)
	}
}

func TestConvertChain(t *testing.T) {
	_, svc := newSvc()
	_ = svc.UpdateRates([]*model.Rate{
		{Pair: "USD/EUR", Value: 0.9},
		{Pair: "USD/CNY", Value: 7.2},
	})
	v, err := svc.Convert("EUR", "CNY", 100)
	if err != nil {
		t.Fatal(err)
	}
	if v != 800 {
		t.Fatalf("convert=%v want 800", v)
	}
	if _, err := svc.Convert("EUR", "JPY", 100); !errors.Is(err, store.ErrRateNotFound) {
		t.Fatalf("errors.Is=%v err=%v", errors.Is(err, store.ErrRateNotFound), err)
	}
}

func TestSubscribeRejectsInvalid(t *testing.T) {
	_, svc := newSvc()
	_ = svc.UpdateRates([]*model.Rate{{Pair: "USD/CNY", Value: 7.2}})
	if _, err := svc.Subscribe("USD/CNY", 7.0, "sideways"); err == nil {
		t.Fatal("invalid direction accepted")
	}
	if _, err := svc.Subscribe("", 7.0, model.DirAbove); err == nil {
		t.Fatal("empty pair accepted")
	}
	if _, err := svc.Subscribe("EUR/JPY", 7.0, model.DirAbove); !errors.Is(err, store.ErrRateNotFound) {
		t.Fatalf("missing pair err=%v", err)
	}
}

func TestSubscribeEvaluate(t *testing.T) {
	_, svc := newSvc()
	_ = svc.UpdateRates([]*model.Rate{{Pair: "USD/CNY", Value: 6.9}})
	_, err := svc.Subscribe("USD/CNY", 7.0, model.DirAbove)
	if err != nil {
		t.Fatal(err)
	}
	subs := svc.ListSubscriptions()
	if len(subs) != 1 {
		t.Fatalf("subs=%v", subs)
	}
	triggered, err := svc.EvaluateSub(subs[0])
	if err != nil {
		t.Fatal(err)
	}
	if triggered {
		t.Fatal("should not trigger below target")
	}
	_ = svc.UpdateRates([]*model.Rate{{Pair: "USD/CNY", Value: 7.1}})
	triggered, err = svc.EvaluateSub(subs[0])
	if err != nil {
		t.Fatal(err)
	}
	if !triggered {
		t.Fatal("should trigger above target")
	}
	// 已通知的订阅不应重复触发
	triggered, _ = svc.EvaluateSub(subs[0])
	if triggered {
		t.Fatal("notified subscription triggered again")
	}
}

func TestListCurrencies(t *testing.T) {
	_, svc := newSvc()
	cs := svc.ListCurrencies()
	if len(cs) < 30 {
		t.Fatalf("currencies=%d", len(cs))
	}
	for i := 1; i < len(cs); i++ {
		if cs[i-1].Code > cs[i].Code {
			t.Fatalf("currencies not sorted: %s > %s", cs[i-1].Code, cs[i].Code)
		}
	}
}

func TestSubscriptionBatchesSorted(t *testing.T) {
	_, svc := newSvc()
	_ = svc.UpdateRates([]*model.Rate{{Pair: "USD/CNY", Value: 7.2}})
	_, _ = svc.Subscribe("USD/CNY", 7.0, model.DirAbove)
	_, _ = svc.Subscribe("USD/CNY", 7.5, model.DirAbove)
	batches := svc.SubscriptionBatches()
	if len(batches) == 0 || len(batches[0]) == 0 || batches[0][0].ID != "sub-1" {
		t.Fatalf("batches=%v", batches)
	}
}
