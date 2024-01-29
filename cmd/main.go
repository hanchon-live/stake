package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/hanchon-live/stake/src/components"
	"github.com/hanchon-live/stake/src/components/wallet"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	secret_key string = "cookie-key"
	port       string = "3000"
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

	server.POST("/wallets", func(c echo.Context) error {
		fmt.Println("--------")
		c.Request().ParseForm()
		providers, ok := c.Request().Form["providers"]
		if !ok {
			providers = []string{}
		}
		return wallet.WalletList(providers).Render(c.Request().Context(), c.Response().Writer)
	})

	server.POST("/currentwallet", func(c echo.Context) error {
		c.Request().ParseForm()
		account, ok := c.Request().Form["accounts"]
		fmt.Println(account)
		value := "0x..."
		if ok && len(account) > 0 {
			value = account[0]
		}
		return c.String(http.StatusOK, value)
	})

	server.Logger.Fatal(server.Start(":" + port))
}
