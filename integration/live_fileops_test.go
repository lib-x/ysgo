package integration

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/lib-x/ysgo"
)

func TestLiveSandboxFileOps(t *testing.T) {
	user := os.Getenv("YSGO_LIVE_USER")
	pass := os.Getenv("YSGO_LIVE_PASS")
	if user == "" || pass == "" {
		t.Skip("missing YSGO_LIVE_USER/YSGO_LIVE_PASS")
	}
	if os.Getenv("YSGO_LIVE_WRITE") != "1" {
		t.Skip("set YSGO_LIVE_WRITE=1 to enable sandbox write test")
	}

	client := ysgo.NewClient(user, pass)
	if _, err := client.PrepareAdminSession(); err != nil {
		t.Fatalf("PrepareAdminSession failed: %v", err)
	}

	now := time.Now().Unix()
	title := fmt.Sprintf("ysgo-live-fileops-%d", now)
	if err := client.AddDirectory(&ysgo.DirectorySettingsRequest{
		Number:       "0",
		Title:        title,
		Description:  "live fileops sandbox",
		OpenPassword: "",
		SortNumber:   "0",
		OpenMethod:   "0",
		FileSort:     "1",
		Permissions:  "000101",
		Time:         "",
		SortWeight:   "0",
	}); err != nil {
		t.Fatalf("AddDirectory failed: %v", err)
	}

	dirs, err := client.GetDirectoryListParsed()
	if err != nil {
		t.Fatalf("GetDirectoryListParsed failed: %v", err)
	}
	var dirNumber string
	for _, d := range dirs.List {
		if d.Title == title {
			dirNumber = fmt.Sprintf("%d", d.Number)
			break
		}
	}
	if dirNumber == "" {
		t.Fatal("sandbox directory not found")
	}
	defer func() { _ = client.DeleteDirectory(dirNumber) }()

	entry, err := client.AddEntry(&ysgo.AddEntryRequest{
		DirectoryNumber: dirNumber,
		Title:           "live-title",
		Content:         "https://example.com",
		Sequence:        1000,
		Public:          true,
	})
	if err != nil {
		t.Fatalf("AddEntry failed: %v", err)
	}
	if _, err := client.CreateSubdirectory(dirNumber, "", "nested/live", 1200); err != nil {
		t.Fatalf("CreateSubdirectory failed: %v", err)
	}
	if err := client.MoveEntries(&ysgo.MoveEntriesRequest{
		SourceDirectoryNumber: dirNumber,
		TargetDirectoryNumber: dirNumber,
		SourcePath:            "",
		TargetPath:            "nested/live",
		FileNumbers:           []int{entry.Number},
	}); err != nil {
		t.Fatalf("MoveEntries failed: %v", err)
	}
	if err := client.SetFileVisibility(dirNumber, "", entry.Number, false); err != nil {
		t.Fatalf("SetFileVisibility failed: %v", err)
	}

	files, err := client.GetFileListParsed(&ysgo.FileListRequest{DirectoryNumber: dirNumber, FileNumber: "0"})
	if err != nil {
		t.Fatalf("GetFileListParsed failed: %v", err)
	}
	foundMoved := false
	for _, f := range files.Files {
		if f.Number == entry.Number {
			if f.Subdirectory != "nested/live" {
				t.Fatalf("entry not moved: got %q", f.Subdirectory)
			}
			if f.Visible {
				t.Fatalf("entry visibility not updated")
			}
			foundMoved = true
			break
		}
	}
	if !foundMoved {
		t.Fatal("moved entry not found")
	}
}
