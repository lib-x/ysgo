package integration

import (
	"os"
	"testing"
	"time"

	"github.com/lib-x/ysgo"
)

func TestLiveSpacePasswordSession(t *testing.T) {
	user := os.Getenv("YSGO_LIVE_USER")
	pass := os.Getenv("YSGO_LIVE_PASS")
	spacePassword := os.Getenv("YSGO_LIVE_SPACE_PASSWORD")
	if user == "" || pass == "" || spacePassword == "" {
		t.Skip("missing YSGO_LIVE_USER/YSGO_LIVE_PASS/YSGO_LIVE_SPACE_PASSWORD")
	}

	client := ysgo.NewClient(
		user,
		pass,
		ysgo.WithSpacePassword(spacePassword),
		ysgo.WithTimeout(15*time.Second),
	)
	resp, err := client.PrepareSession()
	if err != nil {
		t.Fatalf("PrepareSession failed: %v", err)
	}
	if resp.Token == "" {
		t.Fatal("expected session token after space password verification")
	}
	if resp.Space.UploadAddress == "" {
		t.Fatal("expected upload address after space password verification")
	}
}
