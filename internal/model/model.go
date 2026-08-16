package model

import "sort"

type Currency struct {
	Code   string
	Name   string
	Symbol string
}

type Rate struct {
	Pair      string
	Value     float64
	UpdatedAt int64
}

type Subscription struct {
	ID        string
	Pair      string
	Target    float64
	Direction string
	Triggered bool
	Notified  bool
}

type Notification struct {
	SubID     string
	Pair      string
	Value     float64
	Target    float64
	Direction string
}

type Summary struct {
	Checked   int
	Triggered int
	Failed    int
}

const (
	DirAbove = "above"
	DirBelow = "below"
)

func PairKey(base, quote string) string {
	return base + "/" + quote
}

func ValidRate(r *Rate) bool {
	return r != nil && r.Pair != "" && r.Value > 0
}

func ValidSubscription(s *Subscription) bool {
	return s != nil && s.ID != "" && s.Pair != "" && s.Target > 0 &&
		(s.Direction == DirAbove || s.Direction == DirBelow)
}

func DirectRate(rates map[string]*Rate, base, quote string) (float64, bool) {
	r, ok := rates[PairKey(base, quote)]
	if !ok || r == nil {
		return 0, false
	}
	return r.Value, true
}

func ChainRate(rates map[string]*Rate, base, quote string) (float64, bool) {
	if base == quote {
		return 1, true
	}
	if v, ok := DirectRate(rates, base, quote); ok {
		return v, true
	}
	baseInUSD, ok1 := rateViaUSD(rates, base)
	quoteInUSD, ok2 := rateViaUSD(rates, quote)
	if !ok1 || !ok2 {
		return 0, false
	}
	return baseInUSD / quoteInUSD, true
}

// rateViaUSD 返回 1 单位 code 折合多少 USD。
func rateViaUSD(rates map[string]*Rate, code string) (float64, bool) {
	if code == "USD" {
		return 1, true
	}
	if v, ok := DirectRate(rates, code, "USD"); ok {
		return v, true
	}
	if v, ok := DirectRate(rates, "USD", code); ok {
		return 1 / v, true
	}
	return 0, false
}

func Convert(rate, amount float64) float64 {
	return rate * amount
}

func TriggerMet(value, target float64, direction string) bool {
	switch direction {
	case DirAbove:
		return value >= target
	case DirBelow:
		return value <= target
	default:
		return false
	}
}

func SortCurrencies(cs []Currency) []Currency {
	sort.SliceStable(cs, func(i, j int) bool { return cs[i].Code < cs[j].Code })
	return cs
}

func SortSubscriptions(subs []*Subscription) []*Subscription {
	sort.SliceStable(subs, func(i, j int) bool { return subs[i].ID < subs[j].ID })
	return subs
}

func SortRates(rs []*Rate) []*Rate {
	sort.SliceStable(rs, func(i, j int) bool { return rs[i].Pair < rs[j].Pair })
	return rs
}

func BuildSubscriptionBatches(subs []*Subscription, size int) [][]*Subscription {
	if size <= 0 {
		size = 1
	}
	out := make([][]*Subscription, 0, (len(subs)+size-1)/size)
	for i := 0; i < len(subs); i += size {
		end := i + size
		if end > len(subs) {
			end = len(subs)
		}
		batch := make([]*Subscription, end-i)
		copy(batch, subs[i:end])
		out = append(out, batch)
	}
	return out
}

func MergeSummary(dst, src Summary) Summary {
	dst.Checked += src.Checked
	dst.Triggered += src.Triggered
	dst.Failed += src.Failed
	return dst
}

func SupportedCurrencies() []Currency {
	return []Currency{
		{Code: "USD", Name: "US Dollar", Symbol: "$"},
		{Code: "EUR", Name: "Euro", Symbol: "€"},
		{Code: "GBP", Name: "British Pound", Symbol: "£"},
		{Code: "JPY", Name: "Japanese Yen", Symbol: "¥"},
		{Code: "CNY", Name: "Chinese Yuan", Symbol: "¥"},
		{Code: "HKD", Name: "Hong Kong Dollar", Symbol: "HK$"},
		{Code: "AUD", Name: "Australian Dollar", Symbol: "A$"},
		{Code: "CAD", Name: "Canadian Dollar", Symbol: "C$"},
		{Code: "CHF", Name: "Swiss Franc", Symbol: "CHF"},
		{Code: "SGD", Name: "Singapore Dollar", Symbol: "S$"},
		{Code: "KRW", Name: "South Korean Won", Symbol: "₩"},
		{Code: "INR", Name: "Indian Rupee", Symbol: "₹"},
		{Code: "BRL", Name: "Brazilian Real", Symbol: "R$"},
		{Code: "MXN", Name: "Mexican Peso", Symbol: "Mex$"},
		{Code: "RUB", Name: "Russian Ruble", Symbol: "₽"},
		{Code: "ZAR", Name: "South African Rand", Symbol: "R"},
		{Code: "NZD", Name: "New Zealand Dollar", Symbol: "NZ$"},
		{Code: "SEK", Name: "Swedish Krona", Symbol: "kr"},
		{Code: "NOK", Name: "Norwegian Krone", Symbol: "kr"},
		{Code: "DKK", Name: "Danish Krone", Symbol: "kr"},
		{Code: "THB", Name: "Thai Baht", Symbol: "฿"},
		{Code: "MYR", Name: "Malaysian Ringgit", Symbol: "RM"},
		{Code: "IDR", Name: "Indonesian Rupiah", Symbol: "Rp"},
		{Code: "PHP", Name: "Philippine Peso", Symbol: "₱"},
		{Code: "VND", Name: "Vietnamese Dong", Symbol: "₫"},
		{Code: "AED", Name: "UAE Dirham", Symbol: "د.إ"},
		{Code: "SAR", Name: "Saudi Riyal", Symbol: "﷼"},
		{Code: "TRY", Name: "Turkish Lira", Symbol: "₺"},
		{Code: "PLN", Name: "Polish Zloty", Symbol: "zł"},
		{Code: "CZK", Name: "Czech Koruna", Symbol: "Kč"},
		{Code: "HUF", Name: "Hungarian Forint", Symbol: "Ft"},
		{Code: "ILS", Name: "Israeli New Shekel", Symbol: "₪"},
	}
}
