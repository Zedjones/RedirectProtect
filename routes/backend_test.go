package routes

import (
	"errors"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/franela/goblin"
	"github.com/go-bongo/bongo"
	mock_echo "github.com/zedjones/redirectprotect/test/mocks"

	"github.com/golang/mock/gomock"
)

func startMongoDocker(start bool) {
	composePath := os.Getenv("COMPOSEPATH")
	var cmd *exec.Cmd
	if start {
		cmd = exec.Command("/usr/bin/docker-compose", "up")
	} else {
		cmd = exec.Command("/usr/bin/docker-compose", "down")
	}
	cmd.Dir = composePath
	dur, _ := time.ParseDuration(".5s")
	time.Sleep(dur)
	cmd.Start()
}

func TestRegisterURL(t *testing.T) {
	startMongoDocker(true)
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
	})
}

func testNoURLPassphrase(t *testing.T) error {
	ctrl := gomock.NewController(t)
	m := mock_echo.NewMockContext(ctrl)

	gomock.InOrder(
		m.EXPECT().QueryParam("url"),
		m.EXPECT().QueryParam("passphrase"),
		m.EXPECT().QueryParam("ttl"),
	)

	m.EXPECT().String(http.StatusBadRequest, "URL or passphrase not provided")

	return RegisterURL(m)
}

func testBadDuration(t *testing.T) error {
	ctrl := gomock.NewController(t)
	m := mock_echo.NewMockContext(ctrl)

	gomock.InOrder(
		m.EXPECT().QueryParam("url").Return("some_url"),
		m.EXPECT().QueryParam("passphrase").Return("some_passphrase"),
		m.EXPECT().QueryParam("ttl").Return("not_a_duration"),
	)

	m.EXPECT().String(http.StatusInternalServerError, "Error parsing duration")

	return RegisterURL(m)
}

func testBadURL(t *testing.T) error {
	ctrl := gomock.NewController(t)
	m := mock_echo.NewMockContext(ctrl)

	gomock.InOrder(
		m.EXPECT().QueryParam("url").Return("magnet:?some_magnet_url"),
		m.EXPECT().QueryParam("passphrase").Return("some_passphrase"),
		m.EXPECT().QueryParam("ttl").Return(""),
	)

	m.EXPECT().String(http.StatusBadRequest, "Invalid URL provided")

	return RegisterURL(m)
}

func testGeneratePasswordFail(t *testing.T) error {
	ctrl := gomock.NewController(t)
	m := mock_echo.NewMockContext(ctrl)

	oldGenerate := generateFromPassword
	defer func() { generateFromPassword = oldGenerate }()

	generateFromPassword = func(pass []byte, cost int) ([]byte, error) {
		return nil, errors.New("Failed to generate password")
	}

	gomock.InOrder(
		m.EXPECT().QueryParam("url").Return("google.com"),
		m.EXPECT().QueryParam("passphrase").Return("some_passphrase"),
		m.EXPECT().QueryParam("ttl").Return(""),
	)

	m.EXPECT().String(http.StatusInternalServerError, "Failed to generate password")

	return RegisterURL(m)
}

func testDatabaseConnFail(t *testing.T) error {
	ctrl := gomock.NewController(t)
	m := mock_echo.NewMockContext(ctrl)

	oldGenerate := generateFromPassword
	defer func() { generateFromPassword = oldGenerate }()

	generateFromPassword = func(pass []byte, cost int) ([]byte, error) {
		return []byte("some test"), nil
	}

	oldGetConn := getConnection
	defer func() { getConnection = oldGetConn }()

	getConnection = func() (*bongo.Connection, error) {
		return nil, errors.New("some error")
	}

	gomock.InOrder(
		m.EXPECT().QueryParam("url").Return("google.com"),
		m.EXPECT().QueryParam("passphrase").Return("some_passphrase"),
		m.EXPECT().QueryParam("ttl").Return(""),
	)

	m.EXPECT().String(http.StatusInternalServerError, "Failed to acquire database connection")

	return RegisterURL(m)
}
