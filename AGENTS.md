# Repository Guidelines

## Project Structure & Module Organization
`main.go` is the interactive CLI entrypoint. `cmd/quick-save/main.go` provides a shortcut command that runs the quick-save flow directly. Reusable code lives under `pkg/`: `pkg/core` loads environment values, `pkg/file_manager` handles directory copy and cleanup, `pkg/io_manager` contains console input/menu logic, and `pkg/runner` orchestrates save/restore actions. Local sample data lives in `origin/` and `backups/`; both paths are also referenced from `.env`.

## Build, Test, and Development Commands
Use standard Go tooling:

- `go run .` runs the interactive restore/save menu.
- `go run ./cmd/quick-save` runs the quick-save entrypoint.
- `go build ./...` verifies all packages compile.
- `go test ./...` runs the full test suite.
- `gofmt -w main.go cmd/quick-save/main.go pkg/**/*.go` formats touched files before review.

The application expects `.env` to define `ORIGIN_DIR` and `BACKUP_DIR` with local filesystem paths.

## Coding Style & Naming Conventions
Follow idiomatic Go: tabs for indentation, `gofmt` output as the source of truth, exported identifiers in `PascalCase`, internal helpers in `camelCase`. Preserve the existing package layout and names such as `file_manager` and `io_manager` rather than renaming folders opportunistically. Keep CLI output short and action-oriented. Prefer small functions in `pkg/runner` and isolate filesystem behavior in `pkg/file_manager`.

## Testing Guidelines
There are currently no committed `_test.go` files, so new contributions should add focused unit tests with the code they change. Place tests beside the package under test, use Go's `testing` package, and prefer table-driven cases where branches are simple. For filesystem logic, use `t.TempDir()` instead of the checked-in `origin/` or `backups/` folders. Run `go test ./...` before opening a PR.

## Commit & Pull Request Guidelines
Recent history uses short, scope-first subjects such as `Saver, Restorer :: imp` and `Console List :: imp`. Keep commits concise, imperative, and tied to one change. For pull requests, include:

- a short summary of behavior changes,
- the commands you ran (`go build ./...`, `go test ./...`),
- any `.env` or path assumptions reviewers need to reproduce the change.

This is a terminal application, so screenshots are usually unnecessary; include terminal output only when it clarifies a menu or restore/save flow change.
