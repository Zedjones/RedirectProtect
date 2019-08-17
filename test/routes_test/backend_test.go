package routes_test

import (
	"net/http"
	"testing"

	"github.com/zedjones/redirectprotect/routes"

	mock_echo "github.com/zedjones/redirectprotect/test/mocks"

	"github.com/golang/mock/gomock"
)

func TestRegisterURL_NoURLPassphrase(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := mock_echo.NewMockContext(ctrl)

	gomock.InOrder(
		m.EXPECT().QueryParam("url"),
		m.EXPECT().QueryParam("passphrase"),
		m.EXPECT().QueryParam("ttl"),
	)

	m.EXPECT().String(http.StatusBadRequest, "URL or passphrase not provided")

	routes.RegisterURL(m)
}
