package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SSSagatov/AI-chess-trainer/internal/chesscom"
	"github.com/SSSagatov/AI-chess-trainer/internal/config"
	"github.com/SSSagatov/AI-chess-trainer/internal/db/postgres"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.NewDBConfig()
	if err != nil {
		return fmt.Errorf("load database config: %w", err)
	}

	connectCtx, cancelConnect := context.WithTimeout(ctx, 10*time.Second)
	pool, err := postgres.Connect(connectCtx, cfg)
	cancelConnect()
	if err != nil {
		return err
	}
	defer pool.Close()

	gameHandler, err := chesscom.NewHandler(pool)
	if err != nil {
		return fmt.Errorf("initialize game handler: %w", err)
	}
	handler, err := gameHandler.PageHandler()
	if err != nil {
		return err
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	server := &http.Server{
		Addr:              ":" + port,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}
	listener, err := net.Listen("tcp", server.Addr)
	if err != nil {
		return fmt.Errorf("listen on %s: %w", server.Addr, err)
	}

	serveErr := make(chan error, 1)
	go func() {
		serveErr <- server.Serve(listener)
	}()
	log.Printf("HTTP server listening on %s", listener.Addr())

	select {
	case err := <-serveErr:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return fmt.Errorf("serve HTTP: %w", err)
	case <-ctx.Done():
		log.Println("shutting down HTTP server")
	}

	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()
	if err := server.Shutdown(shutdownCtx); err != nil {
		_ = server.Close()
		return fmt.Errorf("shutdown HTTP server: %w", err)
	}
	if err := <-serveErr; err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("serve HTTP: %w", err)
	}
	return nil
}
