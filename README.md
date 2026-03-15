# Go Restorer

`go-restorer` is a small terminal utility for creating folder snapshots as `.zip` archives
and restoring them later. It is intended for local workflows where you want a quick way to save
the current state of a directory, experiment, and roll back to an earlier backup.

The application works with two filesystem locations:

- `ORIGIN_DIR`: the directory you want to save and restore
- `BACKUP_DIR`: the directory where snapshots are stored

Each saved backup is created as a `.zip` archive with a numeric prefix such as
`1.first.zip`, `2.quick_save.zip`, `3.before-refactor.zip`.

## Features

- Interactive terminal menu for `Save`, `Restore`, and `List`
- Quick-save entrypoint for one-command snapshots
- Backup storage as `.zip` archives
- Restore flow that unpacks an archive into a temporary directory before replacing the origin folder
- Plain text fallback when `/dev/tty` is not available
- `golangci-lint` configuration and CI workflow included

## Project Layout

- `main.go`: main interactive CLI entrypoint
- `cmd/quick-save/main.go`: shortcut command that runs quick-save directly
- `pkg/core`: environment loading
- `pkg/file_manager`: archive create/extract, copy, overwrite, cleanup, backup listing
- `pkg/io_manager`: menu rendering and console input
- `pkg/runner`: save and restore orchestration
- `origin/`: sample source directory for local experiments
- `backups/`: sample archive directory for local experiments

## Requirements

- Go `1.22+`
- A `.env` file in the project root

Example `.env`:

```env
ORIGIN_DIR=./origin
BACKUP_DIR=./backups
```

You can also point both variables to absolute paths:

```env
ORIGIN_DIR=/home/user/projects/app
BACKUP_DIR=/home/user/projects/app-backups
```

## Run

Interactive menu:

```bash
go run .
```

Quick save:

```bash
go run ./cmd/quick-save
```

Build:

```bash
go build ./...
```

Test:

```bash
go test ./...
```

Lint:

```bash
golangci-lint run
```

## Usage

### 1. Save a backup

Run:

```bash
go run .
```

Choose `Save`, then enter a name:

```text
Choose action:
1. Save
2. Restore
3. List
Choice: 1
Enter backup name: before-config-change
Folder ./origin saved to ./backups/4.before-config-change.zip
```

The application automatically prefixes the backup with the next numeric index
and stores it as a `.zip` file.

### 2. Restore a backup

Run:

```bash
go run .
```

Choose `Restore`, then select one of the saved backups:

```text
Choose action:
1. Save
2. Restore
3. List
Choice: 2
Choose restore file:
1. 1.initial
2. 2.quick_save
3. 4.before-config-change
Choice: 3
Folder ./origin restored from ./backups/4.before-config-change.zip
```

Before restoring, the current state of `ORIGIN_DIR` is copied into a temporary snapshot directory.
The selected archive is unpacked into a separate temporary directory, and only then copied into `ORIGIN_DIR`.
If the final restore copy fails, the application attempts to roll back to the temporary snapshot.

### 3. List backups

Run:

```bash
go run .
```

Choose `List`:

```text
Choose action:
1. Save
2. Restore
3. List
Choice: 3
1.initial
2.quick_save
4.before-config-change
```

### 4. Quick save

For a fast snapshot without opening the menu:

```bash
go run ./cmd/quick-save
```

This creates a backup named like:

```text
Folder ./origin saved to ./backups/5.quick_save.zip
```

## Interaction Modes

When the program can access `/dev/tty`, it uses the arrow-key menu.
When `/dev/tty` is unavailable, it falls back to a plain numbered prompt through standard input.

This makes the tool usable both in an interactive terminal and in simpler shell environments.

## How Backups Are Created

- `Save` packs `ORIGIN_DIR` into a new `.zip` archive under `BACKUP_DIR`
- Existing archive files are not overwritten
- The next backup index is calculated from the highest existing numeric prefix
- `Restore` unpacks the selected archive and replaces the contents of `ORIGIN_DIR`

## Development Notes

Formatting:

```bash
gofmt -w main.go cmd/quick-save/main.go pkg/**/*.go
```

Useful local commands:

```bash
go test ./...
go build ./...
golangci-lint run
```

## Limitations

- The tool assumes `.env` exists in the repository root
- Backups are standard `.zip` archives rather than incremental backups
- Very large folders may take longer to save and restore because of archive creation and extraction
