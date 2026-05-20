package main

import (
	"context"
	"log"
	"time"

	"mertani/internal/app"
	"mertani/internal/config"
	"mertani/internal/database"
)

func main() {
	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := database.NewPostgresConnection(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	server := app.NewServer(cfg, db)
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
