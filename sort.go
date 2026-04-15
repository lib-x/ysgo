package ysgo

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

const (
	directorySortAscending  = 1
	directorySortDescending = 2
	sortSettingsPath        = "/nkj/kjxg.aspx"
)

// SetDirectorySortMode sets the global remote directory ordering mode.
func (c *YSClient) SetDirectorySortMode(mode int) error {
	return c.SetDirectorySortModeContext(context.Background(), mode)
}

// SetDirectorySortModeContext sets the global remote directory ordering mode with context control.
func (c *YSClient) SetDirectorySortModeContext(ctx context.Context, mode int) error {
	if mode != directorySortAscending && mode != directorySortDescending {
		return fmt.Errorf("set directory sort mode: invalid mode %d", mode)
	}

	form := url.Values{}
	form.Add("mlpx", fmt.Sprintf("%d", mode))

	httpReq, err := c.newRequestWithContext(ctx, http.MethodPost, sortSettingsPath+"?cz=xgmlpx", form)
	if err != nil {
		return err
	}

	if err := c.do(httpReq, nil, http.StatusOK); err != nil {
		return fmt.Errorf("set directory sort mode: %w", err)
	}

	return nil
}

// UpdateDirectorySort sets the sort number for a specific directory entry.
func (c *YSClient) UpdateDirectorySort(directoryNumber string, sortNumber int) error {
	return c.UpdateDirectorySortContext(context.Background(), directoryNumber, sortNumber)
}

// UpdateDirectorySortContext sets the sort number for a specific directory entry with context control.
func (c *YSClient) UpdateDirectorySortContext(ctx context.Context, directoryNumber string, sortNumber int) error {
	form := url.Values{}
	form.Add("mlbh", directoryNumber)
	form.Add("pxbh", fmt.Sprintf("%d", sortNumber))

	httpReq, err := c.newRequestWithContext(ctx, http.MethodPost, directoryPath+"?cz=xgpxbh", form)
	if err != nil {
		return err
	}

	if err := c.do(httpReq, nil, http.StatusOK); err != nil {
		return fmt.Errorf("update directory sort: %w", err)
	}

	return nil
}

// UpdateFileSort sets the sort sequence for a specific file entry.
func (c *YSClient) UpdateFileSort(directoryNumber string, fileNumber int, sequence int, openPassword string) error {
	return c.UpdateFileSortContext(context.Background(), directoryNumber, fileNumber, sequence, openPassword)
}

// UpdateFileSortContext sets the sort sequence for a specific file entry with context control.
func (c *YSClient) UpdateFileSortContext(ctx context.Context, directoryNumber string, fileNumber int, sequence int, openPassword string) error {
	form := url.Values{}
	form.Add("mlbh", directoryNumber)
	form.Add("wjbh", fmt.Sprintf("%d", fileNumber))
	form.Add("jsq", fmt.Sprintf("%d", sequence))
	form.Add("kqmm", openPassword)

	httpReq, err := c.newRequestWithContext(ctx, http.MethodPost, directoryPath+"?cz=wjxgpxbh", form)
	if err != nil {
		return err
	}

	if err := c.do(httpReq, nil, http.StatusOK); err != nil {
		return fmt.Errorf("update file sort: %w", err)
	}

	return nil
}
