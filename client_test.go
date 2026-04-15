package ysgo

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestInitSessionUsesRandomTokenThenServerToken(t *testing.T) {
	var authHeaders []string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeaders = append(authHeaders, r.Header.Get("Authorization"))
		if r.URL.Path != "/nkj/csxx.aspx" || r.URL.RawQuery != "cz=dq" {
			t.Fatalf("unexpected path: %s?%s", r.URL.Path, r.URL.RawQuery)
		}
		_, _ = io.WriteString(w, `{"token":"server-token","yh":{"isgly":false},"kj":{"scdz":"https://upload.example.com"}}`)
	}))
	defer ts.Close()

	client := NewClient("czyt", "pass", WithAPIBaseURL(ts.URL))
	if client.GetAuthToken() == "" {
		t.Fatal("expected initial auth token")
	}
	if client.GetAuthToken() == "server-token" {
		t.Fatal("expected random initial token before init")
	}

	resp, err := client.InitSession()
	if err != nil {
		t.Fatalf("InitSession returned error: %v", err)
	}
	if resp.Token != "server-token" {
		t.Fatalf("expected server token, got %q", resp.Token)
	}
	if client.GetAuthToken() != "server-token" {
		t.Fatalf("expected client token to update, got %q", client.GetAuthToken())
	}
	if len(authHeaders) != 1 {
		t.Fatalf("expected 1 request, got %d", len(authHeaders))
	}
	if authHeaders[0] != "Bearer czyt;"+strings.TrimPrefix(authHeaders[0], "Bearer czyt;") {
		// structural guard only.
	}
	if !strings.HasPrefix(authHeaders[0], "Bearer czyt;") {
		t.Fatalf("unexpected auth header: %q", authHeaders[0])
	}
}

func TestGetFileListUsesSessionProtocol(t *testing.T) {
	var gotAuth string
	var gotForm url.Values

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/nkj/wj.aspx" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		gotAuth = r.Header.Get("Authorization")
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		gotForm, err = url.ParseQuery(string(body))
		if err != nil {
			t.Fatalf("parse body: %v", err)
		}
		_, _ = io.WriteString(w, `{"ok":true}`)
	}))
	defer ts.Close()

	client := NewClient("czyt", "pass", WithAPIBaseURL(ts.URL), WithAuthToken("fixed-token"))
	body, err := client.GetFileList(&FileListRequest{
		DirectoryNumber: "1445856",
		OpenPassword:    "",
		FileNumber:      "0",
		IP1:             "1.2.3.4",
	})
	if err != nil {
		t.Fatalf("GetFileList returned error: %v", err)
	}
	if string(body) != `{"ok":true}` {
		t.Fatalf("unexpected body: %s", string(body))
	}
	if gotAuth != "Bearer czyt;fixed-token" {
		t.Fatalf("unexpected auth header: %q", gotAuth)
	}
	if gotForm.Get("mlbh") != "1445856" || gotForm.Get("wjbh") != "0" || gotForm.Get("ip1") != "1.2.3.4" {
		t.Fatalf("unexpected form: %#v", gotForm)
	}
}

func TestLoginParsesJSONResponse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawQuery != "cz=yzglmm" {
			t.Fatalf("unexpected query: %s", r.URL.RawQuery)
		}
		_, _ = io.WriteString(w, `{"yh":{"isgly":true},"ml":{"xzpz":"d","xzpzsj":"now","bh":1,"scpz":"u"},"kj":{"scdz":"https://upload.example.com","jsq":2,"jsqj":3}}`)
	}))
	defer ts.Close()

	client := NewClient("czyt", "secret", WithAPIBaseURL(ts.URL), WithAuthToken("fixed-token"))
	resp, err := client.Login()
	if err != nil {
		t.Fatalf("Login returned error: %v", err)
	}
	if !resp.User.IsAdmin || resp.Directory.Number != 1 || resp.Space.Counter != 2 {
		t.Fatalf("unexpected response: %#v", resp)
	}
}

func TestLoginMapsCredentialError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(httpStatusAlreadyReported)
		_, _ = io.WriteString(w, "ERR\n需提供:glmm")
	}))
	defer ts.Close()

	client := NewClient("czyt", "bad", WithAPIBaseURL(ts.URL), WithAuthToken("fixed-token"))
	_, err := client.Login()
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestGetDirectoryInfoMapsAdminRequired(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(httpStatusAlreadyReported)
		_, _ = io.WriteString(w, "ERR\n管理员登陆后才可以编辑目录")
	}))
	defer ts.Close()

	client := NewClient("czyt", "bad", WithAPIBaseURL(ts.URL), WithAuthToken("fixed-token"))
	_, err := client.GetDirectoryInfo("1445856")
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, ErrAdminRequired) {
		t.Fatalf("expected ErrAdminRequired, got %v", err)
	}
}

