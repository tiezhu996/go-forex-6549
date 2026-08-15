package main

import (
	"fmt"

	"forex/internal/config"
	"forex/internal/service"
	"forex/internal/store"
)

func main() {
	cfg := config.Load()
	st := store.New()
	svc := service.New(st, cfg)
	_ = svc
	fmt.Println("forex ready")
}
