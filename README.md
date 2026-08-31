# gotox

`gotox` is a small Go infrastructure toolkit: configuration, logging, cache, HTTP helpers, resource lifecycle, Redis, MongoDB, and GORM integration.

The project follows the principles in [hugo2lee/engineering-philosophy](https://github.com/hugo2lee/engineering-philosophy): keep the architecture small, make meaningful boundaries explicit, prefer concrete types by default, and add abstractions only when real change pressure justifies them.

## Boundary rules

- Keep infrastructure packages concrete. Do not add domain/application/port layers inside `gotox` merely to imitate hexagonal architecture.
- Do not introduce mutable package-level runtime dependencies. Dependencies such as loggers belong on instances/builders or explicit call paths.
- Constructors used for pure wiring must not perform hidden network/database I/O. Use `NewWith...` for existing clients and `Dial...`/`Dial` for connection creation and verification.
- Provider packages should not force broad provider-owned interfaces on consumers. Consumers should define the smallest interface they actually need.
- Keep orchestration details inside the owner. For example, `Resource.Close` reports an error; `ResourcexGroup` owns concurrency and waiting.
- GORM-specific convenience models are persistence types, not domain entity base types.

## Tests

The default test suite must work from a fresh checkout without private configuration or external services:

```bash
go test ./...
go test -race ./...
go vet ./...
```

Tests that require Redis, MongoDB, MySQL, PostgreSQL, private configuration, or another external service must use the `integration` build tag. The integration suite can be compile-checked without starting those services:

```bash
go test -tags=integration -run '^$' ./...
```

Run actual integration tests only in an environment where the required services and configuration are explicitly provided:

```bash
go test -tags=integration ./...
```

## Compatibility policy

Prefer incremental migrations for exported APIs. When an old API is misleading but widely usable, introduce the clearer API first, keep a compatibility wrapper or alias when practical, mark it deprecated, and remove it only with explicit breaking-change evidence and release intent.