func TestDeleteDirectoryUsesRealAction(t *testing.T) {
	var gotForm url.Values

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawQuery != "cz=del" {
			t.Fatalf("unexpected query: %s", r.URL.RawQuery)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		gotForm, err = url.ParseQuery(string(body))
		if err != nil {
			t.Fatalf("parse body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := NewClient("czyt", "secret", WithAPIBaseURL(ts.URL), WithAuthToken("fixed-token"))
	if err := client.DeleteDirectory("123456"); err != nil {
		t.Fatalf("DeleteDirectory returned error: %v", err)
	}
	if gotForm.Get("mlbh") != "123456" {
		t.Fatalf("unexpected form: %#v", gotForm)
	}
}

func TestCompatibilityAliases(t *testing.T) {
	client := NewClient("czyt", "secret", WithApiBaseUrl("https://example.com"))
	if got := client.GetApiBaseUrl(); got != "https://example.com" {
		t.Fatalf("unexpected api base url: %q", got)
	}
}

func TestDeleteFilesUsesRealAction(t *testing.T) {
	var gotForm url.Values

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/nkj/wj.aspx" || r.URL.RawQuery != "cz=del" {
			t.Fatalf("unexpected path: %s?%s", r.URL.Path, r.URL.RawQuery)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		gotForm, err = url.ParseQuery(string(body))
		if err != nil {
			t.Fatalf("parse body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := NewClient("czyt", "secret", WithAPIBaseURL(ts.URL), WithAuthToken("fixed-token"))
	err := client.DeleteFiles(&DeleteFilesRequest{
		DirectoryNumber: "1445856",
		OpenPassword:    "",
		FileNumbers:     []string{"1", "2"},
		XFileNumbers:    []string{"3"},
		LinkNumbers:     []string{"4"},
		Subdirectories:  []string{"a/b", "c"},
	})
	if err != nil {
		t.Fatalf("DeleteFiles returned error: %v", err)
	}

	if gotForm.Get("mlbh") != "1445856" {
		t.Fatalf("unexpected mlbh: %#v", gotForm)
	}
	if gotForm.Get("wjs") != "1,2" {
		t.Fatalf("unexpected wjs: %#v", gotForm)
	}
	if gotForm.Get("xwjs") != "3" {
		t.Fatalf("unexpected xwjs: %#v", gotForm)
	}
	if gotForm.Get("links") != "4" {
		t.Fatalf("unexpected links: %#v", gotForm)
	}
	if gotForm.Get("zmls") != "a/b//<c" {
		t.Fatalf("unexpected zmls: %#v", gotForm)
	}
}

func TestGetDirectoryListUsesRootDirectoryEndpoint(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/nkj/ml.aspx" || r.URL.RawQuery != "" {
			t.Fatalf("unexpected path: %s?%s", r.URL.Path, r.URL.RawQuery)
		}
		_, _ = io.WriteString(w, `{"lb":[{"bh":1,"bt":"sandbox"}],"mlpx":2}`)
	}))
	defer ts.Close()

	client := NewClient("czyt", "secret", WithAPIBaseURL(ts.URL), WithAuthToken("fixed-token"))
	body, err := client.GetDirectoryList()
	if err != nil {
		t.Fatalf("GetDirectoryList returned error: %v", err)
	}
	if string(body) != `{"lb":[{"bh":1,"bt":"sandbox"}],"mlpx":2}` {
		t.Fatalf("unexpected body: %s", string(body))
	}
}

func TestGetDirectoryListParsed(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, `{"lb":[{"bh":1,"bt":"sandbox","sm":"desc","qx":"000101","kqfs":0,"wjpx":"1","pxbh":2,"qxz":true,"qck":true,"qsc":false,"qgl":false}],"mlpx":2}`)
	}))
	defer ts.Close()

	client := NewClient("czyt", "secret", WithAPIBaseURL(ts.URL), WithAuthToken("fixed-token"))
	resp, err := client.GetDirectoryListParsed()
	if err != nil {
		t.Fatalf("GetDirectoryListParsed returned error: %v", err)
	}
	if len(resp.List) != 1 || resp.List[0].Title != "sandbox" || resp.SortMode != 2 {
		t.Fatalf("unexpected response: %#v", resp)
	}
}

