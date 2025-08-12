package ysgo

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	periodicCheckPath = "/nkj/dsdq.aspx"
	fileListPath      = "/nkj/wj.aspx"
	directoryPath     = "/nkj/ml.aspx"
)

func (c *YSClient) PeriodicCheck(req *PeriodicCheckRequest) error {
	reqURL := c.apiBaseUrl + periodicCheckPath

	form := url.Values{}
	form.Add("mlbh", req.DirectoryNumber)
	form.Add("kqmm", req.OpenPassword)
	form.Add("wjbh", req.FileNumber)
	form.Add("gxxmsj", req.UpdateModTime)

	httpReq, err := http.NewRequest("POST", reqURL, bytes.NewBufferString(form.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create periodic check request: %w", err)
	}

	token := c.generateAuthToken()
	httpReq.Header.Set("Authorization", token.String())
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	httpReq.Header.Set("DNT", "1")
	httpReq.Header.Set("Pragma", "no-cache")

	resp, err := c.c.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to execute periodic check request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return NewYSError(resp.StatusCode, "Periodic check failed", resp.Status)
	}

	return nil
}

func (c *YSClient) GetFileList(req *FileListRequest) ([]byte, error) {
	reqURL := c.apiBaseUrl + fileListPath

	form := url.Values{}
	form.Add("mlbh", req.DirectoryNumber)
	form.Add("kqmm", req.OpenPassword)
	form.Add("wjbh", req.FileNumber)

	httpReq, err := http.NewRequest("POST", reqURL, bytes.NewBufferString(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create file list request: %w", err)
	}

	token := c.generateAuthToken()
	httpReq.Header.Set("Authorization", token.String())
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	httpReq.Header.Set("DNT", "1")
	httpReq.Header.Set("Pragma", "no-cache")

	resp, err := c.c.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute file list request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, NewYSError(resp.StatusCode, "Get file list failed", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read file list response body: %w", err)
	}

	return body, nil
}

func (c *YSClient) GetDirectoryInfo(directoryNumber string) ([]byte, error) {
	reqURL := c.apiBaseUrl + directoryPath + "?bjbh=" + directoryNumber

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create directory info request: %w", err)
	}

	token := c.generateAuthToken()
	req.Header.Set("Authorization", token.String())
	req.Header.Set("DNT", "1")
	req.Header.Set("Pragma", "no-cache")

	resp, err := c.c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute directory info request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, NewYSError(resp.StatusCode, "Get directory info failed", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory info response body: %w", err)
	}

	return body, nil
}

func (c *YSClient) SetDirectorySettings(operation string, req *DirectorySettingsRequest) error {
	reqURL := c.apiBaseUrl + directoryPath + "?cz=" + operation

	form := url.Values{}
	form.Add("bh", req.Number)
	form.Add("bt", req.Title)
	form.Add("sm", req.Description)
	form.Add("kqmm", req.OpenPassword)
	form.Add("pxbh", req.SortNumber)
	form.Add("kqfs", req.OpenMethod)
	form.Add("wjpx", req.FileSort)
	form.Add("qx", req.Permissions)
	form.Add("sj", req.Time)
	form.Add("pxz", req.SortWeight)

	httpReq, err := http.NewRequest("POST", reqURL, bytes.NewBufferString(form.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create directory settings request: %w", err)
	}

	token := c.generateAuthToken()
	httpReq.Header.Set("Authorization", token.String())
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	httpReq.Header.Set("DNT", "1")
	httpReq.Header.Set("Pragma", "no-cache")

	resp, err := c.c.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to execute directory settings request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return NewYSError(resp.StatusCode, "Set directory settings failed", resp.Status)
	}

	return nil
}

func (c *YSClient) AddDirectory(req *DirectorySettingsRequest) error {
	return c.SetDirectorySettings("add", req)
}

func (c *YSClient) UpdateDirectory(req *DirectorySettingsRequest) error {
	return c.SetDirectorySettings("edit", req)
}

func (c *YSClient) DeleteDirectory(req *DirectorySettingsRequest) error {
	return c.SetDirectorySettings("del", req)
}
