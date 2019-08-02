package routes

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/zedjones/redirectprotect/db"
	"golang.org/x/crypto/bcrypt"
)

func RegisterURL(c echo.Context) error {
	url := c.QueryParam("url")
	passphrase := c.QueryParam("passphrase")
	bytes, err := bcrypt.GenerateFromPassword([]byte(passphrase), 17)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	newRedirect := db.Redirect{URL: url, Password: string(bytes)}
	return nil
}
