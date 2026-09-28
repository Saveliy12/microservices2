package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com.Saveliy12/microservices2.git/order/internal/app"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("warning: failed to load .env file: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	a, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("failed to init app: %v", err)
	}

	if err := a.Run(ctx); err != nil {
		log.Printf("app exited with error: %v", err)
	}
}
