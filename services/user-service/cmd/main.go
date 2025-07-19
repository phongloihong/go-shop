package main

import (
	"fmt"
	"github.com/phongloihong/go-shop/services/user-service/internal/config"
	"github.com/phongloihong/go-shop/services/user-service/internal/delivery/http/router"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Println("Error loading configuration:", err)
		return
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	router.InitRouter(e)
	
	// Check if server config is nil and provide default port
	port := 8080
	if cfg.Server != nil {
		port = cfg.Server.Port
	}
	
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}
