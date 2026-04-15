# ysgo v0.1.0 发布说明

`ysgo v0.1.0` 是项目的首个**生产可用**版本。

本版本重点不是“接口数量”，而是把永硕 E 盘 / ysepan 的真实协议链路补齐、验证并稳定下来，使其能够被 Go 项目直接作为依赖接入，而不是停留在实验性质的协议脚本层面。

## 本版本重点能力

- 真实会话初始化与管理员登录。
- 空间访问密码支持。
- 目录与文件列表的结构化 API。
- 文件上传、下载、删除、恢复。
- 子目录解析、创建、删除。
- 文件条目新增、更新、移动、可见性切换。
- 排序相关接口。
- `Context` + 默认超时支持。
- functional options 风格客户端配置。
- 上传地址安全校验。
- 并发安全的会话状态管理。

## 推荐初始化方式

```go
client := ysgo.NewClient(
    user,
    managePass,
    ysgo.WithSpacePassword("<space-password>"),
    ysgo.WithTimeout(15*time.Second),
    ysgo.WithManagementDirectory("1445856"),
)
```

## 真实站点验证范围

本版本已经通过真实 ysepan/ys168 环境的 sandbox 联调验证，覆盖：

- 目录新增 / 更新 / 删除。
- 文件上传 / 列表 / 下载 / 删除 / 恢复。
- 子目录创建 / 解析 / 删除。
- 文件条目新增 / 更新 / 移动 / 可见性切换。
- 空间访问密码会话初始化。

## 质量门禁

- `go test ./...`
- `go vet ./...`
- `go test -race ./...`

## 模块路径

```bash
go get github.com/lib-x/ysgo@v0.1.0
```

## 说明

- live integration tests 通过环境变量显式启用。
- 排序接口已实现且协议正确，自动化验证默认避免对真实现有数据做持久排序修改。
- 上传地址默认要求 `https` 且属于 ysepan/ys168 主机域名，loopback 仅用于测试。
