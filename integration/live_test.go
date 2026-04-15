package integration

import (
	"os"
	"testing"

	"github.com/lib-x/ysgo"
)

func TestLiveSessionAndList(t *testing.T) {
	user := os.Getenv("YSGO_LIVE_USER")
	pass := os.Getenv("YSGO_LIVE_PASS")
	if user == "" || pass == "" {
		t.Skip("missing YSGO_LIVE_USER/YSGO_LIVE_PASS")
	}

	client := ysgo.NewClient(user, pass)
	if _, err := client.InitSession(); err != nil {
		t.Fatalf("InitSession failed: %v", err)
	}
	if _, err := client.GetDirectoryListParsed(); err != nil {
		t.Fatalf("GetDirectoryListParsed failed: %v", err)
	}
}
