package store

import (
	"testing"

	"forex/internal/model"
)

func TestPutGetRate(t *testing.T) {
	s := New()
	if err := s.PutRate(&model.Rate{Pair: "USD/CNY", Value: 7.2}); err != nil {
		t.Fatal(err)
	}
	if err := s.PutRate(&model.Rate{Pair: "USD/CNY", Value: 7.3}); err != ErrRateExists {
		t.Fatalf("dup err=%v", err)
	}
	if _, err := s.GetRate("USD/CNY"); err != nil {
		t.Fatal(err)
	}
	if _, err := s.GetRate("EUR/CNY"); err != ErrRateNotFound {
		t.Fatalf("missing err=%v", err)
	}
}

func TestRatePairsFresh(t *testing.T) {
	s := New()
	_ = s.PutRate(&model.Rate{Pair: "USD/CNY", Value: 7.2})
	_ = s.PutRate(&model.Rate{Pair: "USD/EUR", Value: 0.9})
	pairs := s.RatePairs()
	pairs[0] = "XXX/YYY"
	if s.RatePairs()[0] != "USD/CNY" {
		t.Fatal("RatePairs returned aliased slice")
	}
}

func TestListRatesFresh(t *testing.T) {
	s := New()
	_ = s.PutRate(&model.Rate{Pair: "USD/CNY", Value: 7.2})
	rates := s.ListRates()
	rates[0] = &model.Rate{Pair: "XXX/YYY", Value: 1}
	if r, _ := s.GetRate("USD/CNY"); r.Value != 7.2 {
		t.Fatal("ListRates mutated store entry")
	}
}

func TestAddGetSubscription(t *testing.T) {
	s := New()
	id, err := s.AddSubscription(&model.Subscription{Pair: "USD/CNY", Target: 7.0, Direction: model.DirAbove})
	if err != nil {
		t.Fatal(err)
	}
	if id == "" {
		t.Fatal("empty sub id")
	}
	if _, err := s.GetSubscription(id); err != nil {
		t.Fatal(err)
	}
	if _, err := s.GetSubscription("nope"); err != ErrSubscriptionNotFound {
		t.Fatalf("missing err=%v", err)
	}
	if _, err := s.AddSubscription(nil); err == nil {
		t.Fatal("nil subscription accepted")
	}
}

func TestListSubscriptionsFresh(t *testing.T) {
	s := New()
	_, _ = s.AddSubscription(&model.Subscription{Pair: "USD/CNY", Target: 7.0, Direction: model.DirAbove})
	subs := s.ListSubscriptions()
	subs[0] = &model.Subscription{ID: "hacked"}
	if s.ListSubscriptions()[0].ID != "sub-1" {
		t.Fatal("ListSubscriptions returned aliased slice")
	}
}

func TestMarkNotified(t *testing.T) {
	s := New()
	id, _ := s.AddSubscription(&model.Subscription{Pair: "USD/CNY", Target: 7.0, Direction: model.DirAbove})
	if err := s.MarkNotified(id); err != nil {
		t.Fatal(err)
	}
	sub, _ := s.GetSubscription(id)
	if !sub.Triggered || !sub.Notified {
		t.Fatalf("not marked: %+v", sub)
	}
	if err := s.MarkNotified("nope"); err != ErrSubscriptionNotFound {
		t.Fatalf("missing err=%v", err)
	}
}
