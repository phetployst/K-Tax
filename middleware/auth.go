package middleware

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

func BasicAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		username := c.Request().Header.Get("username")
		password := c.Request().Header.Get("password")

		adminUsername, usernameExists := os.LookupEnv("ADMIN_USERNAME")
		adminPassword, passwordExists := os.LookupEnv("ADMIN_PASSWORD")

		if !usernameExists || !passwordExists || username != adminUsername || password != adminPassword {
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
		}

		return next(c)
	}
}
