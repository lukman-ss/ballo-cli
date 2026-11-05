# ballo-cli

Minimal CLI for scaffolding and running [Ballo](https://github.com/lukman-ss/ballo) applications.

## Features

- `serve` — HTTP server with health probes
- `generate` — scaffold a new Ballo project
- Structured exit codes (0 success, 1 error, 2 usage, 3 fatal)

## Install

```sh
go install github.com/lukman-ss/ballo-cli@latest
```

## Usage

```
ballo-cli serve [:8080]       # start server
ballo-cli generate <name>     # scaffold project
ballo-cli help
```

## Tests

```sh
go test ./...
```

## License

MIT
