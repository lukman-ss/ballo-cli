# ballo-cli

CLI for scaffolding and running [Ballo](https://github.com/lukman-ss/ballo) applications.

## Features

- `serve` — runs the Go application in the current project directory (builds and starts it), with SIGINT/SIGTERM handling
- `generate` — scaffold valid Ballo projects with proper module paths; offline by default
- Structured exit codes (0 success, 1 error, 2 usage, 3 fatal)
- Path traversal protection, Go module validation, atomic generation (no partial writes)
- Offline by default; `go mod tidy`/`go build` verification opt-in via `--verify`

## Install

```sh
go install github.com/lukman-ss/ballo-cli@latest
```

## Usage

```
ballo-cli serve [:8080]         # run current project (builds and starts the app in cwd)
ballo-cli generate <name> [flags]   # scaffold project (validates module path, compiles if --verify)
  --verify                          run go mod tidy and build checks (default: skip)
ballo-cli help                  # show this help
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

By default `generate` is offline — it only writes files. Pass `--verify` to run
`go mod tidy` and `go build` against the generated project. If verification
fails the project directory is kept; the error message reports the project
location and the cause.

### `serve`

Runs the Go application in the current directory: it builds the project and
starts the resulting binary. On SIGINT/SIGTERM the process is terminated and
gracefully shut down. Requires a `go.mod` in the current directory; use
`generate` to create a project first.

## Tests

```sh
go test ./...
```

Generated projects are only verified when `--verify` is passed.

## Changelog

See [CHANGELOG.md](CHANGELOG.md).

## License

MIT
