package internal

import (
	"errors"
	"testing"
	"time"

	"github.com/go-bongo/bongo"
	"github.com/golang/mock/gomock"
	"github.com/zedjones/redirectprotect/db"
	"github.com/zedjones/redirectprotect/test/mocks"

	"github.com/franela/goblin"
)

func TestStartTimeCheck(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("Start Time Check", func() {
		g.It("should fail when an invalid duration is provided", func() {
			g.Assert(testParseDurationFail()).Equal(errors.New("some error"))
		})
		g.It("should return when a duration of 0s is provided", func() {
			g.Assert(testNoDuration()).Equal(nil)
		})
		g.It("should succeed when nothing above is happening", func() {
			g.Assert(timeCheckTestSuccess(t, g)).Equal(nil)
		})
	})
}

func testParseDurationFail() error {
	oldDuration := parseDuration
	defer func() { parseDuration = oldDuration }()

	parseDuration = func(s string) (time.Duration, error) {
		return 0, errors.New("some error")
	}

	return StartTimeCheck(&db.Redirect{}, db.BongoCollection{})
}

func testNoDuration() error {
	return StartTimeCheck(&db.Redirect{TTL: "0s"}, db.BongoCollection{})
}

func timeCheckTestSuccess(t *testing.T, g *goblin.G) error {
	ctrl := gomock.NewController(t)
	mockCollection := mocks.NewMockCollection(ctrl)

	oldNow := now
	defer func() { now = oldNow }()

	currTime := time.Now()
	now = func() time.Time {
		return currTime
	}

	redir := db.Redirect{
		TTL: "1s",
		DocumentBase: bongo.DocumentBase{
			Created: currTime,
		},
	}

	oldSleep := sleep
	defer func() { sleep = oldSleep }()

	sleep = func(dur time.Duration) {
		zero, _ := time.ParseDuration("1s")
		g.Assert(dur).Equal(zero)
	}

	mockCollection.EXPECT().DeleteDocument(&redir).Return(nil)

	return StartTimeCheck(&redir, mockCollection)
}
