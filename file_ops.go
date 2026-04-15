package ysgo

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
)

// EntryKind identifies the remote entry kind used by file management APIs.
type EntryKind string

const (
	EntryKindFile EntryKind = "f"
	EntryKindLink EntryKind = "l"
	EntryKindText EntryKind = "w"
)

// AddEntryRequest describes a text/link entry creation request.
type AddEntryRequest struct {
	DirectoryNumber string
	OpenPassword    string
	Title           string
	Content         string
	Subdirectory    string
	Sequence        int
	Public          bool
}

// UpdateEntryRequest describes a file/link/text entry update request.
type UpdateEntryRequest struct {
	DirectoryNumber string
	OpenPassword    string
	FileNumber      int
	Kind            EntryKind
	Title           string
	Content         string
	Subdirectory    string
	Public          bool
}

// MoveEntriesRequest describes a cross-directory or cross-subdirectory move operation.
type MoveEntriesRequest struct {
	SourceDirectoryNumber string
	TargetDirectoryNumber string
	SourceOpenPassword    string
	TargetOpenPassword    string
	SourcePath            string
	TargetPath            string
	FileNumbers           []int
	Subdirectories        []string
}

// AddEntry creates a text or link entry in a directory.
func (c *YSClient) AddEntry(req *AddEntryRequest) (*RemoteFile, error) {
	return c.AddEntryContext(context.Background(), req)
}

// AddEntryContext creates a text or link entry in a directory with context control.
func (c *YSClient) AddEntryContext(ctx context.Context, req *AddEntryRequest) (*RemoteFile, error) {
	if req == nil {
		return nil, fmt.Errorf("add entry: nil request")
	}
	if req.DirectoryNumber == "" {
		return nil, fmt.Errorf("add entry: empty directory number")
	}
	if req.Title == "" && req.Content == "" {
		return nil, fmt.Errorf("add entry: title and content cannot both be empty")
	}

	form := url.Values{}
	form.Add("mlbh", req.DirectoryNumber)
	form.Add("bt", req.Title)
	form.Add("wjm", req.Content)
	form.Add("kqmm", req.OpenPassword)
	form.Add("pdgk", boolToIntString(req.Public))
	form.Add("jsq", fmt.Sprintf("%d", req.Sequence))
	form.Add("zml", req.Subdirectory)
	form.Add("qlpd", "0")

	httpReq, err := c.newRequestWithContext(ctx, http.MethodPost, fileListPath+"?cz=addlj", form)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Number int    `json:"bh"`
		Time   string `json:"sj"`
	}
	if err := c.do(httpReq, &resp, http.StatusOK); err != nil {
		return nil, fmt.Errorf("add entry: %w", err)
	}

	entry := &RemoteFile{
		Number:       resp.Number,
		FileName:     req.Content,
		Title:        req.Title,
		Subdirectory: req.Subdirectory,
		Time:         resp.Time,
		Visible:      req.Public,
		Sequence:     req.Sequence,
	}
	return entry, nil
}

// CreateSubdirectory creates a logical subdirectory marker entry inside a directory.
func (c *YSClient) CreateSubdirectory(directoryNumber, openPassword, path string, sequence int) (*RemoteFile, error) {
	return c.CreateSubdirectoryContext(context.Background(), directoryNumber, openPassword, path, sequence)
}

// CreateSubdirectoryContext creates a logical subdirectory marker entry inside a directory with context control.
func (c *YSClient) CreateSubdirectoryContext(ctx context.Context, directoryNumber, openPassword, path string, sequence int) (*RemoteFile, error) {
	path = strings.Trim(strings.TrimSpace(path), "/")
	if path == "" {
		return nil, fmt.Errorf("create subdirectory: empty path")
	}

	form := url.Values{}
	form.Add("mlbh", directoryNumber)
	form.Add("bt", "")
	form.Add("wjm", "")
	form.Add("kqmm", openPassword)
	form.Add("pdgk", "1")
	form.Add("jsq", fmt.Sprintf("%d", sequence))
	form.Add("zml", path)
	form.Add("qlpd", "1")

	httpReq, err := c.newRequestWithContext(ctx, http.MethodPost, fileListPath+"?cz=addlj", form)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Number int    `json:"bh"`
		Time   string `json:"sj"`
	}
	if err := c.do(httpReq, &resp, http.StatusOK); err != nil {
		return nil, fmt.Errorf("create subdirectory: %w", err)
	}

	name := path
	if idx := strings.LastIndex(path, "/"); idx >= 0 {
		name = path[idx+1:]
	}
	return &RemoteFile{
		Number:       resp.Number,
		Title:        name,
		FileName:     "",
		Subdirectory: path,
		Time:         resp.Time,
		Visible:      true,
		CanPreview:   true,
		Sequence:     sequence,
	}, nil
}

