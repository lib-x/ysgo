package ysgo

type DirectoryListResponse struct {
	List     []DirectoryEntry `json:"lb"`
	SortMode int              `json:"mlpx"`
}

type DirectoryEntry struct {
	Number        int    `json:"bh"`
	Title         string `json:"bt"`
	Description   string `json:"sm"`
	Permissions   string `json:"qx"`
	OpenMethod    int    `json:"kqfs"`
	FileSort      string `json:"wjpx"`
	SortNumber    int    `json:"pxbh"`
	CanDownload   bool   `json:"qxz"`
	CanView       bool   `json:"qck"`
	CanUpload     bool   `json:"qsc"`
	CanManage     bool   `json:"qgl"`
	UploadToken   string `json:"scpz,omitempty"`
	DownloadToken string `json:"xzpz,omitempty"`
	Verified      bool   `json:"yzpd,omitempty"`
	NeedsSyncTime bool   `json:"needgxsj,omitempty"`
	Time          string `json:"sj,omitempty"`
}

type FileListResponse struct {
	Space     FileListSpaceInfo     `json:"kj"`
	Directory FileListDirectoryInfo `json:"ml"`
	Files     []RemoteFile          `json:"lb"`
}

type FileListSpaceInfo struct {
	DownloadLimitMessage string `json:"jzxzxx"`
}

type FileListDirectoryInfo struct {
	DownloadToken string `json:"xzpz,omitempty"`
	UploadToken   string `json:"scpz,omitempty"`
	OpenMethod    int    `json:"kqfs,omitempty"`
	Verified      bool   `json:"yzpd,omitempty"`
}

type RemoteFile struct {
	Number         int    `json:"bh"`
	FileName       string `json:"wjm"`
	Title          string `json:"bt"`
	Subdirectory   string `json:"zml"`
	Time           string `json:"sj"`
	Size           int64  `json:"dx,omitempty"`
	Server         string `json:"fwq,omitempty"`
	FileToken      string `json:"pz,omitempty"`
	DateToken      string `json:"rq,omitempty"`
	DownloadCount  int    `json:"dayjsq,omitempty"`
	Visible        bool   `json:"pdgk,omitempty"`
	CanPreview     bool   `json:"qlpd,omitempty"`
	Sequence       int    `json:"jsq,omitempty"`
	IsImagePreview bool   `json:"xt,omitempty"`
	IsDeleted      bool   `json:"isdel,omitempty"`
	IsAdded        bool   `json:"isadd,omitempty"`
}

type UploadTokenResponse struct {
	Space struct {
		UploadAddress string `json:"scdz"`
		UploadMessage string `json:"jzscxx"`
	} `json:"kj"`
	Directory struct {
		Number      int    `json:"mlbh"`
		UploadToken string `json:"scpz"`
	} `json:"ml"`
}

type UploadRequest struct {
	DirectoryNumber string
	OpenPassword    string
	Subdirectory    string
	FileName        string
	Reader          FileChunkReader
	Size            int64
	Public          bool
}

type UploadResult struct {
	FileNumber int    `json:"bh"`
	Server     string `json:"fwq"`
	DateToken  string `json:"rq"`
	FileToken  string `json:"pz"`
	Time       string `json:"sj"`
	Size       int64  `json:"dx,omitempty"`
	HC         int    `json:"hc,omitempty"`
	XT         int    `json:"xt,omitempty"`
}
