# Release Notes - 0.1.0

## Summary

`ysgo` is now in a production-usable state for real ysepan/ys168 automation workflows.

This release focuses on protocol correctness, operational safety, and real-site verification rather than just mock-only SDK behavior.

## Highlights

- Real session and auth lifecycle support.
- Optional space access password support.
- Typed directory and file listing.
- Real upload, download, delete, restore, and subdirectory workflows.
- File management operations: add, update, move, visibility toggle.
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

The SDK has been exercised against real ysepan/ys168 flows using sandbox-only mutations for:

- directory create / update / delete
- file upload / list / download / delete / restore
- subdirectory create / discover / delete
- file entry add / move / visibility toggle
- space password gated session initialization

## Operational notes

- Live integration tests are env-gated by design.
- Sorting APIs are implemented and protocol-tested, but real-site persistent sort mutation is intentionally conservative in automated verification.
- Upload host validation defaults to `https` + ysepan/ys168 host allowlist, with loopback allowed for tests.