// UpdateEntry updates an existing file, link, or text entry.
func (c *YSClient) UpdateEntry(req *UpdateEntryRequest) error {
	return c.UpdateEntryContext(context.Background(), req)
}

// UpdateEntryContext updates an existing file, link, or text entry with context control.
func (c *YSClient) UpdateEntryContext(ctx context.Context, req *UpdateEntryRequest) error {
	if req == nil {
		return fmt.Errorf("update entry: nil request")
	}
	if req.DirectoryNumber == "" {
		return fmt.Errorf("update entry: empty directory number")
	}
	if req.FileNumber <= 0 {
		return fmt.Errorf("update entry: invalid file number %d", req.FileNumber)
	}

	kind := req.Kind
	if kind == "" {
		kind = detectEntryKind(req.Content)
	}
	if kind == EntryKindFile {
		ext := strings.ToLower(filepath.Ext(req.Content))
		if ext == "" {
			return fmt.Errorf("update entry: file content must include file name with extension")
		}
	}

	form := url.Values{}
	form.Add("lx", string(kind))
	form.Add("ymlbh", req.DirectoryNumber)
	form.Add("xmlbh", req.DirectoryNumber)
	form.Add("wjbh", fmt.Sprintf("%d", req.FileNumber))
	form.Add("wjm", req.Content)
	form.Add("bt", req.Title)
	form.Add("zml", req.Subdirectory)
	form.Add("pdgk", boolToIntString(req.Public))
	form.Add("kqmm", req.OpenPassword)

	httpReq, err := c.newRequestWithContext(ctx, http.MethodPost, fileListPath+"?cz=xgwj", form)
	if err != nil {
		return err
	}

	if err := c.do(httpReq, nil, http.StatusOK); err != nil {
		return fmt.Errorf("update entry: %w", err)
	}
	return nil
}

// SetFileVisibility toggles visitor visibility for a file entry.
func (c *YSClient) SetFileVisibility(directoryNumber, openPassword string, fileNumber int, visible bool) error {
	return c.SetFileVisibilityContext(context.Background(), directoryNumber, openPassword, fileNumber, visible)
}

// SetFileVisibilityContext toggles visitor visibility for a file entry with context control.
func (c *YSClient) SetFileVisibilityContext(ctx context.Context, directoryNumber, openPassword string, fileNumber int, visible bool) error {
	form := url.Values{}
	form.Add("mlbh", directoryNumber)
	form.Add("wjs", fmt.Sprintf("%d", fileNumber))
	form.Add("pdgk", boolToIntString(visible))
	form.Add("kqmm", openPassword)

	httpReq, err := c.newRequestWithContext(ctx, http.MethodPost, fileListPath+"?cz=gggk", form)
	if err != nil {
		return err
	}

	if err := c.do(httpReq, nil, http.StatusOK); err != nil {
		return fmt.Errorf("set file visibility: %w", err)
	}
	return nil
}

// RestoreDeletedFile restores a previously deleted file from the remote recycle bin.
func (c *YSClient) RestoreDeletedFile(directoryNumber, openPassword string, fileNumber int) error {
	return c.RestoreDeletedFileContext(context.Background(), directoryNumber, openPassword, fileNumber)
}

// RestoreDeletedFileContext restores a previously deleted file from the remote recycle bin with context control.
func (c *YSClient) RestoreDeletedFileContext(ctx context.Context, directoryNumber, openPassword string, fileNumber int) error {
	form := url.Values{}
	form.Add("mlbh", directoryNumber)
	form.Add("kqmm", openPassword)
	form.Add("wjbh", fmt.Sprintf("%d", fileNumber))

	httpReq, err := c.newRequestWithContext(ctx, http.MethodPost, fileListPath+"?cz=hfdel", form)
	if err != nil {
		return err
	}

	if err := c.do(httpReq, nil, http.StatusOK); err != nil {
		return fmt.Errorf("restore deleted file: %w", err)
	}
	return nil
}