func TestGetFileListParsed(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, `{"kj":{"jzxzxx":""},"ml":{"xzpz":"tok","kqfs":0,"yzpd":false},"lb":[{"bh":7,"wjm":"a.txt","bt":"A","zml":"","sj":"2026-01-01T00:00:00","dx":3,"fwq":"C","pz":"pt"}]}`)
	}))
	defer ts.Close()

	client := NewClient("czyt", "secret", WithAPIBaseURL(ts.URL), WithAuthToken("fixed-token"))
	resp, err := client.GetFileListParsed(&FileListRequest{DirectoryNumber: "1", FileNumber: "0"})
	if err != nil {
		t.Fatalf("GetFileListParsed returned error: %v", err)
	}
	if len(resp.Files) != 1 || resp.Files[0].FileName != "a.txt" || resp.Directory.DownloadToken != "tok" {
		t.Fatalf("unexpected response: %#v", resp)
	}
}

func TestUploadBytes(t *testing.T) {
	var uploadContentType string
	var uploadBody []byte
	var uploadRequests int

	uploadServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uploadRequests++
		uploadContentType = r.Header.Get("Content-Type")
		var err error
		uploadBody, err = io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read upload body: %v", err)
		}
		_, _ = io.WriteString(w, `{"gid":"gid-1","bh":7,"fwq":"C","rq":"260101","pz":"file-token","sj":"2026-01-01T00:00:00"}`)
	}))
	defer uploadServer.Close()

	apiServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/nkj/ml.aspx" && r.URL.RawQuery == "cz=hqpz" {
			_, _ = io.WriteString(w, `{"kj":{"scdz":"`+uploadServer.URL+`","jzscxx":""},"ml":{"mlbh":1,"scpz":"upload-token"}}`)
			return
		}
		t.Fatalf("unexpected API request: %s?%s", r.URL.Path, r.URL.RawQuery)
	}))
	defer apiServer.Close()

	client := NewClient("czyt", "secret", WithAPIBaseURL(apiServer.URL), WithAuthToken("fixed-token"))
	result, err := client.UploadBytes("1", "", "subdir", "hello.txt", []byte("hello world"))
	if err != nil {
		t.Fatalf("UploadBytes returned error: %v", err)
	}
	if result.FileNumber != 7 || result.Server != "C" || result.FileToken != "file-token" {
		t.Fatalf("unexpected upload result: %#v", result)
	}
	if uploadRequests != 1 {
		t.Fatalf("expected 1 upload request, got %d", uploadRequests)
	}
	if !strings.HasPrefix(uploadContentType, "multipart/form-data;") {
		t.Fatalf("unexpected content type: %q", uploadContentType)
	}
	bodyStr := string(uploadBody)
	if !strings.Contains(bodyStr, "name=\"dlmc\"") || !strings.Contains(bodyStr, "czyt") {
		t.Fatalf("missing dlmc field: %s", bodyStr)
	}
	if !strings.Contains(bodyStr, "name=\"scpz\"") || !strings.Contains(bodyStr, "upload-token") {
		t.Fatalf("missing scpz field: %s", bodyStr)
	}
	if !strings.Contains(bodyStr, "name=\"zmlmc\"") || !strings.Contains(bodyStr, "subdir") {
		t.Fatalf("missing zmlmc field: %s", bodyStr)
	}
	if !strings.Contains(bodyStr, "name=\"file\"; filename=\"aGVsbG8udHh0\"") {
		t.Fatalf("missing encoded filename: %s", bodyStr)
	}
	if !strings.Contains(bodyStr, "hello world") {
		t.Fatalf("missing file data: %s", bodyStr)
	}
}

func TestBuildDownloadURL(t *testing.T) {
	client := NewClient("czyt", "secret")
	url, err := client.BuildDownloadURL(1445856, "token123", RemoteFile{
		Server:    "D",
		FileToken: "file-token",
		FileName:  "hello.txt",
	}, &DownloadURLOptions{ForceDownload: true, TextTimestamp: "12345"})
	if err != nil {
		t.Fatalf("BuildDownloadURL returned error: %v", err)
	}
	want := "https://ys-d.ysepan.com/wap/czyt/_token123/file-token/hello.txt?12345&lx=xz"
	if url != want {
		t.Fatalf("unexpected url:\nwant: %s\ngot:  %s", want, url)
	}
}

