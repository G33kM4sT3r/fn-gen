# fn-gen

> A deterministic feature name generator CLI tool that creates creative, reproducible feature names for your projects. (Don't take it too seriously.)

[![Go](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go&logoColor=white)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Buy Me A Coffee](https://img.shields.io/badge/Buy_Me_A_Coffee-FFDD00?logo=buy-me-a-coffee&logoColor=black)](https://buymeacoffee.com/martin.willig)

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Flags](#flags)
- [Modes](#modes)
- [Seed Mechanism](#seed-mechanism)
- [Project Structure](#project-structure)
- [Development](#development)
- [License](#license)

## Features

- **Deterministic Generation** – Same seed always produces the same name
- **Multiple Languages** – Support for English (`en`) and German (`de`)
- **4 Creative Modes** – From minimal to full buzzword bingo
- **Batch Generation** – Generate multiple names at once
- **Zero Dependencies** – Pure Go, no external runtime required

## Installation

### From Source

```bash
git clone https://github.com/G33kM4sT3r/fn-gen.git
cd fn-gen
make build
```

The binary will be available in `./dist/fn-gen`.

### Install to PATH

```bash
make install
```

### Cross-Platform Builds

```bash
make build-all
```

Builds for: `linux/amd64`, `linux/arm64`, `darwin/amd64`, `darwin/arm64`, `windows/amd64`

## Usage

```bash
fn-gen [flags]
```

### Examples

```bash
# Generate a single startup-style name (default)
fn-gen
# → "Dynamic Workflow Hub"

# Generate 5 enterprise names in German
fn-gen -lang de -mode enterprise -count 5

# Generate a reproducible name with a seed
fn-gen -seed "my-feature-123"
# → Always produces the same name

# Minimal mode for simpler names
fn-gen -mode minimal
# → "Scalable Core"
```

## Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `-lang` | string | `en` | Language for word selection (`en`, `de`) |
| `-mode` | string | `startup` | Generation mode (see [Modes](#-modes)) |
| `-seed` | string | `""` | Deterministic seed for reproducible output |
| `-count` | int | `1` | Number of names to generate |

### Flag Details

#### `-lang`

Selects the language for the word pool. Each language has its own set of adjectives, buzzwords, core terms, and suffixes optimized for that language.

```bash
fn-gen -lang en  # English: "Smart Workflow Engine"
fn-gen -lang de  # German: "Dynamische Daten Pipeline"
```

#### `-mode`

Controls the complexity and style of generated names. See [Modes](#modes) for details.

#### `-seed`

Enables deterministic name generation. When provided, the same seed will always produce the same output, making it perfect for:
- Feature branch naming
- Reproducible builds
- Team coordination
- CI/CD pipelines

#### `-count`

Generate multiple names in a single invocation:

```bash
fn-gen -count 3
# → "Scalable Event Processor"
# → "Cloud-Native Data Hub"
# → "Modular Workflow Engine"
```

## Modes

Each mode defines a pattern that determines which word categories are combined:

| Mode | Pattern | Example Output |
|------|---------|----------------|
| `minimal` | adjective + core | "Scalable Core" |
| `startup` | adjective + core + suffix | "Dynamic Workflow Hub" |
| `enterprise` | adjective + buzzword + core + suffix | "Unified Cloud Integration Platform" |
| `bullshit` | adjective + buzzword + buzzword + core + suffix | "Synergized AI-Powered Blockchain Data Engine" |

### Word Categories

Each mode draws from these word pools:

- **Adjectives** – Descriptive words (Smart, Dynamic, Scalable, ...)
- **Buzzwords** – Trendy tech terms (Cloud, AI-Assisted, Serverless, ...)
- **Core** – Central concept words (Workflow, Data, Integration, ...)
- **Suffix** – Ending words (Hub, Engine, Platform, ...)

## Seed Mechanism

The seed mechanism is the heart of fn-gen's deterministic generation. Understanding how it works helps you leverage it effectively.

### How It Works

```
┌─────────────────────────────────────────────────────────────┐
│                      Seed Construction                      │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│   With custom seed:    seed = "my-feature-123"              │
│                                                             │
│   Without seed:        seed = "{lang}-{mode}-{index}-{date}"│
│                        Example: "en-startup-0-2026-01-15"   │
│                                                             │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                   Per-Word Hash Generation                  │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│   For each word position (i) and category (key):            │
│                                                             │
│   hashInput = "{seed}-{i}-{key}"                            │
│   Example:   "my-feature-123-0-adjectives"                  │
│                                                             │
│   hash = SHA256(hashInput)                                  │
│   index = hash[0:8] as uint64                               │
│   selectedWord = wordList[index % len(wordList)]            │
│                                                             │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                       Final Output                          │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│   parts = ["Dynamic", "Workflow", "Hub"]                    │
│   output = "Dynamic Workflow Hub"                           │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### Seed Properties

1. **Deterministic**: Same seed + same flags = same output, always
2. **Distributed**: SHA-256 ensures uniform distribution across word lists
3. **Independent per position**: Each word in the pattern is hashed independently
4. **Collision-resistant**: Different seeds produce different outputs

### Automatic Seed (No `-seed` flag)

When no seed is provided, fn-gen creates an automatic seed using:

```
{lang}-{mode}-{index}-{date}
```

This means:
- Names are reproducible **within the same day** (when called with same parameters)
- Different indices (`-count > 1`) produce different names
- Tomorrow's names will be different

### Custom Seed Use Cases

| Use Case | Seed Strategy | Example |
|----------|---------------|---------|
| Feature branches | Issue/ticket ID | `-seed "JIRA-1234"` |
| Reproducible demos | Static string | `-seed "demo-2024"` |
| User-specific | Username + date | `-seed "alice-2024-01"` |
| Build artifacts | Commit hash | `-seed "a1b2c3d4"` |

### Combining Seed with Count

When using `-count > 1`, the index is incorporated into the hash:

```bash
fn-gen -seed "project-x" -count 3
```

Internally generates:
```
Hash("project-x" + "-0-" + "adjectives") → word 1 for name 1
Hash("project-x" + "-1-" + "adjectives") → word 1 for name 2
Hash("project-x" + "-2-" + "adjectives") → word 1 for name 3
```

This ensures each name in a batch is unique but still reproducible.

## Project Structure

```
fn-gen/
├── cmd/fn-gen/          # CLI entry point
│   └── main.go
├── internal/
│   ├── cli/             # Flag parsing and configuration
│   │   └── flags.go
│   ├── generator/       # Core generation logic
│   │   ├── generator.go # Name generation
│   │   ├── modes.go     # Mode patterns
│   │   └── seed.go      # Hash function
│   └── words/           # Word data and loader
│       ├── loader.go
│       └── data/
│           ├── en/      # English word sets
│           │   ├── bullshit.json
│           │   ├── enterprise.json
│           │   ├── minimal.json
│           │   └── startup.json
│           └── de/      # German word sets
│               └── ...
├── Makefile
└── README.md
```

## Development

```bash
# Run locally
make run

# Run tests
make test

# Format code
make fmt

# Lint
make lint

# Full check (fmt + vet + lint + test)
make check
```

## License

MIT License – see [LICENSE](LICENSE) for details.
