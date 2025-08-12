package ysgo

import (
	"fmt"
	"log"
	"testing"
)

var (
	testUser = "demo"
	testPass = "demo"
)

func TestYSClient_Login(t *testing.T) {
	client := NewClient(testUser, testPass)

	loginResp, err := client.Login()
	if err != nil {
		log.Printf("Login failed: %v", err)
		return
	}

	fmt.Printf("Login successful. Admin: %v, Directory: %d\n",
		loginResp.User.IsAdmin, loginResp.Directory.Number)
	fmt.Printf("Upload URL: %s\n", loginResp.Space.UploadAddress)
}

func TestYSClient_GetFileList(t *testing.T) {
	client := NewClient(testUser, testPass)

	fileListReq := &FileListRequest{
		DirectoryNumber: "1445856",
		OpenPassword:    "",
		FileNumber:      "0",
	}

	files, err := client.GetFileList(fileListReq)
	if err != nil {
		log.Printf("Get file list failed: %v", err)
		return
	}

	fmt.Printf("File list response: %s\n", string(files))
}

func TestYSClient_PeriodicCheck(t *testing.T) {
	client := NewClient(testUser, testPass)

	checkReq := &PeriodicCheckRequest{
		DirectoryNumber: "1359095",
		OpenPassword:    "",
		FileNumber:      "0",
		UpdateModTime:   "1359095,2024-11-05T22:50:46.227",
	}

	err := client.PeriodicCheck(checkReq)
	if err != nil {
		log.Printf("Periodic check failed: %v", err)
		return
	}

	fmt.Println("Periodic check successful")
}

func TestYSClient_DirectoryManagement(t *testing.T) {
	client := NewClient(testUser, testPass)

	dirReq := &DirectorySettingsRequest{
		Number:       "1367421",
		Title:        "merlin-plugin",
		Description:  "梅林R8000插件备份",
		OpenPassword: "",
		SortNumber:   "0",
		OpenMethod:   "0",
		FileSort:     "1",
		Permissions:  "000100",
		Time:         "2021-01-26T13:06:00",
		SortWeight:   "0",
	}

	err := client.AddDirectory(dirReq)
	if err != nil {
		log.Printf("Add directory failed: %v", err)
		return
	}

	fmt.Println("Directory added successfully")
}

func TestYSClient_CompleteWorkflow(t *testing.T) {
	client := NewClient(testUser, testPass, WithApiBaseUrl("http://c6.ysepan.com"))

	loginResp, err := client.Login()
	if err != nil {
		log.Printf("Login failed: %v", err)
		return
	}

	fmt.Printf("Logged in successfully. Download token: %s\n", loginResp.Directory.DownloadToken)

	fileListReq := &FileListRequest{
		DirectoryNumber: fmt.Sprintf("%d", loginResp.Directory.Number),
		OpenPassword:    "",
		FileNumber:      "0",
	}

	files, err := client.GetFileList(fileListReq)
	if err != nil {
		log.Printf("Get file list failed: %v", err)
	} else {
		fmt.Printf("Retrieved file list: %d bytes\n", len(files))
	}

	err = client.Logout()
	if err != nil {
		log.Printf("Logout failed: %v", err)
	} else {
		fmt.Println("Logged out successfully")
	}
}
