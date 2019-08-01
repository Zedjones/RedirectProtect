package routes

import (
	"github.com/labstack/echo"
)

func RegisterURL(c echo.Context) {
	url := c.QueryParam("url")
	passphrase := c.QueryParam("passphrase")
}
