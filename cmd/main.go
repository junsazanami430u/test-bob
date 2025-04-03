package main

import (
	"github.com/labstack/echo/v4"
)

func main() {
	r := echo.New()
	r.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello World!")
	})
	r.GET("/ping", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"message": "pong",
		})
	})
	r.Start(":8080")
}
