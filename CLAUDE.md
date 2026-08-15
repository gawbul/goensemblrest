# CLAUDE.md

> See [AGENTS.md](file:///Users/gawbul/Documents/Code/goensemblrest/AGENTS.md) for the complete, authoritative guide to this codebase, architecture, endpoint catalog, and contributing rules.

---

## Quick Reference

`goensemblrest` is an idiomatic, zero-external-dependency Go client library for the [Ensembl REST API](https://rest.ensembl.org/).

- **Language / Version**: Go `1.26.x` (target), compatible with Go 1.24+
- **Dependencies**: Go Standard Library only (no third-party dependencies).

### Essential Commands

```bash
# Run lint, race-detected tests, and build (default)
make all

# Run unit tests with race detection
make test-race

# Run linter
make lint

# Format code
make format

# Build all packages and examples
make build

# Run interactive example
make example

# Run live smoke tests against rest.ensembl.org
make test-live
```

For complete architecture details, design patterns, testing strategies, and guidelines for adding endpoints, refer to [AGENTS.md](file:///Users/gawbul/Documents/Code/goensemblrest/AGENTS.md).
