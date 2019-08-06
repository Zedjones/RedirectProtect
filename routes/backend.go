package routes

import (
	"net/http"
	"strings"
	"time"

	"github.com/zedjones/redirectprotect/internal"

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

	if url == "" || passphrase == "" {
		return c.String(http.StatusBadRequest, "URL or passphrase not provided")
	}
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
		TTL: duration.String(), Path: uuid.New().String()}

	connection, err := db.GetConnection()
	collection := connection.Collection(db.CollectionName)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to acquired database connection")
	}

	err = collection.Save(&newRedirect)
	go internal.StartTimeCheck(&newRedirect, collection)
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
	return c.Render(http.StatusOK, "redir.html", nil)
}

func CheckPassphrase(c echo.Context) error {
	redir := &db.Redirect{}
	path := c.QueryParam("path")
	passphrase := c.QueryParam("passphrase")
	connection, err := db.GetConnection()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to acquire database connection")
	}
	err = connection.Collection(db.CollectionName).FindOne(bson.M{"path": path}, redir)
	if err != nil {
		return c.String(http.StatusBadRequest, "Shortened URL does not exist")
	}
	err = bcrypt.CompareHashAndPassword([]byte(redir.Password), []byte(passphrase))
	if err != nil {
		c.Redirect(http.StatusFound, redir.URL)
	}
	return nil
}
