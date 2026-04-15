# ysgo

永硕 E 盘 Go SDK。

## 当前能力

- 初始化会话 `InitSession`
- 管理员登录 `Login`
- 获取目录列表 `GetDirectoryList` / `GetDirectoryListParsed`
- 获取文件列表 `GetFileList` / `GetFileListParsed`
- 获取目录信息 `GetDirectoryInfo`
- 新增/更新/删除目录
- 删除文件/链接/子目录 `DeleteFiles`
- 获取上传凭证 `GetUploadToken`
- 上传文件 `UploadFile` / `UploadBytes`
- 下载链接生成 `BuildDownloadURL`
- 实际下载 `DownloadToWriter` / `DownloadBytes`
- 子目录解析 / 删除
- 文件条目新增 / 更新 / 移动 / 可见性切换 / 恢复删除
- 排序接口

## 安装

```bash
go get github.com/lib-x/ysgo
```


## 推荐配置方式

推荐优先使用 functional options 集中配置客户端：

```go
client := ysgo.NewClient(
    user,
    managePass,
    ysgo.WithSpacePassword("110119"),
    ysgo.WithTimeout(15*time.Second),
    ysgo.WithManagementDirectory("1445856"),
)
```

## 基本用法

```go
package main

import (
    "fmt"
    "log"

    "github.com/lib-x/ysgo"
)

func main() {
    client := ysgo.NewClient("your-endpoint", "your-manage-password")

    session, err := client.InitSession()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("upload addr:", session.Space.UploadAddress)

    files, err := client.GetFileListParsed(&ysgo.FileListRequest{
        DirectoryNumber: "1445856",
        FileNumber:      "0",
        IP1:             session.Space.IP1,
    })
    if err != nil {
        log.Fatal(err)
    }

    for _, f := range files.Files {
        fmt.Println(f.Number, f.FileName, f.Size)
    }
}
```

## 管理目录

```go
_, err := client.PrepareAdminSession()
if err != nil {
    log.Fatal(err)
}

err = client.AddDirectory(&ysgo.DirectorySettingsRequest{
    Number:       "0",
    Title:        "ysgo-sdk-sandbox",
    Description:  "sandbox only",
    OpenPassword: "",
    SortNumber:   "0",
    OpenMethod:   "0",
    FileSort:     "1",
    Permissions:  "000101",
    Time:         "",
    SortWeight:   "0",
})
```

## 上传文件

```go
_, err := client.PrepareSession()
if err != nil {
    log.Fatal(err)
}

result, err := client.UploadBytes("1445856", "", "", "hello.txt", []byte("hello ysgo"))
if err != nil {
    log.Fatal(err)
}

fmt.Println(result.FileNumber, result.FileToken, result.Server)
```

## 下载链接

```go
files, err := client.GetFileListParsed(&ysgo.FileListRequest{
    DirectoryNumber: "1445856",
    FileNumber:      "0",
})
if err != nil {
    log.Fatal(err)
}

if len(files.Files) > 0 {
    downloadURL, err := client.BuildDownloadURL(
        1445856,
        files.Directory.DownloadToken,
        files.Files[0],
        &ysgo.DownloadURLOptions{ForceDownload: true},
    )
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(downloadURL)
}
```

## 实际下载文件

```go
files, err := client.GetFileListParsed(&ysgo.FileListRequest{
    DirectoryNumber: "1445856",
    FileNumber:      "0",
})
if err != nil {
    log.Fatal(err)
}

u, err := client.BuildDownloadURL(1445856, files.Directory.DownloadToken, files.Files[0], nil)
if err != nil {
    log.Fatal(err)
}

body, err := client.DownloadBytes(u)
if err != nil {
    log.Fatal(err)
}
fmt.Println(len(body))
```

## 子目录

> 当前 SDK 通过文件列表解析子目录树；上传到一个新的 `Subdirectory` 时会自动形成对应子目录。

