package main

import (
	"context"
	"fmt"
	"log"

	"github.com/phongloihong/go-shop/services/user-service/internal/config"
	"github.com/phongloihong/go-shop/services/user-service/internal/delivery/http/router"
	"github.com/phongloihong/go-shop/services/user-service/internal/infrastructure/database/postgres"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

	startHTTPServer(cfg)
}

func startHTTPServer(cfg *config.Config) error {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Initialize routes
	router.InitRouter(e)

	// Check if server config is nil and provide default port
	port := 8080
	if cfg.Server != nil {
		port = cfg.Server.Port
	}

	return e.Start(fmt.Sprintf(":%d", port))
}
