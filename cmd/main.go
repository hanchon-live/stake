package main

import (
	"github.com/gorilla/sessions"
	"github.com/hanchon-live/stake/src/components"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	secret_key string = "cookie-key"
)

func main() {
	server := echo.New()
	server.Static("/public/assets/", "./public/assets/")

	server.Use(middleware.Logger())
	server.Use(session.Middleware(sessions.NewCookieStore([]byte(secret_key))))

	server.GET("/", func(c echo.Context) error {
		component := components.Body()
		return component.Render(c.Request().Context(), c.Response().Writer)
	})

	server.Logger.Fatal(server.Start(":3000"))
}