func TestUpdateDirectorySortUsesRealAction(t *testing.T) {
	var gotForm url.Values
	var gotQuery string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.RawQuery
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		gotForm, err = url.ParseQuery(string(body))
		if err != nil {
			t.Fatalf("parse body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := NewClient("czyt", "secret", WithAPIBaseURL(ts.URL), WithAuthToken("fixed-token"))
	if err := client.UpdateDirectorySort("123", 999); err != nil {
		t.Fatalf("UpdateDirectorySort returned error: %v", err)
	}
	if gotQuery != "cz=xgpxbh" {
		t.Fatalf("unexpected query: %s", gotQuery)
	}
	if gotForm.Get("mlbh") != "123" || gotForm.Get("pxbh") != "999" {
		t.Fatalf("unexpected form: %#v", gotForm)
	}
}

func TestUpdateFileSortUsesRealAction(t *testing.T) {
	var gotForm url.Values
	var gotQuery string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.RawQuery
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		gotForm, err = url.ParseQuery(string(body))
		if err != nil {
			t.Fatalf("parse body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := NewClient("czyt", "secret", WithAPIBaseURL(ts.URL), WithAuthToken("fixed-token"))
	if err := client.UpdateFileSort("123", 456, 789, "pass"); err != nil {
		t.Fatalf("UpdateFileSort returned error: %v", err)
	}
	if gotQuery != "cz=wjxgpxbh" {
		t.Fatalf("unexpected query: %s", gotQuery)
	}
	if gotForm.Get("mlbh") != "123" || gotForm.Get("wjbh") != "456" || gotForm.Get("jsq") != "789" || gotForm.Get("kqmm") != "pass" {
		t.Fatalf("unexpected form: %#v", gotForm)
	}
}

func TestSetDirectorySortModeUsesRealAction(t *testing.T) {
	var gotForm url.Values
	var gotPath string
	var gotQuery string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotQuery = r.URL.RawQuery
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		gotForm, err = url.ParseQuery(string(body))
		if err != nil {
			t.Fatalf("parse body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := NewClient("czyt", "secret", WithAPIBaseURL(ts.URL), WithAuthToken("fixed-token"))
	if err := client.SetDirectorySortMode(directorySortDescending); err != nil {
		t.Fatalf("SetDirectorySortMode returned error: %v", err)
	}
	if gotPath != "/nkj/kjxg.aspx" || gotQuery != "cz=xgmlpx" {
		t.Fatalf("unexpected request: %s?%s", gotPath, gotQuery)
	}
	if gotForm.Get("mlpx") != "2" {
		t.Fatalf("unexpected form: %#v", gotForm)
	}
}

func TestListSubdirectories(t *testing.T) {
	files := []RemoteFile{
		{Subdirectory: "nested/a"},
		{Subdirectory: "nested/b"},
		{Subdirectory: "nested2"},
		{Subdirectory: "nested/a/deeper"},
	}
	got := ListSubdirectories(files, "")
	if len(got) != 2 || got[0].Path != "nested" || got[1].Path != "nested2" {
		t.Fatalf("unexpected root subdirs: %#v", got)
	}
	got = ListSubdirectories(files, "nested")
	if len(got) != 2 || got[0].Path != "nested/a" || got[1].Path != "nested/b" {
		t.Fatalf("unexpected nested subdirs: %#v", got)
	}
}

func TestDeleteSubdirectoryUsesDeleteFilesProtocol(t *testing.T) {
	var gotForm url.Values
	var gotQuery string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.RawQuery
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		gotForm, err = url.ParseQuery(string(body))
		if err != nil {
			t.Fatalf("parse body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := NewClient("czyt", "secret", WithAPIBaseURL(ts.URL), WithAuthToken("fixed-token"))
	if err := client.DeleteSubdirectory("123", "", "nested/path"); err != nil {
		t.Fatalf("DeleteSubdirectory returned error: %v", err)
	}
	if gotQuery != "cz=del" {
		t.Fatalf("unexpected query: %s", gotQuery)
	}
	if gotForm.Get("zmls") != "nested/path" {
		t.Fatalf("unexpected zmls: %#v", gotForm)
	}
}

func TestDownloadBytes(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, "hello")
	}))
	defer ts.Close()

	client := NewClient("czyt", "secret")
	body, err := client.DownloadBytes(ts.URL)
	if err != nil {
		t.Fatalf("DownloadBytes returned error: %v", err)
	}
	if string(body) != "hello" {
		t.Fatalf("unexpected body: %s", string(body))
	}
}

func TestDefaultClientTimeout(t *testing.T) {
	client := NewClient("czyt", "secret")
	if client.c.Timeout != defaultTimeout {
		t.Fatalf("expected default timeout %v, got %v", defaultTimeout, client.c.Timeout)
	}
}

func TestContextCancellationPropagates(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
	}))
	defer ts.Close()

	client := NewClient("czyt", "secret", WithAPIBaseURL(ts.URL), WithAuthToken("fixed-token"))
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := client.GetDirectoryListContext(ctx)
	if err == nil {
		t.Fatal("expected cancellation error")
	}
}

func TestMoveEntriesUsesRealActions(t *testing.T) {
	queries := make([]string, 0)
	forms := make([]url.Values, 0)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		queries = append(queries, r.URL.RawQuery)
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		form, err := url.ParseQuery(string(body))
		if err != nil {
			t.Fatalf("parse body: %v", err)
		}
		forms = append(forms, form)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := NewClient("czyt", "secret", WithAPIBaseURL(ts.URL), WithAuthToken("fixed-token"))
	err := client.MoveEntries(&MoveEntriesRequest{
		SourceDirectoryNumber: "1",
		TargetDirectoryNumber: "2",
		SourceOpenPassword:    "a",
		TargetOpenPassword:    "b",
		SourcePath:            "src",
		TargetPath:            "dst",
		FileNumbers:           []int{11, 12},
		Subdirectories:        []string{"child"},
	})
	if err != nil {
		t.Fatalf("MoveEntries returned error: %v", err)
	}
	if len(queries) != 2 || queries[0] != "cz=plzy" || queries[1] != "cz=xgzml" {
		t.Fatalf("unexpected queries: %#v", queries)
	}
	if forms[0].Get("wjs") != "11,12" || forms[0].Get("xzml") != "dst" {
		t.Fatalf("unexpected move files form: %#v", forms[0])
	}
	if forms[1].Get("yzml") != "src/child" || forms[1].Get("xzml") != "dst/child" {
		t.Fatalf("unexpected move subdir form: %#v", forms[1])
	}
}

func TestSetFileVisibilityUsesRealAction(t *testing.T) {
	var gotForm url.Values
	var gotQuery string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.RawQuery
		body, _ := io.ReadAll(r.Body)
		gotForm, _ = url.ParseQuery(string(body))
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()
	client := NewClient("czyt", "secret", WithAPIBaseURL(ts.URL), WithAuthToken("fixed-token"))
	if err := client.SetFileVisibility("1", "pw", 9, false); err != nil {
		t.Fatalf("SetFileVisibility returned error: %v", err)
	}
	if gotQuery != "cz=gggk" || gotForm.Get("wjs") != "9" || gotForm.Get("pdgk") != "0" {
		t.Fatalf("unexpected form/query: %s %#v", gotQuery, gotForm)
	}
}

func TestRestoreDeletedFileUsesRealAction(t *testing.T) {
	var gotForm url.Values
	var gotQuery string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.RawQuery
		body, _ := io.ReadAll(r.Body)
		gotForm, _ = url.ParseQuery(string(body))
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()
	client := NewClient("czyt", "secret", WithAPIBaseURL(ts.URL), WithAuthToken("fixed-token"))
	if err := client.RestoreDeletedFile("1", "pw", 9); err != nil {
		t.Fatalf("RestoreDeletedFile returned error: %v", err)
	}
	if gotQuery != "cz=hfdel" || gotForm.Get("wjbh") != "9" {
		t.Fatalf("unexpected form/query: %s %#v", gotQuery, gotForm)
	}
}

func TestValidateUploadAddress(t *testing.T) {
	cases := []struct {
		addr    string
		wantErr bool
	}{
		{"https://ys-D.ysepan.com/wap/nkjup.aspx", false},
		{"https://y.ys168.com:8000/wap/x", false},
		{"http://ys-D.ysepan.com/wap/nkjup.aspx", true},
		{"https://evil.example.com/upload", true},
	}
	for _, tc := range cases {
		err := validateUploadAddress(tc.addr, allowedUploadHosts)
		if (err != nil) != tc.wantErr {
			t.Fatalf("validateUploadAddress(%q) err=%v wantErr=%v", tc.addr, err, tc.wantErr)
		}
	}
}

func TestPrepareAdminSession(t *testing.T) {
	calls := make([]string, 0)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, r.URL.Path+"?"+r.URL.RawQuery)
		switch r.URL.RawQuery {
		case "cz=dq":
			_, _ = io.WriteString(w, `{"token":"server-token","yh":{"isgly":false},"kj":{"scdz":"https://upload.example.com"}}`)
		case "cz=yzglmm":
			_, _ = io.WriteString(w, `{"yh":{"isgly":true},"ml":{"xzpz":"d","xzpzsj":"now","bh":1,"scpz":"u"},"kj":{"scdz":"https://upload.example.com","jsq":2,"jsqj":3}}`)
		default:
			t.Fatalf("unexpected query: %s", r.URL.RawQuery)
		}
	}))
	defer ts.Close()

	client := NewClient("czyt", "secret", WithAPIBaseURL(ts.URL), WithAuthToken("fixed-token"))
	resp, err := client.PrepareAdminSession()
	if err != nil {
		t.Fatalf("PrepareAdminSession returned error: %v", err)
	}
	if !resp.User.IsAdmin {
		t.Fatalf("expected admin response: %#v", resp)
	}
	if len(calls) != 2 || calls[0] != "/nkj/csxx.aspx?cz=dq" || calls[1] != "/nkj/csxx.aspx?cz=yzglmm" {
		t.Fatalf("unexpected calls: %#v", calls)
	}
}

