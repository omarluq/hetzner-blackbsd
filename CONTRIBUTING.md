# Contributing to blackbsd

Thank you for your interest in contributing! This guide will help you get started.

## Development Setup

### Prerequisites

- [Crystal](https://crystal-lang.org/install/) >= 1.18.2
- [Hace](https://github.com/ralsina/hace) (task runner)
- [Lefthook](https://github.com/evilmartians/lefthook) (git hooks)
- [yamlfmt](https://github.com/google/yamlfmt) (YAML formatter)

### Getting Started

```bash
git clone https://github.com/omarluq/blackbsd.git
cd blackbsd
shards install
lefthook install
```

## Available Hace Tasks

```bash
bin/hace              # Run format check, lint, and tests (default)
bin/hace build        # Build debug binary
bin/hace build:release # Build release binary
bin/hace spec         # Run tests
bin/hace format       # Format source files
bin/hace format:check # Check formatting
bin/hace ameba        # Run linter
bin/hace ameba:fix    # Run linter with auto-fix
bin/hace clean        # Remove build artifacts
```

## Running Tests

```bash
crystal spec -v
# or
bin/hace spec
```

## Code Coverage

```bash
crystal build run_specs.cr -o bin/run_specs
kcov --clean --include-path=./src ./coverage ./bin/run_specs
```

## Pre-commit Hooks

Lefthook runs automatically on commit:

- **Crystal format** - Auto-formats `.cr` files
- **Ameba** - Lints Crystal code
- **yamlfmt** - Formats YAML files

## Code Style

- Follow Crystal's standard formatting (`crystal tool format`)
- Keep code clean per Ameba's rules
- Write specs for new functionality

## Submitting Changes

1. Fork the repository
2. Create a feature branch (`git checkout -b my-feature`)
3. Make your changes with tests
4. Ensure all checks pass (`bin/hace`)
5. Commit your changes
6. Push and open a Pull Request