// MoveEntries moves file entries and/or subdirectories between paths or directories.
func (c *YSClient) MoveEntries(req *MoveEntriesRequest) error {
	return c.MoveEntriesContext(context.Background(), req)
}

// MoveEntriesContext moves file entries and/or subdirectories between paths or directories with context control.
func (c *YSClient) MoveEntriesContext(ctx context.Context, req *MoveEntriesRequest) error {
	if req == nil {
		return fmt.Errorf("move entries: nil request")
	}
	if req.SourceDirectoryNumber == "" || req.TargetDirectoryNumber == "" {
		return fmt.Errorf("move entries: source/target directory number required")
	}
	if len(req.FileNumbers) == 0 && len(req.Subdirectories) == 0 {
		return fmt.Errorf("move entries: nothing to move")
	}

	if len(req.FileNumbers) > 0 {
		form := url.Values{}
		form.Add("ymlbh", req.SourceDirectoryNumber)
		form.Add("xmlbh", req.TargetDirectoryNumber)
		form.Add("xzml", req.TargetPath)
		form.Add("wjs", joinIntCSV(req.FileNumbers))
		form.Add("kqmm0", req.SourceOpenPassword)
		form.Add("kqmm1", req.TargetOpenPassword)

		httpReq, err := c.newRequestWithContext(ctx, http.MethodPost, fileListPath+"?cz=plzy", form)
		if err != nil {
			return err
		}
		if err := c.do(httpReq, nil, http.StatusOK); err != nil {
			return fmt.Errorf("move entries files: %w", err)
		}
	}

	for _, subdir := range req.Subdirectories {
		sourcePath := joinPath(req.SourcePath, subdir)
		if req.SourcePath == "" {
			sourcePath = strings.Trim(strings.TrimSpace(subdir), "/")
		}
		form := url.Values{}
		form.Add("ymlbh", req.SourceDirectoryNumber)
		form.Add("xmlbh", req.TargetDirectoryNumber)
		form.Add("yzml", sourcePath)
		form.Add("xzml", joinPath(req.TargetPath, subdir))
		form.Add("kqmm0", req.SourceOpenPassword)
		form.Add("kqmm1", req.TargetOpenPassword)

		httpReq, err := c.newRequestWithContext(ctx, http.MethodPost, fileListPath+"?cz=xgzml", form)
		if err != nil {
			return err
		}
		if err := c.do(httpReq, nil, http.StatusOK); err != nil {
			return fmt.Errorf("move entries subdirectory %s: %w", subdir, err)
		}
	}

	return nil
}

// GetFilesInSubdirectory returns file entries that belong exactly to path.
func (c *YSClient) GetFilesInSubdirectory(req *FileListRequest, path string) ([]RemoteFile, error) {
	return c.GetFilesInSubdirectoryContext(context.Background(), req, path)
}

// GetFilesInSubdirectoryContext returns file entries that belong exactly to path with context control.
func (c *YSClient) GetFilesInSubdirectoryContext(ctx context.Context, req *FileListRequest, path string) ([]RemoteFile, error) {
	resp, err := c.GetFileListParsedContext(ctx, req)
	if err != nil {
		return nil, err
	}
	path = strings.Trim(strings.TrimSpace(path), "/")
	files := make([]RemoteFile, 0)
	for _, file := range resp.Files {
		if strings.Trim(strings.TrimSpace(file.Subdirectory), "/") == path {
			files = append(files, file)
		}
	}
	return files, nil
}

func boolToIntString(v bool) string {
	if v {
		return "1"
	}
	return "0"
}

func detectEntryKind(content string) EntryKind {
	if looksLikeURL(content) {
		return EntryKindLink
	}
	return EntryKindText
}

func looksLikeURL(content string) bool {
	content = strings.TrimSpace(strings.ToLower(content))
	return strings.HasPrefix(content, "http://") || strings.HasPrefix(content, "https://") || strings.HasPrefix(content, "ftp://")
}

func joinIntCSV(values []int) string {
	if len(values) == 0 {
		return ""
	}
	parts := make([]string, 0, len(values))
	for _, v := range values {
		parts = append(parts, fmt.Sprintf("%d", v))
	}
	return strings.Join(parts, ",")
}

func joinPath(parent, child string) string {
	parent = strings.Trim(strings.TrimSpace(parent), "/")
	child = strings.Trim(strings.TrimSpace(child), "/")
	if parent == "" {
		return child
	}
	if child == "" {
		return parent
	}
	return parent + "/" + child
}