func TestAddEntryUsesRealAction(t *testing.T) {
	var gotForm url.Values
	var gotQuery string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.RawQuery
		body, _ := io.ReadAll(r.Body)
		gotForm, _ = url.ParseQuery(string(body))
		_, _ = io.WriteString(w, `{"bh":123,"sj":"2026-01-01T00:00:00"}`)
	}))
	defer ts.Close()

	client := NewClient("czyt", "secret", WithAPIBaseURL(ts.URL), WithAuthToken("fixed-token"))
	entry, err := client.AddEntry(&AddEntryRequest{
		DirectoryNumber: "1",
		OpenPassword:    "pw",
		Title:           "title",
		Content:         "https://example.com",
		Subdirectory:    "nested",
		Sequence:        77,
		Public:          true,
	})
	if err != nil {
		t.Fatalf("AddEntry returned error: %v", err)
	}
	if gotQuery != "cz=addlj" {
		t.Fatalf("unexpected query: %s", gotQuery)
	}
	if gotForm.Get("mlbh") != "1" || gotForm.Get("bt") != "title" || gotForm.Get("wjm") != "https://example.com" || gotForm.Get("zml") != "nested" {
		t.Fatalf("unexpected form: %#v", gotForm)
	}
	if entry.Number != 123 || entry.Title != "title" {
		t.Fatalf("unexpected entry: %#v", entry)
	}
}

