package routes

import (
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/zedjones/redirectprotect/db"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

func RegisterURL(c echo.Context) error {
	var duration time.Duration
	url := c.QueryParam("url")
	passphrase := c.QueryParam("passphrase")
	durationStr := c.QueryParam("ttl")
	var err error
	if durationStr != "" {
		duration, err = time.ParseDuration(durationStr)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Error parsing duration")
		}
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(passphrase), 17)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	newRedirect := db.Redirect{URL: url, Password: string(bytes),
		TTL: &duration, Path: uuid.New().String()}
	connection, err := db.GetConnection()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to acquired database connection")
	}
	err = connection.Collection(db.CollectionName).Save(&newRedirect)
	return err
}

func GetRedirect(c echo.Context) error {
	var err error
	redir := &db.Redirect{}
	connection, err := db.GetConnection()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to acquire database connection")
	}
	path := strings.TrimPrefix(c.Request().URL.Path, "/")
	err = connection.Collection(db.CollectionName).FindOne(bson.M{"path": path}, redir)
	if err != nil {
		return c.String(http.StatusBadRequest, "Shortened URL does not exist")
	}
	return err
}
