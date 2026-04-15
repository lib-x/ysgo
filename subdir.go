package ysgo

import (
	"context"
	"fmt"
	"sort"
	"strings"
)

// SubdirectoryEntry describes a discovered logical subdirectory.
type SubdirectoryEntry struct {
	Name string
	Path string
}

// ListSubdirectories derives a unique one-level subdirectory listing from file entries.
func ListSubdirectories(files []RemoteFile, basePath string) []SubdirectoryEntry {
	prefix := strings.Trim(strings.TrimSpace(basePath), "/")
	seen := make(map[string]struct{})
	entries := make([]SubdirectoryEntry, 0)

	for _, file := range files {
		zml := strings.Trim(strings.TrimSpace(file.Subdirectory), "/")
		if zml == "" {
			continue
		}

		if prefix != "" {
			if zml == prefix {
				continue
			}
			wantPrefix := prefix + "/"
			if !strings.HasPrefix(zml, wantPrefix) {
				continue
			}
			zml = strings.TrimPrefix(zml, wantPrefix)
		}

		part := zml
		if idx := strings.Index(part, "/"); idx >= 0 {
			part = part[:idx]
		}
		if part == "" {
			continue
		}

		fullPath := part
		if prefix != "" {
			fullPath = prefix + "/" + part
		}
		if _, ok := seen[fullPath]; ok {
			continue
		}
		seen[fullPath] = struct{}{}
		entries = append(entries, SubdirectoryEntry{Name: part, Path: fullPath})
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Path < entries[j].Path
	})

	return entries
}

// GetSubdirectories returns unique one-level subdirectories rooted under basePath.
func (c *YSClient) GetSubdirectories(req *FileListRequest, basePath string) ([]SubdirectoryEntry, error) {
	return c.GetSubdirectoriesContext(context.Background(), req, basePath)
}

// GetSubdirectoriesContext returns unique one-level subdirectories rooted under basePath with context control.
func (c *YSClient) GetSubdirectoriesContext(ctx context.Context, req *FileListRequest, basePath string) ([]SubdirectoryEntry, error) {
	resp, err := c.GetFileListParsedContext(ctx, req)
	if err != nil {
		return nil, err
	}
	return ListSubdirectories(resp.Files, basePath), nil
}

// DeleteSubdirectory deletes a logical subdirectory and all entries beneath it.
func (c *YSClient) DeleteSubdirectory(directoryNumber, openPassword, path string) error {
	return c.DeleteSubdirectoryContext(context.Background(), directoryNumber, openPassword, path)
}

// DeleteSubdirectoryContext deletes a logical subdirectory and all entries beneath it with context control.
func (c *YSClient) DeleteSubdirectoryContext(ctx context.Context, directoryNumber, openPassword, path string) error {
	path = strings.Trim(strings.TrimSpace(path), "/")
	if path == "" {
		return fmt.Errorf("delete subdirectory: empty path")
	}
	return c.DeleteFilesContext(ctx, &DeleteFilesRequest{
		DirectoryNumber: directoryNumber,
		OpenPassword:    openPassword,
		Subdirectories:  []string{path},
	})
}
