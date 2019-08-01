package routes

import (
	"net/http"

	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

func RegisterURL(c echo.Context) error {
	url := c.QueryParam("url")
	passphrase := c.QueryParam("passphrase")
	bytes, err := bcrypt.GenerateFromPassword([]byte(passphrase), 17)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	url, bytes = url, bytes
	return nil
}
