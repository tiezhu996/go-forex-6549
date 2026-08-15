package model

import "testing"

func TestValidRate(t *testing.T) {
	if !ValidRate(&Rate{Pair: "USD/CNY", Value: 7.2}) {
		t.Fatal("valid rate rejected")
	}
	if ValidRate(nil) || ValidRate(&Rate{Pair: "USD/CNY"}) || ValidRate(&Rate{Pair: "", Value: 1}) {
		t.Fatal("invalid rate accepted")
	}
}

func TestValidSubscription(t *testing.T) {
	ok := &Subscription{ID: "sub-1", Pair: "USD/CNY", Target: 7.0, Direction: DirAbove}
	if !ValidSubscription(ok) {
		t.Fatal("valid subscription rejected")
	}
	for _, s := range []*Subscription{nil, {}, {ID: "x", Pair: "USD/CNY", Target: 0, Direction: DirAbove}, {ID: "x", Pair: "USD/CNY", Target: 7, Direction: "sideways"}} {
		if ValidSubscription(s) {
			t.Fatalf("invalid subscription accepted: %+v", s)
		}
	}
}

func TestChainRate(t *testing.T) {
	rates := map[string]*Rate{
		"USD/EUR": {Pair: "USD/EUR", Value: 0.9},
		"USD/CNY": {Pair: "USD/CNY", Value: 7.2},
	}
	if v, ok := ChainRate(rates, "USD", "CNY"); !ok || v != 7.2 {
		t.Fatalf("direct chain rate=%v ok=%v", v, ok)
	}
	if v, ok := ChainRate(rates, "EUR", "CNY"); !ok || v != 8.0 {
		t.Fatalf("pivot chain rate=%v ok=%v", v, ok)
	}
	if v, ok := ChainRate(rates, "EUR", "USD"); !ok || v != 1/0.9 {
		t.Fatalf("invert chain rate=%v ok=%v", v, ok)
	}
	if _, ok := ChainRate(rates, "EUR", "JPY"); ok {
		t.Fatal("unreachable chain rate should fail")
	}
}

func TestTriggerMet(t *testing.T) {
	if !TriggerMet(7.2, 7.0, DirAbove) {
		t.Fatal("above should trigger")
	}
	if !TriggerMet(7.0, 7.0, DirAbove) {
		t.Fatal("above boundary should trigger")
	}
	if TriggerMet(6.9, 7.0, DirAbove) {
		t.Fatal("below target should not trigger above")
	}
	if !TriggerMet(6.9, 7.0, DirBelow) {
		t.Fatal("below should trigger")
	}
	if TriggerMet(7.1, 7.0, DirBelow) {
		t.Fatal("above target should not trigger below")
	}
}

func TestSortRatesDoesNotAliasInput(t *testing.T) {
	in := []*Rate{{Pair: "B"}, {Pair: "A"}}
	out := SortRates(in)
	if out[0].Pair != "A" || out[1].Pair != "B" {
		t.Fatalf("sort order wrong: %v", out)
	}
}

func TestBuildSubscriptionBatchesFresh(t *testing.T) {
	subs := []*Subscription{{ID: "1"}, {ID: "2"}, {ID: "3"}}
	batches := BuildSubscriptionBatches(subs, 2)
	if len(batches) != 2 {
		t.Fatalf("batch count=%d", len(batches))
	}
	batches[0][0] = &Subscription{ID: "x"}
	if subs[0].ID != "1" {
		t.Fatal("mutating batch corrupted input")
	}
}

func TestMergeSummary(t *testing.T) {
	got := MergeSummary(Summary{Checked: 1, Failed: 1}, Summary{Checked: 2, Triggered: 3, Failed: 4})
	if got.Checked != 3 || got.Triggered != 3 || got.Failed != 5 {
		t.Fatalf("merge=%+v", got)
	}
}

func TestSupportedCurrencies(t *testing.T) {
	cs := SupportedCurrencies()
	if len(cs) < 30 {
		t.Fatalf("supported currencies=%d, want >=30", len(cs))
	}
}
