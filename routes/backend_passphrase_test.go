package routes

import (
	"errors"
	"net/http"
	"testing"

	"github.com/franela/goblin"
	"github.com/golang/mock/gomock"
	"github.com/zedjones/redirectprotect/db"
	"github.com/zedjones/redirectprotect/test/mocks"
	"gopkg.in/mgo.v2/bson"
)

func TestCheckPassphrase(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("Check Passphrase", func() {
		g.It("should fail when it cannot acquire a database connection", func() {
			g.Assert(passTestDatabaseConnFail(t)).Equal(nil)
		})
		g.It("should fail when path does not exist in DB", func() {
			g.Assert(passTestPathIncorrect(t)).Equal(nil)
		})
		g.It("should fail when password is incorrect", func() {
			g.Assert(hashPassError(t, g)).Equal(nil)
		})
		g.It("should succeed when nothing above is happening", func() {
			g.Assert(passSuccess(t)).Equal(nil)
		})
	})
}

func passTestDatabaseConnFail(t *testing.T) error {
	ctrl := gomock.NewController(t)
	mockEcho := mocks.NewMockContext(ctrl)

	oldGetConn := getConnection
	defer func() { getConnection = oldGetConn }()

	getConnection = func() (db.Connection, error) {
		return nil, errors.New("some error")
	}

	gomock.InOrder(
		mockEcho.EXPECT().QueryParam("path"),
		mockEcho.EXPECT().QueryParam("passphrase"),
	)

	mockEcho.EXPECT().String(http.StatusInternalServerError, "Failed to acquire database connection")

	return CheckPassphrase(mockEcho)
}

func passTestPathIncorrect(t *testing.T) error {
	ctrl := gomock.NewController(t)
	mockEcho := mocks.NewMockContext(ctrl)
	mockConnection := mocks.NewMockConnection(ctrl)
	mockCollection := mocks.NewMockCollection(ctrl)

	oldGetConn := getConnection
	defer func() { getConnection = oldGetConn }()

	getConnection = func() (db.Connection, error) {
		return mockConnection, nil
	}

	gomock.InOrder(
		mockEcho.EXPECT().QueryParam("path").Return("some path"),
		mockEcho.EXPECT().QueryParam("passphrase"),
	)

	redir := db.Redirect{}
	mockCollection.EXPECT().FindOne(bson.M{"path": "some path"}, &redir).Return(errors.New("Some error"))
	mockConnection.EXPECT().Collection("redirections").Return(mockCollection)

	mockEcho.EXPECT().String(http.StatusBadRequest, "Shortened URL does not exist")

	return CheckPassphrase(mockEcho)
}

func hashPassError(t *testing.T, g *goblin.G) error {
	ctrl := gomock.NewController(t)
	mockEcho := mocks.NewMockContext(ctrl)
	mockConnection := mocks.NewMockConnection(ctrl)
	mockCollection := mocks.NewMockCollection(ctrl)

	oldGetConn := getConnection
	defer func() { getConnection = oldGetConn }()

	getConnection = func() (db.Connection, error) {
		return mockConnection, nil
	}

	oldCompare := compareHashAndPassword
	defer func() { compareHashAndPassword = oldCompare }()

	compareHashAndPassword = func(a []byte, b []byte) error {
		g.Assert(a).Equal([]byte("testing"))
		g.Assert(b).Equal([]byte("my_pass"))
		return errors.New("some error")
	}

	gomock.InOrder(
		mockEcho.EXPECT().QueryParam("path").Return("some path"),
		mockEcho.EXPECT().QueryParam("passphrase").Return("my_pass"),
	)

	mockCollection.EXPECT().FindOne(bson.M{"path": "some path"}, &db.Redirect{}).
		DoAndReturn(func(b bson.M, redir *db.Redirect) error { redir.Password = "testing"; return nil })
	mockConnection.EXPECT().Collection("redirections").Return(mockCollection)

	mockEcho.EXPECT().String(http.StatusBadRequest, "Bad password provided.")

	return CheckPassphrase(mockEcho)
}

func passSuccess(t *testing.T) error {
	ctrl := gomock.NewController(t)
	mockEcho := mocks.NewMockContext(ctrl)
	mockConnection := mocks.NewMockConnection(ctrl)
	mockCollection := mocks.NewMockCollection(ctrl)

	oldGetConn := getConnection
	defer func() { getConnection = oldGetConn }()

	getConnection = func() (db.Connection, error) {
		return mockConnection, nil
	}

	oldCompare := compareHashAndPassword
	defer func() { compareHashAndPassword = oldCompare }()

	compareHashAndPassword = func([]byte, []byte) error {
		return nil
	}

	gomock.InOrder(
		mockEcho.EXPECT().QueryParam("path").Return("some path"),
		mockEcho.EXPECT().QueryParam("passphrase"),
	)

	mockCollection.EXPECT().FindOne(bson.M{"path": "some path"}, &db.Redirect{}).
		DoAndReturn(func(b bson.M, redir *db.Redirect) error { redir.URL = "test_url"; return nil })
	mockConnection.EXPECT().Collection("redirections").Return(mockCollection)

	mockEcho.EXPECT().JSON(http.StatusOK, map[string]string{"url": "test_url"})

	return CheckPassphrase(mockEcho)
}
