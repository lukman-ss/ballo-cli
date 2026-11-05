# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### Added
- Test `TestGenerateDefaultOffline`: verify generate is offline by default.
- Test `TestVerificationSuccess`: verify `testGeneratedProject` succeeds when valid.
- Test `TestServeNoGoMod`: verify `serve` fails when no `go.mod` present.

### Changed
- `cmdGenerate`: default offline; `--verify` flag triggers `go mod tidy`/`go build`.
- `cmdGenerate`: on verification failure, reports project location + cause; keeps project directory.
- `cmdServe`: runs current project binary (not CLI's Ballo server); builds and starts app in cwd.
- `README.md`: reflect offline-by-default, `--verify` flag, `serve` runs user project.

### Removed
- `--skip-verify` flag (replaced by `--verify`).
- Stale health-check-only `serve` implementation (now runs user project).

## [0.2.0] - 2026-11-05

### Added
- Initial implementation of generation and serve commands.

## [0.1.0] - 2026-10-30

- Initial CLI foundation.
