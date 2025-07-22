package router

import "github.com/labstack/echo/v4"

func InitRouter(e *echo.Echo) {
	e.GET("/health", func(c echo.Context) error {
		return c.String(200, "OK")
	})

	apiVersion1(e)
}

func apiVersion1(e *echo.Echo) {
	v1 := e.Group("/v1")
	initUserRoutesV1(v1)
}
