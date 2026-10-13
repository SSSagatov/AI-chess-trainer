package main

import (
	"context"
	"log"

	"github.com/SSSagatov/AI-chess-trainer/internal/config"
	"github.com/SSSagatov/AI-chess-trainer/internal/db/postgres"
	"github.com/SSSagatov/AI-chess-trainer/migrations"
)

func main() {
	ctx := context.Background()
	cfg, err := config.NewDBConfig()
	if err != nil {
		log.Fatal(err)
	}

	pool, err := postgres.Connect(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	if err := migrations.Up(ctx, pool); err != nil {
		log.Fatal(err)
	}
	log.Println("migrations applied")
}
