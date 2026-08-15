package config

import "testing"

func TestLoadDefaults(t *testing.T) {
	t.Setenv("FOREX_WORKERS", "")
	t.Setenv("FOREX_BATCH_SIZE", "")
	c := Load()
	if c.Workers != 2 || c.BatchSize != 2 || c.USDPivot != "USD" {
		t.Fatalf("unexpected config: %+v", c)
	}
}
