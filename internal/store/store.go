package store

import (
	"errors"
	"strconv"
	"sync"

	"forex/internal/model"
)

var (
	ErrRateNotFound         = errors.New("rate not found")
	ErrRateExists           = errors.New("rate already exists")
	ErrSubscriptionNotFound = errors.New("subscription not found")
)

type Store struct {
	mu        sync.RWMutex
	rates     map[string]*model.Rate
	subs      map[string]*model.Subscription
	rateOrder []string
	subOrder  []string
	nextSubID int
}

func New() *Store {
	return &Store{
		rates:     make(map[string]*model.Rate),
		rateOrder: []string{},
		subOrder:  []string{},
		nextSubID: 1,
	}
}

func (s *Store) PutRate(r *model.Rate) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.rates[r.Pair]; ok {
		return ErrRateExists
	}
	s.rates[r.Pair] = r
	s.rateOrder = append(s.rateOrder, r.Pair)
	return nil
}

func (s *Store) UpsertRate(r *model.Rate) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.rates[r.Pair]; !ok {
		s.rateOrder = append(s.rateOrder, r.Pair)
	}
	s.rates[r.Pair] = r
}

func (s *Store) GetRate(pair string) (*model.Rate, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.rates[pair]
	if !ok {
		return nil, ErrRateNotFound
	}
	return r, nil
}

func (s *Store) ListRates() []*model.Rate {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*model.Rate, 0, len(s.rateOrder))
	for _, p := range s.rateOrder {
		out = append(out, s.rates[p])
	}
	return out
}

func (s *Store) RatePairs() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]string, len(s.rateOrder))
	copy(out, s.rateOrder)
	return out
}

func (s *Store) AddSubscription(sub *model.Subscription) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if sub == nil || sub.Pair == "" {
		return "", errors.New("invalid subscription")
	}
	sub.ID = "sub-" + strconv.Itoa(s.nextSubID)
	s.nextSubID++
	s.subs[sub.ID] = sub
	s.subOrder = append(s.subOrder, sub.ID)
	return sub.ID, nil
}

func (s *Store) GetSubscription(id string) (*model.Subscription, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sub, ok := s.subs[id]
	if !ok {
		return nil, ErrSubscriptionNotFound
	}
	return sub, nil
}

func (s *Store) ListSubscriptions() []*model.Subscription {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*model.Subscription, 0, len(s.subOrder))
	for _, id := range s.subOrder {
		out = append(out, s.subs[id])
	}
	return out
}

func (s *Store) MarkNotified(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	sub, ok := s.subs[id]
	if !ok {
		return ErrSubscriptionNotFound
	}
	sub.Triggered = true
	sub.Notified = true
	return nil
}
