package test_routes

import (
	"testing"

	mock_echo "github.com/zedjones/redirectprotect/test/mocks"

	"github.com/golang/mock/gomock"
)

func TestRegisterURL(t *testing.T) {
	ctrl := gomock.NewController(t)

	m := mock_echo.NewMockContext(ctrl)
	m = m
}
