package router

import "github.com/labstack/echo/v4"

func initUserRoutesV1(e *echo.Group) {
	userGroup := e.Group("/users")
	userGroup.POST("/register", nil)
	userGroup.POST("/update", nil)
	userGroup.PUT("/change-password", nil)
	userGroup.GET("/profile", nil)
}
