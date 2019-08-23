package internal

import (
	"errors"
	"testing"

	"github.com/franela/goblin"
	"github.com/golang/mock/gomock"
	"github.com/zedjones/redirectprotect/db"
	"github.com/zedjones/redirectprotect/test/mocks"
	"gopkg.in/mgo.v2/bson"
)

func TestAddChecks(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("Add Checks", func() {
		g.It("should fail when it cannot acquire a database connection", func() {
			g.Assert(testDatabaseConnFail()).Equal(errors.New("some error"))
		})
		g.It("should fail when it cannot deepcopy a redirection", func() {
			g.Assert(testDeepcopyFail(t, g)).Equal(errors.New("some error"))
		})
		g.It("should succeed when nothing above is happening", func() {
			g.Assert(testSuccess(t, g)).Equal(nil)
		})
	})
}

func testDatabaseConnFail() error {

	oldGetConn := getConnection
	defer func() { getConnection = oldGetConn }()

	getConnection = func() (db.Connection, error) {
		return nil, errors.New("some error")
	}

	return AddChecks()
}

func testDeepcopyFail(t *testing.T, g *goblin.G) error {
	ctrl := gomock.NewController(t)
	mockConnection := mocks.NewMockConnection(ctrl)
	mockCollection := mocks.NewMockCollection(ctrl)
	mockResultSet := mocks.NewMockResultSet(ctrl)

	mockResultSet.EXPECT().Next(&db.Redirect{}).DoAndReturn(
		func(redir *db.Redirect) bool { redir.Path = "test"; return true })
	mockCollection.EXPECT().Find(bson.M{}).Return(mockResultSet)
	mockConnection.EXPECT().Collection("redirections").Return(mockCollection)

	oldGetConn := getConnection
	defer func() { getConnection = oldGetConn }()

	getConnection = func() (db.Connection, error) {
		return mockConnection, nil
	}

	oldCopy := copy
	defer func() { copy = oldCopy }()

	copy = func(new interface{}, old interface{}) error {
		oldRedir := old.(**db.Redirect)
		g.Assert((*oldRedir).Path).Equal("test")
		return errors.New("some error")
	}

	return AddChecks()
}

func testSuccess(t *testing.T, g *goblin.G) error {
	ctrl := gomock.NewController(t)
	mockConnection := mocks.NewMockConnection(ctrl)
	mockCollection := mocks.NewMockCollection(ctrl)
	mockResultSet := mocks.NewMockResultSet(ctrl)

	mockResultSet.EXPECT().Next(&db.Redirect{}).DoAndReturn(
		func(redir *db.Redirect) bool {
			redir.Path = "test"
			return true
		})
	mockResultSet.EXPECT().Next(&db.Redirect{Path: "test"}).DoAndReturn(
		func(redir *db.Redirect) bool {
			redir.Path = "test"
			return false
		})
	mockCollection.EXPECT().Find(bson.M{}).Return(mockResultSet)
	mockConnection.EXPECT().Collection("redirections").Return(mockCollection)

	oldGetConn := getConnection
	defer func() { getConnection = oldGetConn }()

	getConnection = func() (db.Connection, error) {
		return mockConnection, nil
	}

	oldStartTimeCheck := startTimeCheck
	defer func() { startTimeCheck = oldStartTimeCheck }()

	startTimeCheck = func(redir *db.Redirect, collection db.Collection) error {
		g.Assert(redir.Path).Equal("test")
		g.Assert(collection).Equal(mockCollection)
		return nil
	}

	return AddChecks()
}
