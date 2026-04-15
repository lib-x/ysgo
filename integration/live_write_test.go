package integration

import (
	"bytes"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/lib-x/ysgo"
)

func TestLiveSandboxWriteFlow(t *testing.T) {
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
	title := fmt.Sprintf("ysgo-live-sandbox-%d", now)
	fileName := fmt.Sprintf("live-%d.txt", now)
	content := []byte("live integration content")

	if err := client.AddDirectory(&ysgo.DirectorySettingsRequest{
		Number:       "0",
		Title:        title,
		Description:  "live integration sandbox",
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

	res, err := client.UploadBytes(dirNumber, "", "nested/live", fileName, content)
	if err != nil {
		t.Fatalf("UploadBytes failed: %v", err)
	}
	defer func() {
		_ = client.DeleteFiles(&ysgo.DeleteFilesRequest{
			DirectoryNumber: dirNumber,
			FileNumbers:     []string{fmt.Sprintf("%d", res.FileNumber)},
		})
	}()

	files, err := client.GetFileListParsed(&ysgo.FileListRequest{DirectoryNumber: dirNumber, FileNumber: "0"})
	if err != nil {
		t.Fatalf("GetFileListParsed failed: %v", err)
	}

	var target ysgo.RemoteFile
	found := false
	for _, f := range files.Files {
		if f.Number == res.FileNumber {
			target = f
			found = true
			break
		}
	}
	if !found {
		t.Fatal("uploaded file not found")
	}

	downloadURL, err := client.BuildDownloadURL(mustInt(dirNumber), files.Directory.DownloadToken, target, nil)
	if err != nil {
		t.Fatalf("BuildDownloadURL failed: %v", err)
	}
	body, err := client.DownloadBytes(downloadURL)
	if err != nil {
		t.Fatalf("DownloadBytes failed: %v", err)
	}
	if !bytes.Equal(body, content) {
		t.Fatalf("download mismatch: got %q want %q", body, content)
	}
}

func mustInt(s string) int {
	var n int
	_, err := fmt.Sscanf(s, "%d", &n)
	if err != nil {
		panic(err)
	}
	return n
}
