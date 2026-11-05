# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### Added
- Fixed `generate` to produce compilable Ballo project.
- Added strict Go module path validation (including dots and underscores).
- Implemented atomic generation with cleanup on failure.
- Integrated `go mod tidy` and `go build` verification step.
- Updated `serve` to run a real Ballo engine with graceful shutdown.
- Added health (`/healthz`, `/health`) endpoints in generated app.
- Updated tests to cover invalid/valid module names and build verification.

### Changed
- `main.go` imports now include `errors`, `fmt`, `os`, `strings`.
- `cmd_serve.go` now starts a Ballo server with proper middleware.
- `cmd_generate.go` now validates project names and runs module tidy/build.

### Fixed
- Compilation errors due to undefined `NewApp`, `addr`.
- `serve` used custom mux instead of Ballo engine.
- Missing validation allowed path traversal.
- Generation left half‑done files on write errors.
- Signal handling blocked on `srv.Close` without graceful shutdown.

## [0.2.0] - 2026-11-05

### Added
- Initial implementation of generation and serve commands.

## [0.1.0] - 2026-10-30

- Initial CLI foundation.
