package service

import (
	"errors"
	"fmt"

	"forex/internal/config"
	"forex/internal/model"
	"forex/internal/store"
)

type Service struct {
	store     *store.Store
	batchSize int
	usd       string
}

func New(s *store.Store, cfg *config.Config) *Service {
	b := cfg.BatchSize
	if b <= 0 {
		b = 1
	}
	usd := cfg.USDPivot
	if usd == "" {
		usd = "USD"
	}
	return &Service{store: s, batchSize: b, usd: usd}
}

func (svc *Service) QueryRates(pairs []string) ([]*model.Rate, []string) {
	rates := make([]*model.Rate, 0, len(pairs))
	missing := make([]string, 0)
	for _, pair := range pairs {
		r, err := svc.store.GetRate(pair)
		if err != nil {
			missing = append(missing, pair)
			continue
		}
		rates = append(rates, r)
	}
	return rates, missing
}

func (svc *Service) Convert(src, dst string, amount float64) (float64, error) {
	if src == "" || dst == "" || amount < 0 {
		return 0, errors.New("invalid conversion input")
	}
	rates := map[string]*model.Rate{}
	for _, r := range svc.store.ListRates() {
		rates[r.Pair] = r
	}
	rate, ok := model.ChainRate(rates, src, dst)
	if !ok {
		return 0, fmt.Errorf("convert %s->%s: %w", src, dst, store.ErrRateNotFound)
	}
	return model.Convert(rate, amount), nil
}

func (svc *Service) ListCurrencies() []model.Currency {
	currencies := model.SupportedCurrencies()
	model.SortCurrencies(currencies)
	return currencies
}

func (svc *Service) Subscribe(pair string, target float64, direction string) (string, error) {
	if _, err := svc.store.GetRate(pair); err != nil {
		return "", fmt.Errorf("subscribe %s: %w", pair, err)
	}
	sub := &model.Subscription{Pair: pair, Target: target, Direction: direction}
	id, err := svc.store.AddSubscription(sub)
	if err != nil {
		return "", fmt.Errorf("subscribe %s: %w", pair, err)
	}
	return id, nil
}

func (svc *Service) ListSubscriptions() []*model.Subscription {
	return svc.store.ListSubscriptions()
}

func (svc *Service) UpdateRates(batch []*model.Rate) int {
	n := 0
	for _, r := range batch {
		if !model.ValidRate(r) {
			continue
		}
		svc.store.UpsertRate(r)
		n++
	}
	return n
}

func (svc *Service) EvaluateSub(sub *model.Subscription) (bool, error) {
	r, err := svc.store.GetRate(sub.Pair)
	if err != nil {
		return false, fmt.Errorf("evaluate %s: %w", sub.ID, err)
	}
	if sub.Notified {
		return false, nil
	}
	if !model.TriggerMet(r.Value, sub.Target, sub.Direction) {
		return false, nil
	}
	if err := svc.store.MarkNotified(sub.ID); err != nil {
		return false, fmt.Errorf("mark notified %s: %w", sub.ID, err)
	}
	return true, nil
}

func (svc *Service) SubscriptionBatches() [][]*model.Subscription {
	subs := svc.store.ListSubscriptions()
	model.SortSubscriptions(subs)
	return model.BuildSubscriptionBatches(subs, svc.batchSize)
}
