package routes

import (
	"errors"
	"net/http"
	"testing"

	"github.com/franela/goblin"
	"github.com/google/uuid"
	"github.com/zedjones/redirectprotect/db"
	"github.com/zedjones/redirectprotect/test/mocks"

	"github.com/golang/mock/gomock"
)

func TestRegisterURL(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("Register URL", func() {
		g.It("should fail when no passphrase or URL are provided", func() {
			g.Assert(testNoURLPassphrase(t)).Equal(nil)
		})
		g.It("should fail when a bad duration is provided", func() {
			g.Assert(testBadDuration(t)).Equal(nil)
		})
		g.It("should fail when a bad URL is provided", func() {
			g.Assert(testBadURL(t)).Equal(nil)
		})
		g.It("should fail when bcrypt fails to hash the password", func() {
			g.Assert(testGeneratePasswordFail(t)).Equal(nil)
		})
		g.It("should fail when it cannot acquire a database connection", func() {
			g.Assert(testDatabaseConnFail(t)).Equal(nil)
		})
		g.It("should fail when collection fails to save redirect", func() {
			g.Assert(testDatabaseSaveFail(t)).Equal(nil)
		})
		g.It("succeeds when nothing above is happening", func() {
			g.Assert(testSuccessCase(t, g)).Equal(nil)
		})
	})
}

func testNoURLPassphrase(t *testing.T) error {
	ctrl := gomock.NewController(t)
	mockEcho := mocks.NewMockContext(ctrl)

	gomock.InOrder(
		mockEcho.EXPECT().QueryParam("url"),
		mockEcho.EXPECT().QueryParam("passphrase"),
		mockEcho.EXPECT().QueryParam("ttl"),
	)

	mockEcho.EXPECT().String(http.StatusBadRequest, "URL or passphrase not provided")

	return RegisterURL(mockEcho)
}

func testBadDuration(t *testing.T) error {
	ctrl := gomock.NewController(t)
	mockEcho := mocks.NewMockContext(ctrl)

	gomock.InOrder(
		mockEcho.EXPECT().QueryParam("url").Return("some_url"),
		mockEcho.EXPECT().QueryParam("passphrase").Return("some_passphrase"),
		mockEcho.EXPECT().QueryParam("ttl").Return("not_a_duration"),
	)

	mockEcho.EXPECT().String(http.StatusInternalServerError, "Error parsing duration")

	return RegisterURL(mockEcho)
}

func testBadURL(t *testing.T) error {
	ctrl := gomock.NewController(t)
	mockEcho := mocks.NewMockContext(ctrl)

	gomock.InOrder(
		mockEcho.EXPECT().QueryParam("url").Return("magnet:?some_magnet_url"),
		mockEcho.EXPECT().QueryParam("passphrase").Return("some_passphrase"),
		mockEcho.EXPECT().QueryParam("ttl").Return(""),
	)

	mockEcho.EXPECT().String(http.StatusBadRequest, "Invalid URL provided")

	return RegisterURL(mockEcho)
}

func testGeneratePasswordFail(t *testing.T) error {
	ctrl := gomock.NewController(t)
	mockEcho := mocks.NewMockContext(ctrl)

	oldGenerate := generateFromPassword
	defer func() { generateFromPassword = oldGenerate }()

	generateFromPassword = func(pass []byte, cost int) ([]byte, error) {
		return nil, errors.New("Failed to generate password")
	}

	gomock.InOrder(
		mockEcho.EXPECT().QueryParam("url").Return("google.com"),
		mockEcho.EXPECT().QueryParam("passphrase").Return("some_passphrase"),
		mockEcho.EXPECT().QueryParam("ttl").Return(""),
	)

	mockEcho.EXPECT().String(http.StatusInternalServerError, "Failed to generate password")

	return RegisterURL(mockEcho)
}