func TestCreateSubdirectoryUsesRealAction(t *testing.T) {
	var gotForm url.Values
	var gotQuery string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.RawQuery
		body, _ := io.ReadAll(r.Body)
		gotForm, _ = url.ParseQuery(string(body))
		_, _ = io.WriteString(w, `{"bh":124,"sj":"2026-01-01T00:00:00"}`)
	}))
	defer ts.Close()

	client := NewClient("czyt", "secret", WithAPIBaseURL(ts.URL), WithAuthToken("fixed-token"))
	entry, err := client.CreateSubdirectory("1", "pw", "nested/live", 88)
	if err != nil {
		t.Fatalf("CreateSubdirectory returned error: %v", err)
	}
	if gotQuery != "cz=addlj" {
		t.Fatalf("unexpected query: %s", gotQuery)
	}
	if gotForm.Get("qlpd") != "1" || gotForm.Get("zml") != "nested/live" {
		t.Fatalf("unexpected form: %#v", gotForm)
	}
	if entry.Number != 124 || entry.Title != "live" {
		t.Fatalf("unexpected entry: %#v", entry)
	}
}

func TestUpdateEntryUsesRealAction(t *testing.T) {
	var gotForm url.Values
	var gotQuery string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.RawQuery
		body, _ := io.ReadAll(r.Body)
		gotForm, _ = url.ParseQuery(string(body))
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := NewClient("czyt", "secret", WithAPIBaseURL(ts.URL), WithAuthToken("fixed-token"))
	if err := client.UpdateEntry(&UpdateEntryRequest{
		DirectoryNumber: "1",
		OpenPassword:    "pw",
		FileNumber:      9,
		Kind:            EntryKindLink,
		Title:           "new-title",
		Content:         "https://new.example.com",
		Subdirectory:    "nested",
		Public:          true,
	}); err != nil {
		t.Fatalf("UpdateEntry returned error: %v", err)
	}
	if gotQuery != "cz=xgwj" || gotForm.Get("wjbh") != "9" || gotForm.Get("lx") != "l" {
		t.Fatalf("unexpected form/query: %s %#v", gotQuery, gotForm)
	}
}

