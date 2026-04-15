# ysgo v0.1.0

The first production-usable release of `ysgo`.

## Highlights

- Real ysepan/ys168 session lifecycle support.
- Optional space access password support.
- Typed directory and file listing APIs.
- Real upload, download, delete, restore, and subdirectory workflows.
- File entry operations: add, update, move, visibility toggle.
- Functional-option based client configuration.
- Context-aware APIs and default timeout.
- Concurrency-safe mutable client state.
- CI and live integration test scaffolding.

## Recommended client construction

```go
client := ysgo.NewClient(
    user,
    managePass,
    ysgo.WithSpacePassword("<space-password>"),
    ysgo.WithTimeout(15*time.Second),
    ysgo.WithManagementDirectory("1445856"),
)
```

## Verified against real site

This release has been exercised against real ysepan/ys168 flows using sandbox-only mutations for:

- directory create / update / delete
- file upload / list / download / delete / restore
- subdirectory create / discover / delete
- file entry add / update / move / visibility toggle
- space-password gated session initialization

## Quality gates

- `go test ./...`
- `go vet ./...`
- `go test -race ./...`

## Operational notes

- Live integration tests are env-gated by design.
- Sorting APIs are implemented and protocol-tested, while real-site persistent sort mutation remains conservative in automated verification.
- Upload host validation defaults to `https` + ysepan/ys168 allowlist, with loopback allowed for tests.
