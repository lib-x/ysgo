package ysgo

import "time"

type LoginRequest struct {
	ManagementPassword string `json:"glmm" form:"glmm"`
	DirectoryNumber    string `json:"mlbh" form:"mlbh"`
}

type LoginResponse struct {
	User      LoginUserInfo      `json:"yh"`
	Directory LoginDirectoryInfo `json:"ml"`
	Space     LoginSpaceInfo     `json:"kj"`
}

type SessionResponse struct {
	User    LoginUserInfo  `json:"yh"`
	Space   LoginSpaceInfo `json:"kj"`
	Token   string         `json:"token,omitempty"`
	Message string         `json:"fwxz,omitempty"`
	Record  string         `json:"jldl,omitempty"`
}

type LoginUserInfo struct {
	IsAdmin bool `json:"isgly"`
}

type LoginDirectoryInfo struct {
	DownloadToken     string `json:"xzpz"`
	DownloadTokenTime string `json:"xzpzsj"`
	Number            int    `json:"bh"`
	UploadToken       string `json:"scpz"`
}

type LoginSpaceInfo struct {
	UploadAddress string `json:"scdz"`
	Counter       int    `json:"jsq"`
	CounterStage  int    `json:"jsqj"`
	IP0           string `json:"ip0,omitempty"`
	IP1           string `json:"ip1,omitempty"`
}

type PeriodicCheckRequest struct {
	DirectoryNumber string `json:"mlbh" form:"mlbh"`
	OpenPassword    string `json:"kqmm" form:"kqmm"`
	FileNumber      string `json:"wjbh" form:"wjbh"`
	UpdateModTime   string `json:"gxxmsj" form:"gxxmsj"`
}

type FileListRequest struct {
	DirectoryNumber string `json:"mlbh" form:"mlbh"`
	OpenPassword    string `json:"kqmm" form:"kqmm"`
	FileNumber      string `json:"wjbh" form:"wjbh"`
	IP1             string `json:"ip1,omitempty" form:"ip1,omitempty"`
}

type DirectorySettingsRequest struct {
	Number       string `json:"bh" form:"bh"`
	Title        string `json:"bt" form:"bt"`
	Description  string `json:"sm" form:"sm"`
	OpenPassword string `json:"kqmm" form:"kqmm"`
	SortNumber   string `json:"pxbh" form:"pxbh"`
	OpenMethod   string `json:"kqfs" form:"kqfs"`
	FileSort     string `json:"wjpx" form:"wjpx"`
	Permissions  string `json:"qx" form:"qx"`
	Time         string `json:"sj" form:"sj"`
	SortWeight   string `json:"pxz" form:"pxz"`
}

type DeleteFilesRequest struct {
	DirectoryNumber string
	OpenPassword    string
	FileNumbers     []string
	XFileNumbers    []string
	LinkNumbers     []string
	Subdirectories  []string
}

type FileInfo struct {
	Name        string    `json:"name"`
	Size        int64     `json:"size"`
	ModTime     time.Time `json:"mod_time"`
	IsDirectory bool      `json:"is_dir"`
	DownloadURL string    `json:"download_url,omitempty"`
	FileNumber  string    `json:"file_number,omitempty"`
}

type DirectoryInfo struct {
	Number      string     `json:"number"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Files       []FileInfo `json:"files,omitempty"`
}

type AuthToken struct {
	Username string
	Token    string
}

func (t AuthToken) String() string {
	return "Bearer " + t.Username + ";" + t.Token
}