func testDatabaseConnFail(t *testing.T) error {
	ctrl := gomock.NewController(t)
	mockEcho := mocks.NewMockContext(ctrl)

	oldGenerate := generateFromPassword
	defer func() { generateFromPassword = oldGenerate }()

	generateFromPassword = func(pass []byte, cost int) ([]byte, error) {
		return []byte("some test"), nil
	}

	oldGetConn := getConnection
	defer func() { getConnection = oldGetConn }()

	getConnection = func() (db.Connection, error) {
		return nil, errors.New("some error")
	}

	gomock.InOrder(
		mockEcho.EXPECT().QueryParam("url").Return("google.com"),
		mockEcho.EXPECT().QueryParam("passphrase").Return("some_passphrase"),
		mockEcho.EXPECT().QueryParam("ttl").Return(""),
	)

	mockEcho.EXPECT().String(http.StatusInternalServerError, "Failed to acquire database connection")

	return RegisterURL(mockEcho)
}

func testDatabaseSaveFail(t *testing.T) error {
	ctrl := gomock.NewController(t)
	mockEcho := mocks.NewMockContext(ctrl)
	mockConnection := mocks.NewMockConnection(ctrl)
	mockCollection := mocks.NewMockCollection(ctrl)

	oldGenerate := generateFromPassword
	defer func() { generateFromPassword = oldGenerate }()

	generateFromPassword = func(pass []byte, cost int) ([]byte, error) {
		return []byte("some test"), nil
	}

	oldGetConn := getConnection
	defer func() { getConnection = oldGetConn }()

	getConnection = func() (db.Connection, error) {
		return mockConnection, nil
	}

	oldUUID := uuidNew
	defer func() { uuidNew = oldUUID }()

	myUUID := uuid.New()

	uuidNew = func() uuid.UUID {
		return myUUID
	}

	gomock.InOrder(
		mockEcho.EXPECT().QueryParam("url").Return("google.com"),
		mockEcho.EXPECT().QueryParam("passphrase").Return("some_passphrase"),
		mockEcho.EXPECT().QueryParam("ttl").Return(""),
	)

	expectedRedirect := db.Redirect{URL: "http://google.com", Password: "some test", TTL: "0s",
		Path: myUUID.String()}

	mockCollection.EXPECT().Save(&expectedRedirect).Return(errors.New("some error"))
	mockConnection.EXPECT().Collection("redirections").Return(mockCollection)

	mockEcho.EXPECT().String(http.StatusInternalServerError, "Failed to save redirect to the database")

	return RegisterURL(mockEcho)
}

func testSuccessCase(t *testing.T, g *goblin.G) error {
	ctrl := gomock.NewController(t)
	mockEcho := mocks.NewMockContext(ctrl)
	mockConnection := mocks.NewMockConnection(ctrl)
	mockCollection := mocks.NewMockCollection(ctrl)

	oldGenerate := generateFromPassword
	defer func() { generateFromPassword = oldGenerate }()

	generateFromPassword = func(pass []byte, cost int) ([]byte, error) {
		return []byte("some test"), nil
	}

	oldGetConn := getConnection
	defer func() { getConnection = oldGetConn }()

	getConnection = func() (db.Connection, error) {
		return mockConnection, nil
	}

	oldUUID := uuidNew
	defer func() { uuidNew = oldUUID }()

	myUUID := uuid.New()

	uuidNew = func() uuid.UUID {
		return myUUID
	}

	oldStartTimeCheck := startTimeCheck
	defer func() { startTimeCheck = oldStartTimeCheck }()

	startTimeCheck = func(redir *db.Redirect, coll db.Collection) error {
		g.Assert(coll).Equal(mockCollection)
		return nil
	}

	gomock.InOrder(
		mockEcho.EXPECT().QueryParam("url").Return("google.com"),
		mockEcho.EXPECT().QueryParam("passphrase").Return("some_passphrase"),
		mockEcho.EXPECT().QueryParam("ttl").Return(""),
	)

	expectedRedirect := db.Redirect{URL: "http://google.com", Password: "some test", TTL: "0s",
		Path: myUUID.String()}

	mockCollection.EXPECT().Save(&expectedRedirect).Return(nil)
	mockConnection.EXPECT().Collection("redirections").Return(mockCollection)

	mockEcho.EXPECT().String(http.StatusOK, expectedRedirect.Path)

	return RegisterURL(mockEcho)
}
