package ysgo

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
)

const uploadChunkSize int64 = 2 * 1024 * 1024

// FileChunkReader is the minimal random-access reader required for uploads.
type FileChunkReader interface {
	io.ReaderAt
}

// GetUploadToken retrieves the upload token and upload address for a directory.
func (c *YSClient) GetUploadToken(directoryNumber, openPassword string) (*UploadTokenResponse, error) {
	return c.GetUploadTokenContext(context.Background(), directoryNumber, openPassword)
}

// GetUploadTokenContext retrieves the upload token and upload address for a directory with context control.
func (c *YSClient) GetUploadTokenContext(ctx context.Context, directoryNumber, openPassword string) (*UploadTokenResponse, error) {
	form := url.Values{}
	form.Add("mlbh", directoryNumber)
	form.Add("kqmm", openPassword)
	form.Add("lx", "sc")

	httpReq, err := c.newRequestWithContext(ctx, http.MethodPost, directoryPath+"?cz=hqpz", form)
	if err != nil {
		return nil, err
	}

	var resp UploadTokenResponse
	if err := c.do(httpReq, &resp, http.StatusOK); err != nil {
		return nil, fmt.Errorf("get upload token: %w", err)
	}

	return &resp, nil
}

// GetDirectoryListParsed returns the typed directory listing response.
func (c *YSClient) GetDirectoryListParsed() (*DirectoryListResponse, error) {
	return c.GetDirectoryListParsedContext(context.Background())
}

// GetDirectoryListParsedContext returns the typed directory listing response with context control.
func (c *YSClient) GetDirectoryListParsedContext(ctx context.Context) (*DirectoryListResponse, error) {
	body, err := c.GetDirectoryListContext(ctx)
	if err != nil {
		return nil, err
	}

	var resp DirectoryListResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse directory list: %w", err)
	}

	return &resp, nil
}

// GetFileListParsed returns the typed file listing response.
func (c *YSClient) GetFileListParsed(req *FileListRequest) (*FileListResponse, error) {
	return c.GetFileListParsedContext(context.Background(), req)
}

// GetFileListParsedContext returns the typed file listing response with context control.
func (c *YSClient) GetFileListParsedContext(ctx context.Context, req *FileListRequest) (*FileListResponse, error) {
	body, err := c.GetFileListContext(ctx, req)
	if err != nil {
		return nil, err
	}

	var resp FileListResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse file list: %w", err)
	}

	return &resp, nil
}

// UploadBytes uploads an in-memory file payload to a directory.
func (c *YSClient) UploadBytes(directoryNumber, openPassword, subdirectory, fileName string, data []byte) (*UploadResult, error) {
	return c.UploadBytesContext(context.Background(), directoryNumber, openPassword, subdirectory, fileName, data)
}

// UploadBytesContext uploads an in-memory file payload to a directory with context control.
func (c *YSClient) UploadBytesContext(ctx context.Context, directoryNumber, openPassword, subdirectory, fileName string, data []byte) (*UploadResult, error) {
	return c.UploadFileContext(ctx, &UploadRequest{
		DirectoryNumber: directoryNumber,
		OpenPassword:    openPassword,
		Subdirectory:    subdirectory,
		FileName:        fileName,
		Reader:          bytes.NewReader(data),
		Size:            int64(len(data)),
		Public:          true,
	})
}

// UploadFile uploads a ReaderAt-backed file using the real ysepan chunked upload protocol.
func (c *YSClient) UploadFile(req *UploadRequest) (*UploadResult, error) {
	return c.UploadFileContext(context.Background(), req)
}

