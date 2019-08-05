package main

import (
	"github.com/labstack/echo"
	"github.com/zedjones/redirectprotect/internal"
	"github.com/zedjones/redirectprotect/routes"
)

func main() {
	go internal.AddChecks()
	e := echo.New()
	e.POST("/add_redirect", routes.RegisterURL)
	e.POST("/check_pass", routes.CheckPassphrase)
	e.GET("*", routes.GetRedirect)
	e.Start(":1234")
}
