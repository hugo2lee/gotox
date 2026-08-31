# gotox Agent Rules

These are project-specific rules for coding agents working in this repository. They complement the general engineering philosophy at https://github.com/hugo2lee/engineering-philosophy.

## Architecture

- Preserve the current package-by-capability structure. Do not introduce Clean Architecture, DDD, ports/adapters, repository, or use-case directories without concrete change pressure that cannot be handled by the existing package boundaries.
- Prefer concrete structs and functions inside a package. Define interfaces at the consuming boundary when substitution, independent testing, or ownership separation is actually required.
- Never add mutable package-level runtime dependencies such as loggers, database handles, clients, clocks, or configuration. Use explicit constructor/builder/call injection.
- Keep orchestration implementation details with the orchestrator. Do not leak `sync.WaitGroup`, worker-pool mechanics, or similar coordination primitives into capability contracts.

## Infrastructure construction

- Distinguish pure wiring from connection creation.
- Use `NewWith...`-style constructors when a caller already owns a concrete client/database.
- Use `Dial`-style functions for network/database creation and verification.
- Preserve compatibility wrappers for exported APIs unless a breaking change is explicitly approved.
- If a partial connection is created and verification fails, clean it up before returning the error.

## Persistence

- `ormx.GormModel` is a persistence convenience type. Do not treat it as a domain entity base class.
- Do not introduce new generic persistence base entities that leak ORM types/tags into domain behavior.

## Tests and verification

Before considering a change complete, run the repository default baseline:

```bash
go vet ./...
go test ./...
go test -race ./...
go test -tags=integration -run '^$' ./...
```

Default tests must be deterministic and must not require private config, external network endpoints, Redis, MongoDB, MySQL, or PostgreSQL.

Tests requiring external infrastructure must be isolated with the `integration` build tag. Do not weaken or bypass the default baseline to accommodate an integration dependency.

Do not commit real credentials, tokens, secrets, or live-service authentication data into tests or examples.

## Change discipline

- Prefer small vertical changes with an observable reason over repository-wide rewrites.
- Keep behavior-preserving compatibility changes separate from breaking API removals.
- Do not combine mechanical formatting cleanup with architectural changes unless formatting is required for the changed files.
- When CI reveals pre-existing debt, separate it from the current behavior change and make the boundary explicit rather than hiding the failure.
