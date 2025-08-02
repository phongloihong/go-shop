package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/phongloihong/go-shop/services/user-service/internal/config"
	"github.com/phongloihong/go-shop/services/user-service/internal/delivery/connect"
	"github.com/phongloihong/go-shop/services/user-service/internal/infrastructure/database/postgres"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Println("Error loading configuration:", err)
		return
	}

	conn, err := postgres.NewConnection(context.Background(), cfg.Database)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	defer func() {
		conErr := conn.Close(context.Background())
		if conErr != nil {
			fmt.Println("Error when closing database connection:", conErr)
		}
	}()

	fmt.Println("Connected to database successfully")

	startConnectServer(cfg, conn)
}

func startConnectServer(cfg *config.Config, conn *pgx.Conn) {
	server := connect.StartConnect(conn)
	server.Addr = fmt.Sprintf(":%d", cfg.Server.Port)

	// handle graceful shutdown
	go func() {
		fmt.Println("Starting server on", server.Addr)

		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// wait for shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	server.Shutdown(ctx)

	fmt.Println("Server gracefully stopped")
}