// UploadFileContext uploads a ReaderAt-backed file using the real ysepan chunked upload protocol with context control.
func (c *YSClient) UploadFileContext(ctx context.Context, req *UploadRequest) (*UploadResult, error) {
	if req == nil {
		return nil, fmt.Errorf("upload file: nil request")
	}
	if req.DirectoryNumber == "" {
		return nil, fmt.Errorf("upload file: empty directory number")
	}
	if req.FileName == "" {
		return nil, fmt.Errorf("upload file: empty file name")
	}
	if req.Reader == nil {
		return nil, fmt.Errorf("upload file: nil reader")
	}
	if req.Size <= 0 {
		return nil, fmt.Errorf("upload file: invalid size %d", req.Size)
	}
	if strings.ContainsAny(req.FileName, "%/\\'\"<>&:*?|#`") {
		return nil, fmt.Errorf("upload file: invalid file name %q", req.FileName)
	}

	tokenResp, err := c.GetUploadTokenContext(ctx, req.DirectoryNumber, req.OpenPassword)
	if err != nil {
		return nil, err
	}
	if tokenResp.Directory.UploadToken == "" {
		return nil, fmt.Errorf("upload file: empty upload token")
	}
	if tokenResp.Space.UploadAddress == "" {
		return nil, fmt.Errorf("upload file: empty upload address")
	}
	if err := validateUploadAddress(tokenResp.Space.UploadAddress, c.getAllowedUploadHosts()); err != nil {
		return nil, fmt.Errorf("upload file: %w", err)
	}

	totalChunks := int((req.Size + uploadChunkSize - 1) / uploadChunkSize)
	gid := ""
	publicFlag := "0"
	if req.Public {
		publicFlag = "1"
	}

	var result UploadResult
	for chunkIndex := 0; chunkIndex < totalChunks; chunkIndex++ {
		start := int64(chunkIndex) * uploadChunkSize
		end := start + uploadChunkSize
		if end > req.Size {
			end = req.Size
		}
		chunkSize := end - start
		chunkData := make([]byte, chunkSize)
		if _, err := req.Reader.ReadAt(chunkData, start); err != nil {
			return nil, fmt.Errorf("upload file: read chunk %d: %w", chunkIndex, err)
		}

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		if chunkIndex == 0 {
			_ = writer.WriteField("dlmc", c.userEndpoint)
			_ = writer.WriteField("scpz", tokenResp.Directory.UploadToken)
			_ = writer.WriteField("zmlmc", req.Subdirectory)
			_ = writer.WriteField("jsq", "0")
			_ = writer.WriteField("pdgk", publicFlag)
			_ = writer.WriteField("bz", "")
			_ = writer.WriteField("fsize", fmt.Sprintf("%d", req.Size))
			_ = writer.WriteField("isbm", "1")
		}
		_ = writer.WriteField("gid", gid)
		_ = writer.WriteField("chunkIndex", fmt.Sprintf("%d", chunkIndex))
		_ = writer.WriteField("totalfps", fmt.Sprintf("%d", totalChunks))

		part, err := writer.CreateFormFile("file", encodeUploadFileName(req.FileName))
		if err != nil {
			return nil, fmt.Errorf("upload file: create form file: %w", err)
		}
		if _, err := part.Write(chunkData); err != nil {
			return nil, fmt.Errorf("upload file: write chunk %d: %w", chunkIndex, err)
		}
		if err := writer.Close(); err != nil {
			return nil, fmt.Errorf("upload file: close multipart writer: %w", err)
		}

		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenResp.Space.UploadAddress, body)
		if err != nil {
			return nil, fmt.Errorf("upload file: create request: %w", err)
		}
		httpReq.Header.Set("Content-Type", writer.FormDataContentType())

		resp, err := c.getHTTPClient().Do(httpReq)
		if err != nil {
			return nil, fmt.Errorf("upload file: post chunk %d: %w", chunkIndex, err)
		}
		respBody, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			return nil, fmt.Errorf("upload file: read chunk response %d: %w", chunkIndex, readErr)
		}
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("upload file: chunk %d status %d: %s", chunkIndex, resp.StatusCode, strings.TrimSpace(string(respBody)))
		}

		var chunkResp struct {
			GID string `json:"gid"`
			UploadResult
		}
		if err := json.Unmarshal(respBody, &chunkResp); err != nil {
			return nil, fmt.Errorf("upload file: decode chunk response %d: %w", chunkIndex, err)
		}
		if chunkIndex == 0 {
			gid = chunkResp.GID
		}
		if chunkIndex == totalChunks-1 {
			result = chunkResp.UploadResult
		}
	}

	return &result, nil
}

func encodeUploadFileName(name string) string {
	base := filepath.Base(name)
	return base64.StdEncoding.EncodeToString([]byte(base))
}