func TestLoginUsesConfiguredManagementDirectory(t *testing.T) {
	var gotForm url.Values
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		gotForm, _ = url.ParseQuery(string(body))
		_, _ = io.WriteString(w, `{"yh":{"isgly":true},"ml":{"xzpz":"d","xzpzsj":"now","bh":1,"scpz":"u"},"kj":{"scdz":"https://upload.example.com","jsq":2,"jsqj":3}}`)
	}))
	defer ts.Close()

	client := NewClient("czyt", "secret", WithAPIBaseURL(ts.URL), WithAuthToken("fixed-token"), WithManagementDirectory("999"))
	if _, err := client.Login(); err != nil {
		t.Fatalf("Login returned error: %v", err)
	}
	if gotForm.Get("mlbh") != "999" {
		t.Fatalf("unexpected management directory: %#v", gotForm)
	}
}

func TestInitSessionAutoVerifiesSpacePassword(t *testing.T) {
	calls := make([]string, 0)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, r.URL.RawQuery)
		switch r.URL.RawQuery {
		case "cz=dq":
			count := 0
			for _, c := range calls {
				if c == "cz=dq" {
					count++
				}
			}
			if count == 1 {
				_, _ = io.WriteString(w, `{"fwxz":"[需要输入登陆密码].","token":"server-token-1"}`)
			} else {
				_, _ = io.WriteString(w, `{"yh":{"isgly":false},"kj":{"scdz":"https://upload.example.com"},"token":"server-token-2"}`)
			}
		case "cz=yzdlmm":
			_, _ = io.WriteString(w, `server-token-verified`)
		default:
			t.Fatalf("unexpected query: %s", r.URL.RawQuery)
		}
	}))
	defer ts.Close()

	client := NewClient("czyt", "secret", WithAPIBaseURL(ts.URL), WithSpacePassword("110119"))
	resp, err := client.InitSession()
	if err != nil {
		t.Fatalf("InitSession returned error: %v", err)
	}
	if resp.Token != "server-token-2" {
		t.Fatalf("unexpected final token: %#v", resp)
	}
	if len(calls) != 3 || calls[0] != "cz=dq" || calls[1] != "cz=yzdlmm" || calls[2] != "cz=dq" {
		t.Fatalf("unexpected call sequence: %#v", calls)
	}
}

func TestInitSessionRequiresSpacePasswordWithoutOption(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, `{"fwxz":"[需要输入登陆密码].","token":"server-token-1"}`)
	}))
	defer ts.Close()

	client := NewClient("czyt", "secret", WithAPIBaseURL(ts.URL))
	_, err := client.InitSession()
	if err == nil {
		t.Fatal("expected space password error")
	}
	if !errors.Is(err, ErrSpacePasswordRequired) {
		t.Fatalf("expected ErrSpacePasswordRequired, got %v", err)
	}
}

func TestVerifySpacePasswordInvalid(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, "ERR\n登陆密码不正确")
	}))
	defer ts.Close()

	client := NewClient("czyt", "secret", WithAPIBaseURL(ts.URL))
	err := client.VerifySpacePassword("bad")
	if err == nil {
		t.Fatal("expected invalid space password error")
	}
	if !errors.Is(err, ErrInvalidSpacePassword) {
		t.Fatalf("expected ErrInvalidSpacePassword, got %v", err)
	}
}

func TestWithTimeout(t *testing.T) {
	client := NewClient("czyt", "secret", WithTimeout(5*time.Second))
	if client.getHTTPClient().Timeout != 5*time.Second {
		t.Fatalf("unexpected timeout: %v", client.getHTTPClient().Timeout)
	}
}
