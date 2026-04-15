# Changelog

All notable changes to this project will be documented in this file.

The format is based on Keep a Changelog, and this project follows Semantic Versioning in spirit.

## [0.1.0] - 2026-04-15

### Added
- Session lifecycle helpers: `PrepareSession`, `PrepareAdminSession`.
- Default HTTP timeout and `...Context` variants across major network APIs.
- Space access password support via `WithSpacePassword(...)` and `VerifySpacePassword(...)`.
- Structured directory and file listing APIs.
- Upload token retrieval and multipart chunked upload support.
- Download URL generation and real file download helpers.
- Subdirectory discovery, creation, and deletion support.
- File-entry operations: add, update, move, visibility toggle, and restore deleted file.
- Sorting operations for directories and files.
- Functional options for timeout, management directory, and allowed upload hosts.
- Live integration test skeletons and CI workflow.
- Example tests and expanded README documentation.

### Changed
- Replaced timestamp-based auth assumptions with the real ysepan token/session flow.
- Centralized HTTP request/response handling.
- Hardened upload handling with host/scheme validation.
- Added mutex protection around mutable client session state.

### Fixed
- Error mapping for ysepan-specific status/body semantics, including admin-required and space-password-required cases.
- Real protocol handling for directory CRUD, file deletion, upload responses, and download URL construction.

### Security
- Prevents secrets from being stored in source defaults.
- Validates upload destination hosts by default.