```go
_, err = client.UploadBytes("1445856", "", "nested/demo", "hello.txt", []byte("hello"))
if err != nil {
    log.Fatal(err)
}

subs, err := client.GetSubdirectories(&ysgo.FileListRequest{
    DirectoryNumber: "1445856",
    FileNumber:      "0",
}, "nested")
if err != nil {
    log.Fatal(err)
}
for _, s := range subs {
    fmt.Println(s.Path)
}

if err := client.DeleteSubdirectory("1445856", "", "nested/demo"); err != nil {
    log.Fatal(err)
}
```

## Context 与超时

SDK 默认使用 `30s` HTTP 超时，并为主要接口提供 `...Context` 版本。

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

resp, err := client.GetDirectoryListParsedContext(ctx)
if err != nil {
    log.Fatal(err)
}
fmt.Println(len(resp.List))
```

## Live 集成测试

仓库包含基于环境变量的真实集成测试骨架：

```bash
YSGO_LIVE_USER=your-user \
YSGO_LIVE_PASS=your-pass \
YSGO_LIVE_WRITE=1 \
YSGO_LIVE_SPACE_PASSWORD=110119 \
go test ./integration -v
```

注意：真实写操作测试应始终限制在 sandbox 目录内。

## 生产使用建议

- SDK 默认 HTTP 超时为 `30s`。
- 推荐优先使用 `...Context` 版本接口，将超时/取消控制交给调用方。
- 上传地址会校验为 `https`，且主机必须属于 `*.ysepan.com` 或 `*.ys168.com`；loopback 地址仅用于测试。
- 管理密码与 token 不应记录到日志中。
- 真实写操作请始终限制在 sandbox 目录内进行回归验证。

## 支持矩阵

| 能力 | 状态 | 真实验证 |
|---|---|---|
| InitSession / Login | 已支持 | 是 |
| 目录列表 / 文件列表 | 已支持 | 是 |
| 目录新增 / 更新 / 删除 | 已支持 | 是 |
| 文件上传 / 下载 / 删除 | 已支持 | 是 |
| 子目录解析 / 删除 | 已支持 | 是 |
| 排序接口 | 已支持 | 协议已验证，默认不改真实数据 |
| 文件移动 / 公开状态切换 / 恢复删除 | 已支持 | 文件移动与公开状态已在 sandbox 验证；恢复删除已有协议与测试骨架 |

## 推荐生命周期

```go
client := ysgo.NewClient(user, pass)

// 只读能力。
if _, err := client.PrepareSession(); err != nil {
    log.Fatal(err)
}

// 管理能力。
if _, err := client.PrepareAdminSession(); err != nil {
    log.Fatal(err)
}
```

## 已验证说明

真实 ysepan/ys168 sandbox 联调已覆盖：

- 目录新增 / 更新 / 删除
- 文件上传 / 列表 / 下载 / 删除 / 恢复
- 子目录创建 / 解析 / 删除
- 文件条目新增 / 子目录创建 / 条目移动 / 可见性切换

默认不直接对现有真实目录顺序做持久修改，因此排序相关目前以协议级自动化验证为主。

## 管理登录目录配置

如果你的管理登录需要显式指定目录号，可以覆盖默认目录号：

```go
client := ysgo.NewClient(user, pass, ysgo.WithManagementDirectory("1445856"))
```

## 并发说明

`YSClient` 适合被多个 goroutine 共享用于常规请求流程；内部会话 token 读取/更新已做同步保护。  
仍建议避免在外部无序地并发切换客户端配置对象本身。


## 空间访问密码

如果空间启用了访问密码，可以直接通过 functional option 提供：

```go
client := ysgo.NewClient(user, managePass, ysgo.WithSpacePassword("110119"))

_, err := client.PrepareSession()
if err != nil {
    log.Fatal(err)
}
```

SDK 在 `InitSession()` / `PrepareSession()` 遇到“需要输入登陆密码”时会自动尝试验证该空间访问密码并重新初始化会话。
