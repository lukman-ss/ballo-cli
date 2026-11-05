# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### Added
- `--skip-verify` flag to `generate` for offline/skip build checks (default: verify).
- Test `TestTestGeneratedProjectFailureDoesNotDeleteDir`: verification failure does not delete project directory.
- Test `TestGenerateSkipVerify`: verify `--skip-verify` flag works.
- Test `TestGenerateOfflineDoesNotRequireNetwork`: generate works offline by default.
- Cleanup: removed stale `myapp/go.mod` and `myapptest/go.mod` fixtures; now created in `t.TempDir()`.

### Changed
- `cmdGenerate`: add `--skip-verify` flag; verification failure no longer deletes generated project.
- `testGeneratedProject`: on failure, report error without deleting project directory.
- Documentation: README.md updated to describe optional verification and `--skip-verify` flag.

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
