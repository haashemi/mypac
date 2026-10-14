# AGENTS.md

mypac is a small Go CLI for Arch Linux. It compares what is on the system with a hand-written YAML config: package groups, reasons for unowned ("lost") files, and ignores for modified package files.

## Layout

```
main.go              subcommand dispatch, root check, usage text, shared colors, newFlags (-config)
packages.go          `packages` (default): groups `pacman -Qqe` by config.Packages
lostfiles.go         `lostfiles`: groups `lostfiles` output by config.LostFiles reasons
modified.go          `modified`: lists `[modified]` backup files from `pacman -Qii`
internal/config/     config search paths (/etc only), YAML parsing (incl. legacy bare-list format), LostFileReason
internal/arch/       every external command (pacman, lostfiles) and its output parsing
PKGBUILD             `mypac-git` package; builds from GitHub `main`, so local changes need pushing first
```

## Conventions

- mypac runs as root only (`sudo mypac`); `main` refuses otherwise, except for `help`. Never shell out to `sudo`, and don't read user-home paths (`~`, `$XDG_CONFIG_HOME`), since under sudo they're root's.
- Stdlib `flag` with one `flag.FlagSet` per subcommand; no CLI framework. Add a command by writing `runX(args []string) error` and registering it in `commands` in `main.go`, plus a line in `usage()`.
- All `exec.Command` calls live in `internal/arch`. Keep parsing in pure functions (`parseModified`) so they're testable without pacman.
- Commands parsing pacman output set `LC_ALL=C`.
- Output: group headers with `boldGreenP`, items `blueP` indented by `prefixSpace`, errors returned up to `main` (printed red, exit 1). `-paths` flags print bare paths with no decoration, for piping into `tar -T -`.
- Config changes: update the `Config` struct, `TestParse`, and the README config example together.
- Dependencies: `fatih/color` and `gopkg.in/yaml.v3` only. Prefer stdlib over adding more.

## External tools

- `pacman -Qqe` – explicitly installed packages.
- `pacman -Qii` – the `Backup Files` section marks files `[modified]`, `[unmodified]` or `[unreadable]`.
- `lostfiles [-s]` – a bash script that must run as root (a hard dependency in the PKGBUILD). Prints one absolute path per line.

## Checks

```sh
gofmt -l . && go vet ./... && go test ./...
```

Tests cover config parsing, lost-file reason matching, and `pacman -Qii` parsing. Real commands need an Arch system and root: `go build && sudo ./mypac modified`.
