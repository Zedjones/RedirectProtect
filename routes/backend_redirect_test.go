package routes

import (
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/franela/goblin"
	"github.com/golang/mock/gomock"
	"github.com/zedjones/redirectprotect/db"
	"github.com/zedjones/redirectprotect/test/mocks"
	"gopkg.in/mgo.v2/bson"
)

func TestGetRedirect(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("Get Redirect", func() {
		g.It("should fail when it cannot acquire a database connection", func() {
			g.Assert(redirectTestDatabaseConnFail(t)).Equal(nil)
		})
		g.It("should fail when path does not exist in DB", func() {
			g.Assert(testPathIncorrect(t)).Equal(nil)
		})
		g.It("should succeed when nothing above is happening", func() {
			g.Assert(redirectTestSuccessCase(t)).Equal(nil)
		})
	})
}

func redirectTestDatabaseConnFail(t *testing.T) error {
	ctrl := gomock.NewController(t)
	mockEcho := mocks.NewMockContext(ctrl)

	oldGetConn := getConnection
	defer func() { getConnection = oldGetConn }()

	getConnection = func() (db.Connection, error) {
		return nil, errors.New("some error")
	}

	mockEcho.EXPECT().String(http.StatusInternalServerError, "Failed to acquire database connection")

	return GetRedirect(mockEcho)
}

func testPathIncorrect(t *testing.T) error {
	ctrl := gomock.NewController(t)
	mockEcho := mocks.NewMockContext(ctrl)
	mockConnection := mocks.NewMockConnection(ctrl)
	mockCollection := mocks.NewMockCollection(ctrl)

	oldGetConn := getConnection
	defer func() { getConnection = oldGetConn }()

	getConnection = func() (db.Connection, error) {
		return mockConnection, nil
	}

	mockEcho.EXPECT().Request().Return(
		&http.Request{URL: &url.URL{Path: "some path"}},
	)

	redir := db.Redirect{}
	mockCollection.EXPECT().FindOne(bson.M{"path": "some path"}, &redir).Return(errors.New("Some error"))
	mockConnection.EXPECT().Collection("redirections").Return(mockCollection)

	mockEcho.EXPECT().String(http.StatusBadRequest, "Shortened URL does not exist")

	return GetRedirect(mockEcho)
}

func redirectTestSuccessCase(t *testing.T) error {
	ctrl := gomock.NewController(t)
	mockEcho := mocks.NewMockContext(ctrl)
	mockConnection := mocks.NewMockConnection(ctrl)
	mockCollection := mocks.NewMockCollection(ctrl)

	oldGetConn := getConnection
	defer func() { getConnection = oldGetConn }()

	getConnection = func() (db.Connection, error) {
		return mockConnection, nil
	}

	mockEcho.EXPECT().Request().Return(
		&http.Request{URL: &url.URL{Path: "some path"}},
	)

	redir := db.Redirect{}
	mockCollection.EXPECT().FindOne(bson.M{"path": "some path"}, &redir)
	mockConnection.EXPECT().Collection("redirections").Return(mockCollection)

	mockEcho.EXPECT().Render(http.StatusOK, "redir.html", nil)

	return GetRedirect(mockEcho)
}
