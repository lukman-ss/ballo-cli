# ballo-cli

CLI for scaffolding and running [Ballo](https://github.com/lukman-ss/ballo) applications.

## Features

- `serve` — Ballo HTTP server with health probes, graceful shutdown (SIGINT/SIGTERM)
- `generate` — scaffold valid Ballo projects with proper module paths; verification optional
- Structured exit codes (0 success, 1 error, 2 usage, 3 fatal)
- Path traversal protection, Go module validation, atomic generation (no partial writes)
- Works offline by default; `go mod tidy`/`go build` verification is opt-in via normal flow or skippable with `--skip-verify`

## Install

```sh
go install github.com/lukman-ss/ballo-cli@latest
```

## Usage

```
ballo-cli serve [:8080]       # start Ballo server with graceful shutdown
ballo-cli generate <name> [flags]   # scaffold project (validates module path, compiles)
  --skip-verify                     skip go mod tidy and build checks (useful offline)
ballo-cli help                # show this help
```

### `generate`

Creates a new Ballo project with:

- Valid `go.mod` with project module path
- Working `main.go` using `ballo.New()` with Recovery middleware
- Example routes (GET /health, 404/405 handlers)
- Graceful shutdown (SIGINT/SIGTERM)
- `.env.example` template

Project names must be valid Go module paths (e.g. `myapp`, `github.com/user/myapp`).
Path traversal attempts are rejected.

### `serve`

Runs a Ballo application server with `/healthz` and `/health` endpoints. Uses
graceful shutdown with 5-second timeout on SIGINT/SIGTERM.

## Tests

```sh
go test ./...
```

Generated projects are automatically built/verified during generation.

## Changelog

See [CHANGELOG.md](CHANGELOG.md).

## License

MIT
