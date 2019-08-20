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

//Doing this because I want to be able to mock these functions
var (
	generateFromPassword   = bcrypt.GenerateFromPassword
	compareHashAndPassword = bcrypt.CompareHashAndPassword
	getConnection          = db.GetConnection
	startTimeCheck         = internal.StartTimeCheck
	uuidNew                = uuid.New
)

//RegisterURL registers a URL in the database with the corresponding password and
//duration, if provided. 15 rounds of bcrypt are done on the password before
//storing it, and the shortened path returned is a generated UUID
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
	if (strings.Contains(url, ":/") || strings.Contains(url, ":?")) &&
		!(strings.HasPrefix(url, "https://") || strings.HasPrefix(url, "http://")) {
		return c.String(http.StatusBadRequest, "Invalid URL provided")
	} else if !strings.Contains(url, ":/") {
		url = "http://" + url
	}

	bytes, err := generateFromPassword([]byte(passphrase), 15)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	newRedirect := db.Redirect{URL: url, Password: string(bytes),
		TTL: duration.String(), Path: uuidNew().String()}

	connection, err := getConnection()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to acquire database connection")
	}
	collection := connection.Collection(db.CollectionName)

	err = collection.Save(&newRedirect)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to save redirect to the database")
	}
	go startTimeCheck(&newRedirect, collection)
	return c.String(http.StatusOK, newRedirect.Path)
}

//GetRedirect is the default GET handler, and simply returns a page asking for the password if
//the shortened URL exists in the DB. Once the user inputs the password, a call to CheckPassphrase
//is made
func GetRedirect(c echo.Context) error {
	redir := &db.Redirect{}
	connection, err := getConnection()
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

//CheckPassphrase checks that the passphrase is valid for the provided path, and returns a JSON object
//containing the URL that was shortened/encrypted
func CheckPassphrase(c echo.Context) error {
	redir := &db.Redirect{}
	path := c.QueryParam("path")
	passphrase := c.QueryParam("passphrase")
	connection, err := getConnection()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to acquire database connection")
	}
	err = connection.Collection(db.CollectionName).FindOne(bson.M{"path": path}, redir)
	if err != nil {
		return c.String(http.StatusBadRequest, "Shortened URL does not exist")
	}
	err = compareHashAndPassword([]byte(redir.Password), []byte(passphrase))
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad password provided.")
	}
	return c.JSON(http.StatusOK, map[string]string{"url": redir.URL})
}
