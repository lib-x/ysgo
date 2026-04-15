package ysgo

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

const (
	periodicCheckPath = "/nkj/dsdq.aspx"
	fileListPath      = "/nkj/wj.aspx"
	directoryPath     = "/nkj/ml.aspx"
)

func (c *YSClient) PeriodicCheck(req *PeriodicCheckRequest) error {
	return c.PeriodicCheckContext(context.Background(), req)
}

func (c *YSClient) PeriodicCheckContext(ctx context.Context, req *PeriodicCheckRequest) error {
	form := url.Values{}
	form.Add("mlbh", req.DirectoryNumber)
	form.Add("kqmm", req.OpenPassword)
	form.Add("wjbh", req.FileNumber)
	form.Add("gxxmsj", req.UpdateModTime)

	httpReq, err := c.newRequestWithContext(ctx, http.MethodPost, periodicCheckPath, form)
	if err != nil {
		return err
	}

	if err := c.do(httpReq, nil, http.StatusOK); err != nil {
		return fmt.Errorf("periodic check: %w", err)
	}

	return nil
}

func (c *YSClient) GetFileList(req *FileListRequest) ([]byte, error) {
	return c.GetFileListContext(context.Background(), req)
}

func (c *YSClient) GetFileListContext(ctx context.Context, req *FileListRequest) ([]byte, error) {
	form := url.Values{}
	form.Add("mlbh", req.DirectoryNumber)
	form.Add("kqmm", req.OpenPassword)
	form.Add("wjbh", req.FileNumber)
	if req.IP1 != "" {
		form.Add("ip1", req.IP1)
	}

	httpReq, err := c.newRequestWithContext(ctx, http.MethodPost, fileListPath, form)
	if err != nil {
		return nil, err
	}

	var body []byte
	if err := c.do(httpReq, &body, http.StatusOK); err != nil {
		return nil, fmt.Errorf("get file list: %w", err)
	}

	return body, nil
}

func (c *YSClient) GetDirectoryInfo(directoryNumber string) ([]byte, error) {
	return c.GetDirectoryInfoContext(context.Background(), directoryNumber)
}

func (c *YSClient) GetDirectoryInfoContext(ctx context.Context, directoryNumber string) ([]byte, error) {
	httpReq, err := c.newRequestWithContext(ctx, http.MethodGet, directoryPath+"?bjbh="+url.QueryEscape(directoryNumber), nil)
	if err != nil {
		return nil, err
	}

	var body []byte
	if err := c.do(httpReq, &body, http.StatusOK); err != nil {
		return nil, fmt.Errorf("get directory info: %w", err)
	}

	return body, nil
}

func (c *YSClient) GetDirectoryList() ([]byte, error) {
	return c.GetDirectoryListContext(context.Background())
}

func (c *YSClient) GetDirectoryListContext(ctx context.Context) ([]byte, error) {
	httpReq, err := c.newRequestWithContext(ctx, http.MethodGet, directoryPath, nil)
	if err != nil {
		return nil, err
	}

	var body []byte
	if err := c.do(httpReq, &body, http.StatusOK); err != nil {
		return nil, fmt.Errorf("get directory list: %w", err)
	}

	return body, nil
}

func (c *YSClient) SetDirectorySettings(operation string, req *DirectorySettingsRequest) error {
	return c.SetDirectorySettingsContext(context.Background(), operation, req)
}

func (c *YSClient) SetDirectorySettingsContext(ctx context.Context, operation string, req *DirectorySettingsRequest) error {
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

	httpReq, err := c.newRequestWithContext(ctx, http.MethodPost, directoryPath+"?cz="+operation, form)
	if err != nil {
		return err
	}

	if err := c.do(httpReq, nil, http.StatusOK); err != nil {
		return fmt.Errorf("set directory settings %s: %w", operation, err)
	}

	return nil
}

func (c *YSClient) AddDirectory(req *DirectorySettingsRequest) error {
	return c.AddDirectoryContext(context.Background(), req)
}

func (c *YSClient) AddDirectoryContext(ctx context.Context, req *DirectorySettingsRequest) error {
	return c.SetDirectorySettingsContext(ctx, "add", req)
}

func (c *YSClient) UpdateDirectory(req *DirectorySettingsRequest) error {
	return c.UpdateDirectoryContext(context.Background(), req)
}

func (c *YSClient) UpdateDirectoryContext(ctx context.Context, req *DirectorySettingsRequest) error {
	return c.SetDirectorySettingsContext(ctx, "add", req)
}

func (c *YSClient) DeleteDirectory(directoryNumber string) error {
	return c.DeleteDirectoryContext(context.Background(), directoryNumber)
}

func (c *YSClient) DeleteDirectoryContext(ctx context.Context, directoryNumber string) error {
	form := url.Values{}
	form.Add("mlbh", directoryNumber)

	httpReq, err := c.newRequestWithContext(ctx, http.MethodPost, directoryPath+"?cz=del", form)
	if err != nil {
		return err
	}

	if err := c.do(httpReq, nil, http.StatusOK); err != nil {
		return fmt.Errorf("delete directory: %w", err)
	}

	return nil
}

func (c *YSClient) DeleteFiles(req *DeleteFilesRequest) error {
	return c.DeleteFilesContext(context.Background(), req)
}

func (c *YSClient) DeleteFilesContext(ctx context.Context, req *DeleteFilesRequest) error {
	form := url.Values{}
	form.Add("mlbh", req.DirectoryNumber)
	form.Add("kqmm", req.OpenPassword)
	form.Add("wjs", joinCSV(req.FileNumbers))
	form.Add("xwjs", joinCSV(req.XFileNumbers))
	form.Add("links", joinCSV(req.LinkNumbers))
	form.Add("zmls", joinSubdirectories(req.Subdirectories))

	httpReq, err := c.newRequestWithContext(ctx, http.MethodPost, fileListPath+"?cz=del", form)
	if err != nil {
		return err
	}

	if err := c.do(httpReq, nil, http.StatusOK); err != nil {
		return fmt.Errorf("delete files: %w", err)
	}

	return nil
}
