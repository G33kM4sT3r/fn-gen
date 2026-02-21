# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

fn-gen is a deterministic feature name generator CLI tool written in pure Go (zero dependencies). It creates creative, reproducible feature names using SHA-256 hashing for deterministic word selection from categorized word pools.

## Commands

```bash
make build          # Build binary → dist/fn-gen
make run            # Run locally via go run
make test           # Run tests with race detector
make lint           # Run golangci-lint
make fmt            # Format with gofmt
make vet            # Static analysis
make check          # All checks: vet + lint + test
make build-all      # Cross-compile for linux/darwin/windows
```

Run directly: `go run ./cmd/fn-gen -lang en -mode startup -seed "test" -count 3 -explain`

## Architecture

```
cmd/fn-gen/main.go → cli.ParseFlags() → words.Load(lang, mode) → generator.New() → generator.Generate()
```

**cli** (`internal/cli/flags.go`): Parses 5 flags (lang, mode, seed, count, explain) into a `Config` struct using the standard `flag` package.

**words** (`internal/words/`): Loads JSON word files from `internal/words/data/{lang}/{mode}.json`. Each `WordSet` has 4 categories: adjectives, buzzwords, core, suffix. The `Get(key)` method provides dynamic category access.

**generator** (`internal/generator/`): Core logic. Each mode defines a pattern (sequence of word category keys). For each position, it hashes `"{seed}-{position}-{category}"` with SHA-256, takes first 8 bytes as uint64, and mods by list length to select a word. Without a custom seed, auto-generates `"{lang}-{mode}-{index}-{YYYY-MM-DD}"`.

**modes** (`internal/generator/modes.go`): Maps mode names to category patterns — minimal (2 words), startup (3), enterprise (4), bullshit (5).

## Key Design Decisions

- Word data files are read from disk at runtime via `os.ReadFile`, not embedded. The binary must be run from the project root or the data path must be accessible.
- Version, commit hash, and build time are injected via LDFLAGS at build time.
- Languages (en, de) each have 4 mode-specific JSON files with identical structure.
- No tests exist yet; the Makefile `test` target is ready (`go test -race ./...`).
